package models

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	Completed   bool
	CreatedAt   *timestamppb.Timestamp
	CompletedAt *timestamppb.Timestamp
}
