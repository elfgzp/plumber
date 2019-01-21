package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/serializers"
	"net/http"
	"strconv"
)

type TaskListCreate struct {
	Name string `json:"name"`
}

func CheckTaskListCreate(taskListCreate TaskListCreate) []FieldCheckError {
	var errs []FieldCheckError
	if taskListCreate.Name == "" {
		errs = append(errs, FieldCheckError{Field: "name", Error: "Task list name required."})
	}

	if taskListNameExsit(taskListCreate.Name) {
		errs = append(errs, FieldCheckError{Field: "name", Error: "Task list name existed."})
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

	if errs := CheckTaskListCreate(taskListCreate); len(errs) > 0 {
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

	taskList, err := models.GetTaskListBySlug(params["taskSlug"])
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
	Sequence string `json:"sequence"`
}

func CheckTaskListUpdate(taskListUpdate TaskListUpdate, taskList *models.TaskList) (map[string]interface{}, []FieldCheckError) {
	var errs []FieldCheckError
	contents := map[string]interface{}{}
	if taskListUpdate.Name != "" && taskListUpdate.Name != taskList.Name {
		if taskListNameExsit(taskListUpdate.Name) {
			errs = append(errs, FieldCheckError{Field: "name", Error: "Task list name existed."})
		} else {
			contents["name"] = taskListUpdate.Name
		}
	}

	if taskListUpdate.Sequence != "" {
		if sequence, err := strconv.Atoi(taskListUpdate.Sequence); err == nil {
			if sequence != taskList.Sequence {
				contents["sequence"] = sequence
			}
		} else {
			errs = append(errs, FieldCheckError{Field: "sequence", Error: "Sequence must be a number."})
		}

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

	taskList, err := models.GetTaskListBySlug(params["taskSlug"])
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

	contents, errs := CheckTaskListUpdate(taskListUpdate, taskList)
	if len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if err := models.UpdateObject(&taskList, contents); err != nil {
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

	taskList, err := models.GetTaskListBySlug(params["taskSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if taskList == nil {
		helpers.Response404(w, "Task list not found.")
		return
	}

	if err := models.FakedDestroyObject(&taskList); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response204(w)
	}
}
