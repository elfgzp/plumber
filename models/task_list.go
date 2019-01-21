package models

type TaskList struct {
	Model
	Name     string
	Sequence int `gorm:"AUTO_INCREMENT"`
	Tasks    []Task

	Project   Project
	ProjectID uint `gorm:"not null;"`

	Active bool
}

func GetTaskListByName(name string) (*TaskList, error) {
	var tl TaskList
	err := GetObjectByField(&tl, "name", name)
	if err != nil {
		return nil, err
	}
	return &tl, err
}

func GetTaskListBySlug(slug string) (*TaskList, error) {
	var tl TaskList
	err := GetObjectByField(&tl, "slug", slug)
	if err != nil {
		return nil, err
	}
	return &tl, err
}

func CreateTaskList(name string, projectID uint, user *User) (*TaskList, error) {
	t := TaskList{Name: name, ProjectID: projectID}
	t.SetCreatedBy(user)
	if err := db.Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (tl *TaskList) GetTasksLimit(page, limit int) (*[]Task, int, error) {
	var tasks []Task
	total, err := GetObjectsByFieldLimit(&tasks, &Task{}, page, limit, "sequence asc, created_at desc", "task_list_id", tl.ID)
	return &tasks, total, err
}
