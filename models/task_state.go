package models

type TaskState struct {
	Model
	Name     string
	Sequence int
	Tasks    []Task

	Project   Project
	ProjectID int
}
