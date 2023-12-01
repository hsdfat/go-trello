package trello_service

import (
	"fmt"
	"go-trello/logger"
	"regexp"
	"strconv"

	"github.com/adlio/trello"
)

// GetCardsInBoard returns a list of cards visible in the board
func (c *TrelloClient) GetCardsInBoard(id string) (cards []*trello.Card, err error) {
	if c == nil || c.CBoard == nil {
		return nil, fmt.Errorf("no board specified, get board first")

	}
	path := fmt.Sprintf("/board/%s/cards/visible", id)
	err = c.Client.Get(path, trello.Defaults(), &cards)
	if err != nil {
		return nil, err
	}
	logger.Debugln("Number of cards visible", len(cards))
	for _, card := range cards {
		card.SetClient(c.Client)
	}
	return cards, err
}

// FilterTasks gets tasks from a list of cards
func (c *TrelloClient) FilterTasks(cards []*trello.Card) (tasks []*Task, err error) {
	if c == nil || c.CBoard == nil {
		return nil, fmt.Errorf("no board specified, get board first")

	}
	for _, card := range cards {
		ok, hour := ValidateTaskName(card.Name) 
		if ok {
			task := &Task{
				Card: card,
				Hour: hour,
			}
			tasks = append(tasks, task)
		}
	}
	logger.Debugln("Number of tasks", len(tasks))
	return tasks, err
}

// StatisticTask gets tasks by members
func (c *TrelloClient) StatisticTask(tasks []*Task) (err error) {
	if c == nil || c.CBoard == nil {
		return fmt.Errorf("no board specified, get board first")

	}
	for _, task := range tasks {
		if task.Card.IDMembers != nil && len(task.Card.IDMembers) > 0 {
			for _, member := range task.Card.IDMembers {
				if ValidateMember(member) {
					stat, ok := c.MemberStatistics[member]
					if ok {
						stat.NTasks++
						stat.NHours = stat.NHours + task.Hour
						stat.TotalTasks = append(stat.TotalTasks, task)
						if task.Card.IDList == c.DoneList {
							stat.NDoneTasks++
							stat.NDoneHours = stat.NDoneHours + task.Hour
						}
						if !c.CheckCardInSkipList(task.Card) {
							stat.NProgressTasks ++
							stat.NProgressHours = stat.NProgressHours + task.Hour
						}
					}
				}
			}
		}
	}
	return err
}

// ValidateTasksName validates card name is task type or not
func ValidateTaskName(name string) (bool, int) {
	re := regexp.MustCompile(TASK_NAME_PATTERN)
	if !re.MatchString(name) {
		return false, 0
	}
	matches := re.FindStringSubmatch(name)
	fmt.Println(matches)
	if len(matches) < 3 {
		return true, 0
	}
	timeValue := matches[2]
	timeValueInt, err := strconv.Atoi(timeValue)
	if err!= nil {
        return true, 0
    }
	return true, timeValueInt
}

// PrintMemberStatistics prinses the member statistics
func (c *TrelloClient) PrintMemberStatistics() {
	for _, stat := range c.MemberStatistics {
		logger.Debugln("member: ", stat.Name, "\ttask [done/progress/total]: ", stat.NDoneTasks,"/", stat.NProgressTasks, "/", stat.NTasks,
		 "\tHour [done/progress/total]: ", stat.NDoneHours, "/", stat.NProgressHours, "/", stat.NHours)
	}
}

// CheckCardInSkipList returns true if card in the skip list
func (c *TrelloClient) CheckCardInSkipList(card *trello.Card) bool {
	for _, skipList := range c.SkipLists {
		if card.IDList == skipList {
			return true
		}
	}
	return false
}