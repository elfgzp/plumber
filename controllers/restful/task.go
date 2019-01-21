package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/serializers"
	"net/http"
	"time"
)

type CreateTask struct {
	Name             string    `json:"name"`
	Deadline         time.Time `json:"deadline"`
	AssignedUserSlug string    `json:"assigned_user_slug"`
}

func checkTaskCreate(createTask CreateTask) (map[string]interface{}, []FieldCheckError) {
	var errs []FieldCheckError
	contents := map[string]interface{}{}

	if createTask.Name == "" {
		errs = append(errs, FieldCheckError{"name", "Name required."})
	}

	if createTask.AssignedUserSlug != "" {
		if assign, _ := models.GetUserBySlug(createTask.AssignedUserSlug); assign == nil {
			errs = append(errs, FieldCheckError{"assigned_user_slug", "Assigned user not Found."})
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

	if contents, errs := checkTaskCreate(createTask); len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	} else if len(contents) > 0 {
		if val, ok := contents["assign_id"]; ok {
			assignID = val.(*uint)
		}
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

	if task, err := models.CreateTask(createTask.Name, createTask.Deadline, assignID, taskList, ru); err != nil {
		helpers.Response500(w)
		return
	} else {
		helpers.Response201(w, "", serializers.SerializeTask(task))
		return
	}

}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {

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
			ti[i] = serializers.SerializeTask(&t)
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
		helpers.Response404(w, "Task list not found.")
		return
	}

	helpers.Response200(w, "", serializers.SerializeTask(task))
}

func DestroyTaskHandler(w http.ResponseWriter, r *http.Request) {

}
