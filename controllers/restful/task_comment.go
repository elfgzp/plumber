package restful

import (
	"encoding/json"
	"github.com/elfgzp/plumber/helpers"
	"github.com/elfgzp/plumber/models"
	"github.com/elfgzp/plumber/serializers"
	"io/ioutil"
	"net/http"
)

type CreateTaskComment struct {
	Content string `json:"content"`
}

func validCreateTaskComment(createTaskCom CreateTaskComment) []FieldValidError {
	var errs []FieldValidError
	if createTaskCom.Content == "" {
		errs = append(errs, FieldValidError{Field: "content", Error: "Content required."})
	}
	return errs
}

func CreateTaskCommentHandler(w http.ResponseWriter, r *http.Request) {
	var createTaskCom CreateTaskComment

	ru := getRequestUser(r)
	params := getRouteParams(r)
	if err := json.NewDecoder(r.Body).Decode(&createTaskCom); err != nil {
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

	if errs := validCreateTaskComment(createTaskCom); len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if taskCom, err := models.CreateTaskComment(createTaskCom.Content, task.ID, ru); err != nil {
		helpers.Response500(w)
		return
	} else {
		helpers.Response201(w, "", serializers.SerializeTaskComment(taskCom))
		return
	}
}

func validUpdateTaskComment(data map[string]interface{}) (map[string]interface{}, []FieldValidError) {
	var errs []FieldValidError
	contents := map[string]interface{}{}

	contents, errs = validStringField(contents, data, "content", errs)
	if content, ok := contents["content"]; ok {
		if content == "" || content == nil {
			errs = append(errs)
		}
	}
	return contents, errs
}

func UpdateTaskCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	taskComment, _ := models.GetTaskCommentBySlug(params["taskCommentSlug"])

	if taskComment == nil {
		helpers.Response404(w, "Task comment not found.")
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

	contents, errs := validUpdateTaskComment(data)
	if len(errs) > 0 {
		helpers.Response400(w, "", errs)
		return
	}

	if err := models.UpdateObject(&taskComment, contents, ru); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response200(w, "", serializers.SerializeTaskComment(taskComment))
	}
}

func ListTaskCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	taskComments, total, _ := task.GetTaskCommentsLimit(page, limit)
	var taskCommenti []interface{}
	if total != 0 {
		taskCommenti = make([]interface{}, len(*taskComments))
		for i, taskComment := range *taskComments {
			taskCommenti[i] = serializers.SerializeTaskComment(&taskComment)
		}
	} else {
		taskCommenti = make([]interface{}, 0)
	}
	helpers.Response200(w, "", helpers.PagedData{Total: total, Page: page, Limit: limit, Result: taskCommenti})
}

func RetrieveTaskCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	taskComment, _ := models.GetTaskCommentBySlug(params["taskCommentSlug"])

	if taskComment == nil {
		helpers.Response404(w, "Task comment not found.")
		return
	}

	helpers.Response200(w, "", serializers.SerializeTaskComment(taskComment))
}

func DestroyTaskCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	taskComment, err := models.GetTaskCommentBySlug(params["taskCommentSlug"])
	if err != nil {
		helpers.Response500(w)
		return
	}

	if taskComment == nil {
		helpers.Response404(w, "Task comment not found.")
		return
	}

	if err := models.FakedDestroyObject(&taskComment, ru); err != nil {
		helpers.Response500(w)
	} else {
		helpers.Response204(w)
	}
}
