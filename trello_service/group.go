package trello_service

import "time"

type Group struct {
	Time time.Time
	Name string
	DoneTask string
	ProgressTask string			// Doing, Todo, CodeReview, Testing, Design
	SprintBacklogTask string
	PendingTask string
}

