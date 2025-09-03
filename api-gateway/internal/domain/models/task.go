package models

import "time"

type Task struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}
