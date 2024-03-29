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
func (c *TrelloClient) GetCardsInBoard(id string) (cards []*trello.Card, err error, number int) {
	if c == nil || c.Board == nil {
		return nil, fmt.Errorf("no board specified, get board first"), 0
	}
	path := fmt.Sprintf("/board/%s/cards/visible", id)
	err = c.Client.Get(path, trello.Defaults(), &cards)
	if err != nil {
		return nil, err, 0
	}
	// logger.Debugln("Number of cards visible", len(cards))
	i := 0
	for _, card := range cards {
		card.SetClient(c.Client)
		c.Cards[card.ID] = card
		i += 1
	}
	number = i
	return cards, err, number
}

// FilterTasks gets tasks from a list of cards
func (c *TrelloClient) FilterTasks(cards []*trello.Card) (tasks []*Task, err error) {
	logger.Debugln("Filtering tasks")
	if c == nil || c.Board == nil {
		return nil, fmt.Errorf("no board specified, get board first")
	}
	for _, card := range cards {
		ok, hour, isExtraTask, typeOfTask := ValidateTaskName(card.Name)
		// logger.Debug("ok: ", ok)
		// logger.Debug("hour: ", hour)
		memberMap := c.Members
		var idMember string
		for _, idMembers := range card.IDMembers {
			idMember = idMembers
		}
		member := memberMap[idMember]
		if ok {
			task := &Task{
				Card:         card,
				Hour:         hour,
				IsDone:       card.IDList == c.DoneList,
				IsInProgress: c.CheckTaskInProgress(card),
				IsExtra:      isExtraTask,
				TypeOfTask:   typeOfTask, 		// ex: Test hieu nang
				Members:      member,
			}
			// logger.Debug("is Extra Tast: ", isExtraTask)
			creationTime, err := GetCreationTime(card.ID)

			if err == nil {
				task.CreationTime = creationTime
			}
			tasks = append(tasks, task)
		}
	}
	// logger.Debug("Number of tasks", len(tasks))
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
							if task.IsExtra {
								atomic.AddInt32(&stat.NExtraTasks, 1)
								atomic.AddInt32(&stat.NExtraHours, task.Hour)
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
func ValidateTaskName(name string) (bool, int32, bool, string) {
	re := regexp.MustCompile(TASK_NAME_PATTERN)
	if !re.MatchString(name) {
		return false, 0, false, ""
	}
	matches := re.FindStringSubmatch(name)
	extraTask := matches[1]
	typeOfTask := matches[3]
	if (len(matches) < 3) && (extraTask != "Ngoài") {
		return true, 0, false, typeOfTask
	}
	// logger.Info("---------------------------")
	// logger.Info("matches[0]: ", matches[0])		// content of task
	// logger.Info("matches[1]: ", matches[1])		
	// logger.Info("matches[2]: ", matches[2])		// time estimate for task
	// logger.Info("matches[3]: ", matches[3])		// type of task
	timeValue := matches[2]
	timeValueInt, err := strconv.Atoi(timeValue)
	if err != nil {
		return true, 0, false, typeOfTask
	}
	isExtraTask := false
	// logger.Debug("value of extraTask: ", extraTask)
	if extraTask == "Ngoài" {
		isExtraTask = true
		timeValueInt, err = strconv.Atoi(timeValue)
		if err != nil {
			return false, 0, isExtraTask, typeOfTask
		}
	}
	// GroupName
	return true, int32(timeValueInt), isExtraTask, typeOfTask
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

// func CheckCardInExtraList(card *trello.Card, extraLists []string) bool {
// 	for _, extraList := range extraLists {
// 		if card.IDList == extraList {
//             return true
//         }
// 	}
// 	return false
// }

// func (c *TrelloClient) CheckTaskExtra(card *trello.Card) bool {
// 	return !CheckCardInExtraList(card, c.ExtraLists) && card.IDList != c.DoneList
// }
