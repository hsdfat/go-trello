package trello_service

var (
	TASK_NAME_PATTERN   = `^\[(.*)\]\.\[T:(\d+)\]\.\[(.*)\]\..*$`
	DONE_LIST           = "done"
	SPRINT_BACKLOG_LIST = "sprint backlog"
	EPIC_LIST           = "epic"
)
