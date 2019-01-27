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

type CreateTask struct {
	Name             string     `json:"name"`
	Deadline         *time.Time `json:"deadline"`
	AssignedUserSlug string     `json:"assigned_user_slug"`
}

func ValidTaskCreate(createTask CreateTask, taskList *models.TaskList) (map[string]interface{}, []FieldValidError) {
	var errs []FieldValidError
	contents := map[string]interface{}{}

	if createTask.Name == "" {
		errs = append(errs, FieldValidError{"name", "Name required."})
	}

	models.LoadRelatedObject(&taskList, &taskList.Project, "Project")
	if createTask.AssignedUserSlug != "" {
		if assign, _ := models.GetUserBySlug(createTask.AssignedUserSlug); assign == nil {
			errs = append(errs, FieldValidError{"assigned_user_slug", "Assigned user not Found."})
		} else if !taskList.Project.IsProjectMember(assign.ID) {
			errs = append(errs, FieldValidError{"assigned_user_slug", "Assigned user is not project member."})
		} else {
			contents["assign_id"] = &assign.ID
		}
	}

	return contents, errs
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var createTask CreateTask
	var assignID *uint

	ru := getRequestUser(r)
	params := getRouteParams(r)
	if err := json.NewDecoder(r.Body).Decode(&createTask); err != nil {
		helpers.Response400(w, "JSON invalid.", nil)
		return
	}

	taskList, _ := models.GetTaskListBySlug(params["taskListSlug"])
	if taskList == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	models.LoadRelatedObject(&taskList, &taskList.Project, "Project")
	if !taskList.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	if contents, errs := ValidTaskCreate(createTask, taskList); len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	} else if len(contents) > 0 {
		if val, ok := contents["assign_id"]; ok {
			assignID = val.(*uint)
		}
	}

	if task, err := models.CreateTask(createTask.Name, createTask.Deadline, assignID, taskList, ru); err != nil {
		helpers.Response500(w)
		return
	} else {
		helpers.Response201(w, "", serializers.SerializeRetrieveTask(task))
		return
	}

}

//type UpdateTask struct {
//	Name             *string    `json:"name"`
//	Desc             *string    `json:"desc"`
//	Sequence         *int       `json:"sequence"`
//	Deadline         *time.Time `json:"deadline"`
//	Doing            *bool      `json:"doing"`
//	Completed        *bool      `json:"completed"`
//	AssignedUserSlug *string    `json:"assigned_user_slug"`
//}

func validUpdateTask(data map[string]interface{}, task *models.Task) (map[string]interface{}, []FieldValidError) {
	var errs []FieldValidError
	contents := map[string]interface{}{}

	contents, errs = validStringField(contents, data, "name", errs)
	contents, errs = validStringField(contents, data, "desc", errs)
	contents, errs = validIntField(contents, data, "sequence", errs)
	contents, errs = validDatetimeField(contents, data, "deadline", errs)

	models.LoadRelatedObject(&task, &task.Project, "Project")
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

	contents, errs = validBoolField(contents, data, "doing", errs)
	contents, errs = validBoolField(contents, data, "completed", errs)
	return contents, errs
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	ru := getRequestUser(r)
	params := getRouteParams(r)

	taskList, _ := models.GetTaskListBySlug(params["taskListSlug"])
	if taskList == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	models.LoadRelatedObject(&taskList, &taskList.Project, "Project")
	if !taskList.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	task, err := models.GetTaskBySlug(params["taskSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if task == nil {
		helpers.Response404(w, "Task list not found.")
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

	//if err := json.NewDecoder(r.Body).Decode(&updateTask); err != nil {
	//	helpers.Response400(w, "JSON invalid.", nil)
	//	return
	//}

	contents, errs := validUpdateTask(data, task)
	if len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if err := models.UpdateObject(&task, contents, ru); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response200(w, "", serializers.SerializeRetrieveTask(task))
	}

}

func ListTaskHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)
	query := getQuery(r)
	page, limit := getPageLimitQuery(query)

	taskList, _ := models.GetTaskListBySlug(params["taskListSlug"])
	if taskList == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	models.LoadRelatedObject(&taskList, &taskList.Project, "Project")
	if !taskList.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	tasks, total, _ := taskList.GetTasksLimit(page, limit)
	var ti []interface{}
	if total != 0 {
		ti = make([]interface{}, len(*tasks))
		for i, t := range *tasks {
			ti[i] = serializers.SerializeRetrieveListTask(&t)
		}
	} else {
		ti = make([]interface{}, 0)
	}
	helpers.Response200(w, "", helpers.PagedData{Total: total, Page: page, Limit: limit, Result: ti})

}

func RetrieveTaskHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)

	taskList, _ := models.GetTaskListBySlug(params["taskListSlug"])
	if taskList == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	models.LoadRelatedObject(&taskList, &taskList.Project, "Project")
	if !taskList.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	task, err := models.GetTaskBySlug(params["taskSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if task == nil {
		helpers.Response404(w, "Task not found.")
		return
	}

	helpers.Response200(w, "", serializers.SerializeRetrieveTask(task))
}

func DestroyTaskHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)

	taskList, _ := models.GetTaskListBySlug(params["taskListSlug"])
	if taskList == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	models.LoadRelatedObject(&taskList, &taskList.Project, "Project")
	if !taskList.Project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	task, err := models.GetTaskBySlug(params["taskSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if task == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	if err := models.FakedDestroyObject(&task, ru); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response204(w)
	}
}
