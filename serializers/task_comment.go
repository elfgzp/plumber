package serializers

import (
	"github.com/elfgzp/plumber/models"
	"time"
)

type TaskCommentSerializer struct {
	ModelSerializer
	Content   string    `json:"content"`
	NickName  string    `json:"nick_name"`
	Avatar    string    `json:"avatar"`
	UserSlug  string    `json:"user_slug"`
	CreatedAt time.Time `json:"created_at"`
	Slug      string    `json:"slug"`
}

func SerializeTaskComment(taskComment *models.TaskComment) TaskCommentSerializer {
	models.LoadRelatedObject(&taskComment, &taskComment.CreatedBy, "CreatedBy")
	taskCommentSerializer := TaskCommentSerializer{
		Content:   taskComment.Content,
		NickName:  taskComment.CreatedBy.Nickname,
		Avatar:    taskComment.CreatedBy.Avatar,
		UserSlug:  taskComment.CreatedBy.Slug,
		CreatedAt: taskComment.CreatedAt,
		Slug:      taskComment.Slug,
	}
	taskCommentSerializer.serializeBaseField(taskComment)
	return taskCommentSerializer
}
