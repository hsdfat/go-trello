package trello_service

import "go-trello/logger"

// PrintMemberStatistics prinses the member statistics
func (c *TrelloClient) PrintMemberStatistics() {
	for _, stat := range c.MemberStats{
		logger.Debugln("member: ", stat.Name, "\ttask [done/progress/total]: ", stat.NDoneTasks, "/", stat.NProgressTasks, "/", stat.NTasks,
			"\tHour [done/progress/total]: ", stat.NDoneHours, "/", stat.NProgressHours, "/", stat.NHours)
	}
}


