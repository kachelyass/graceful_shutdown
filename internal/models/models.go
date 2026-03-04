package models

type Task struct {
	ID    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
}
