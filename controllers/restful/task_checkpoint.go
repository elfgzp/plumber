package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/serializers"
	"io/ioutil"
	"net/http"
	"time"
)

type CreateTaskCheckPoint struct {
	Desc             string     `json:"desc"`
	Deadline         *time.Time `json:"deadline"`
	AssignedUserSlug string     `json:"assigned_user_slug"`
}

func ValidTaskCheckPointCreate(createTaskCP CreateTaskCheckPoint, task *models.Task) (map[string]interface{}, []FieldValidError) {
	var errs []FieldValidError
	contents := map[string]interface{}{}

	if createTaskCP.Desc == "" {
		errs = append(errs, FieldValidError{"desc", "Desc required."})
	}

	models.LoadRelatedObject(&task, &task.Project, "Project")
	if createTaskCP.AssignedUserSlug != "" {
		if assign, _ := models.GetUserBySlug(createTaskCP.AssignedUserSlug); assign == nil {
			errs = append(errs, FieldValidError{"assigned_user_slug", "Assigned user not Found."})
		} else if !task.Project.IsProjectMember(assign.ID) {
			errs = append(errs, FieldValidError{"assigned_user_slug", "Assigned user is not project member."})
		} else {
			contents["assign_id"] = &assign.ID
		}
	}

	return contents, errs
}

func CreateTaskCheckpointHandler(w http.ResponseWriter, r *http.Request) {
	var createTaskCP CreateTaskCheckPoint
	var assignID *uint

	ru := getRequestUser(r)
	params := getRouteParams(r)
	if err := json.NewDecoder(r.Body).Decode(&createTaskCP); err != nil {
		helpers.Response400(w, "JSON invalid.", nil)
		return
	}

	task, _ := models.GetTaskBySlug(params["taskSlug"])
	if task == nil {
		helpers.Response404(w, "Task not found.")
		return
	}

	models.LoadRelatedObject(&task, &task.Project, "Project")
	if !task.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	if contents, errs := ValidTaskCheckPointCreate(createTaskCP, task); len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	} else if len(contents) > 0 {
		if val, ok := contents["assign_id"]; ok {
			assignID = val.(*uint)
		}
	}

	if taskCP, err := models.CreateTaskCheckpoint(createTaskCP.Desc, createTaskCP.Deadline, assignID, task, ru); err != nil {
		helpers.Response500(w)
		return
	} else {
		helpers.Response201(w, "", serializers.SerializeTaskCheckpoint(taskCP))
		return
	}

}

func ValidUpdateTaskCheckpoint(data map[string]interface{}, task *models.Task) (map[string]interface{}, []FieldValidError) {
	var errs []FieldValidError
	contents := map[string]interface{}{}

	contents, errs = validStringField(contents, data, "desc", errs)
	contents, errs = validIntField(contents, data, "sequence", errs)
	contents, errs = validDatetimeField(contents, data, "deadline", errs)

	if assignedUserSlug, ok := data["assigned_user_slug"]; ok {
		if assignedUserSlug == nil {
			contents["assign_id"] = nil
		} else if _, ok = assignedUserSlug.(string); ok {
			if assign, _ := models.GetUserBySlug(assignedUserSlug.(string)); assign == nil {
				errs = append(errs, FieldValidError{Field: "assigned_user_slug", Error: "Assigned user not Found."})
			} else if !task.Project.IsProjectMember(assign.ID) {
				errs = append(errs, FieldValidError{"assigned_user_slug", "Assigned user is not project member."})
			} else {
				contents["assign_id"] = assign.ID
			}
		} else {
			errs = append(errs, FieldValidError{Field: "assigned_user_slug", Error: "Must be a string."})
		}
	}

	return contents, errs
}

func UpdateTaskCheckpointHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	ru := getRequestUser(r)
	params := getRouteParams(r)

	task, _ := models.GetTaskBySlug(params["taskSlug"])
	if task == nil {
		helpers.Response404(w, "Task not found.")
		return
	}

	models.LoadRelatedObject(&task, &task.Project, "Project")
	if !task.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	taskCP, err := models.GetTaskCheckpointBySlug(params["taskCheckpointSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if taskCP == nil {
		helpers.Response404(w, "Task Checkpoint not found.")
		return
	}
	models.LoadRelatedObject(&task, &task.Project, "Project")
	if !task.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		helpers.Response400(w, "JSON invalid.", nil)
		return
	}

	contents, errs := ValidUpdateTaskCheckpoint(data, task)
	if len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if err := models.UpdateObject(&taskCP, contents, ru); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response200(w, "", serializers.SerializeTaskCheckpoint(taskCP))
	}
}

func ListTaskCheckpointHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)
	query := getQuery(r)
	page, limit := getPageLimitQuery(query)

	task, _ := models.GetTaskBySlug(params["taskSlug"])
	if task == nil {
		helpers.Response404(w, "Task not found.")
		return
	}

	models.LoadRelatedObject(&task, &task.Project, "Project")
	if !task.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	taskCPs, total, _ := task.GetTaskCheckpointsLimit(page, limit)
	var taskCPi []interface{}
	if total != 0 {
		taskCPi = make([]interface{}, len(*taskCPs))
		for i, taskCP := range *taskCPs {
			taskCPi[i] = serializers.SerializeTaskCheckpoint(&taskCP)
		}
	} else {
		taskCPi = make([]interface{}, 0)
	}
	helpers.Response200(w, "", helpers.PagedData{Total: total, Page: page, Limit: limit, Result: taskCPi})
}

func RetrieveTaskCheckpointHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)

	task, _ := models.GetTaskBySlug(params["taskSlug"])
	if task == nil {
		helpers.Response404(w, "Task not found.")
		return
	}

	models.LoadRelatedObject(&task, &task.Project, "Project")
	if !task.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	taskCP, err := models.GetTaskCheckpointBySlug(params["taskCheckpointSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if taskCP == nil {
		helpers.Response404(w, "Task checkpoint not found.")
		return
	}

	helpers.Response200(w, "", serializers.SerializeTaskCheckpoint(taskCP))
}

func DestroyTaskCheckpointHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)

	task, _ := models.GetTaskBySlug(params["taskSlug"])
	if task == nil {
		helpers.Response404(w, "Task not found.")
		return
	}

	models.LoadRelatedObject(&task, &task.Project, "Project")
	if !task.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	taskCP, err := models.GetTaskCheckpointBySlug(params["taskCheckpointSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if taskCP == nil {
		helpers.Response404(w, "Task checkpoint not found.")
		return
	}

	if err := models.FakedDestroyObject(&taskCP, ru); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response204(w)
	}
}
