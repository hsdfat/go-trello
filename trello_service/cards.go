package trello_service

import (
	"fmt"
	"go-trello/logger"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/adlio/trello"
)

// GetCardsInBoard returns a list of cards visible in the board
func (c *TrelloClient) GetCardsInBoard(id string) (cards []*trello.Card, err error) {
	if c == nil || c.Board == nil {
		return nil, fmt.Errorf("no board specified, get board first")

	}
	path := fmt.Sprintf("/board/%s/cards/visible", id)
	err = c.Client.Get(path, trello.Defaults(), &cards)
	if err != nil {
		return nil, err
	}
	// logger.Debugln("Number of cards visible", len(cards))
	for _, card := range cards {
		card.SetClient(c.Client)
		c.Cards[card.ID] = card
	}
	return cards, err
}

// FilterTasks gets tasks from a list of cards
func (c *TrelloClient) FilterTasks(cards []*trello.Card) (tasks []*Task, err error) {
	logger.Debugln("Filtering tasks")
	if c == nil || c.Board == nil {
		return nil, fmt.Errorf("no board specified, get board first")

	}
	for _, card := range cards {
		ok, hour := ValidateTaskName(card.Name)
		if ok {

			task := &Task{
				Card:         card,
				Hour:         hour,
				IsDone:       card.IDList == c.DoneList,
				IsInProgress: c.CheckTaskInProgress(card),
			}
			creationTime, err := GetCreationTime(card.ID)
			if err == nil {
				task.CreationTime = creationTime
			}
			tasks = append(tasks, task)
		}
	}
	logger.Debugln("Number of tasks", len(tasks))
	return tasks, err
}

// StatisticTask gets tasks by members
func (c *TrelloClient) StatisticTask(tasks []*Task) (err error) {
	logger.Debugln("Statistic Tasks")
	if c == nil || c.Board == nil {
		return fmt.Errorf("no board specified, get board first")
	}
	wg := new(sync.WaitGroup)
	wg.Add(3 * len(tasks))
	for _, task := range tasks {
		// TODO: Calculate for total done tasks and progress tasks
		// Member statistics
		go func(task *Task, wg *sync.WaitGroup) {
			defer wg.Done()
			if task.Card.IDMembers != nil && len(task.Card.IDMembers) > 0 {
				for _, member := range task.Card.IDMembers {
					if ValidateMember(member) {
						stat, ok := c.MemberStats[member]
						if ok {
							atomic.AddInt32(&stat.NTasks, 1)
							atomic.AddInt32(&stat.NHours, task.Hour)
							stat.TotalTasks = append(stat.TotalTasks, task)
							if task.IsDone {
								atomic.AddInt32(&stat.NDoneTasks, 1)
								atomic.AddInt32(&stat.NDoneHours, task.Hour)
							}
							if task.IsInProgress {
								atomic.AddInt32(&stat.NProgressTasks, 1)
								atomic.AddInt32(&stat.NProgressHours, task.Hour)
							}
						}
					}
				}
			}
		}(task, wg)

		// Daily Statistics
		go c.DailyTrackingStats.TrackingTaskCreationByDate(task, wg)
		//go c.DailyTrackingStats.TrackingActionByDate(task, wg)
		go c.GetActionsByCard(task, wg)
	}
	wg.Wait()
	return err
}

// ValidateTasksName validates card name is task type or not
func ValidateTaskName(name string) (bool, int32) {
	re := regexp.MustCompile(TASK_NAME_PATTERN)
	if !re.MatchString(name) {
		return false, 0
	}
	matches := re.FindStringSubmatch(name)
	if len(matches) < 3 {
		return true, 0
	}
	timeValue := matches[2]
	timeValueInt, err := strconv.Atoi(timeValue)
	if err != nil {
		return true, 0
	}
	return true, int32(timeValueInt)
}

// CheckCardInSkipList returns true if card in the skip list
func CheckCardInSkipList(card *trello.Card, skipLists []string) bool {
	for _, skipList := range skipLists {
		if card.IDList == skipList {
			return true
		}
	}
	return false
}

// CheckTaskInProgress checks task in progress or not by check card in skip list or done list
func (c *TrelloClient) CheckTaskInProgress(card *trello.Card) bool {

	return !CheckCardInSkipList(card, c.SkipLists) && card.IDList != c.DoneList
}
