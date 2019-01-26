package models

type TaskComment struct {
	Model
	Content string `gorm:"not null;"`

	Task   Task
	TaskID uint `gorm:"not null;"`
}

func GetTaskCommentBySlug(slug string) (*TaskComment, error) {
	var taskCom TaskComment
	if err := GetObjectByField(&taskCom, "slug", slug); err != nil {
		return nil, err
	}
	return &taskCom, nil
}

func CreateTaskComment(content string, taskID uint, user *User) (*TaskComment, error) {
	taskCom := TaskComment{
		Content: content,
		TaskID:  taskID,
	}
	taskCom.SetCreatedBy(user)
	if err := db.Create(&taskCom).Error; err != nil {
		return nil, err
	}
	return &taskCom, nil
}
