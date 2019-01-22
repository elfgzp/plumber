package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/serializers"
	"net/http"
)

type TaskListCreate struct {
	Name string `json:"name"`
}

func ValidTaskListCreate(taskListCreate TaskListCreate) []FieldValidError {
	var errs []FieldValidError
	if taskListCreate.Name == "" {
		errs = append(errs, FieldValidError{Field: "name", Error: "Task list name required."})
	}

	if taskListNameExsit(taskListCreate.Name) {
		errs = append(errs, FieldValidError{Field: "name", Error: "Task list name existed."})
	}
	return errs
}

func CreateTaskListHandler(w http.ResponseWriter, r *http.Request) {
	var taskListCreate TaskListCreate

	ru := getRequestUser(r)
	params := getRouteParams(r)

	project, _ := models.GetProjectBySlug(params["projectSlug"])
	if project == nil {
		helpers.Response404(w, "Project not found.")
		return
	}

	if !project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&taskListCreate); err != nil {
		helpers.Response400(w, "JSON invalid.", nil)
		return
	}

	if errs := ValidTaskListCreate(taskListCreate); len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if taskList, err := models.CreateTaskList(taskListCreate.Name, project.ID, ru); err != nil {
		helpers.Response500(w)
		return
	} else {
		helpers.Response201(w, "", serializers.SerializeTaskList(taskList))
	}
}

func ListTaskListHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)
	query := getQuery(r)
	page, limit := getPageLimitQuery(query)

	project, _ := models.GetProjectBySlug(params["projectSlug"])
	if project == nil {
		helpers.Response404(w, "Project not found.")
		return
	}

	if !project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	taskLists, total, _ := project.GetTaskListsLimit(page, limit)
	var tli []interface{}
	if total != 0 {
		tli = make([]interface{}, len(*taskLists))
		for i, tl := range *taskLists {
			tli[i] = serializers.SerializeTaskList(&tl)
		}
	} else {
		tli = make([]interface{}, 0)
	}
	helpers.Response200(w, "", helpers.PagedData{Total: total, Page: page, Limit: limit, Result: tli})

}

func RetrieveTaskListHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)

	project, _ := models.GetProjectBySlug(params["projectSlug"])
	if project == nil {
		helpers.Response404(w, "Project not found.")
		return
	}

	if !project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	taskList, err := models.GetTaskListBySlug(params["taskListSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if taskList == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	helpers.Response200(w, "", serializers.SerializeTaskList(taskList))
}

type TaskListUpdate struct {
	Name     string `json:"name"`
	Sequence *int   `json:"sequence"`
}

func ValidTaskListUpdate(taskListUpdate TaskListUpdate, taskList *models.TaskList) (map[string]interface{}, []FieldValidError) {
	var errs []FieldValidError
	contents := map[string]interface{}{}
	if taskListUpdate.Name != "" && taskListUpdate.Name != taskList.Name {
		if taskListNameExsit(taskListUpdate.Name) {
			errs = append(errs, FieldValidError{Field: "name", Error: "Task list name existed."})
		} else {
			contents["name"] = taskListUpdate.Name
		}
	}

	if taskListUpdate.Sequence != nil {
		contents["sequence"] = *taskListUpdate.Sequence
	}

	return contents, errs
}

func UpdateTaskListHandler(w http.ResponseWriter, r *http.Request) {
	var taskListUpdate TaskListUpdate
	ru := getRequestUser(r)
	params := getRouteParams(r)

	project, _ := models.GetProjectBySlug(params["projectSlug"])
	if project == nil {
		helpers.Response404(w, "Project not found.")
		return
	}

	if !project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	taskList, err := models.GetTaskListBySlug(params["taskListSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if taskList == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&taskListUpdate); err != nil {
		helpers.Response400(w, "", "JSON invalid.")
		return
	}

	contents, errs := ValidTaskListUpdate(taskListUpdate, taskList)
	if len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if err := models.UpdateObject(taskList, contents, ru); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response200(w, "", serializers.SerializeTaskList(taskList))
	}

}

func DestroyTaskListHandler(w http.ResponseWriter, r *http.Request) {
	ru := getRequestUser(r)
	params := getRouteParams(r)

	project, _ := models.GetProjectBySlug(params["projectSlug"])
	if project == nil {
		helpers.Response404(w, "Project not found.")
		return
	}

	if !project.IsProjectMember(ru.ID) {
		helpers.Response403(w)
		return
	}

	taskList, err := models.GetTaskListBySlug(params["taskListSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if taskList == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	if err := models.FakedDestroyObject(&taskList, ru); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response204(w)
	}
}
