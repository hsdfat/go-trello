package trello_service

import (
	"fmt"
	"go-trello/logger"
	"regexp"

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
func (c *TrelloClient) FilterTasks(cards []*trello.Card) (tasks []*trello.Card, err error) {
	if c == nil || c.CBoard == nil {
		return nil, fmt.Errorf("no board specified, get board first")

	}
	for _, card := range cards {
		if ValidateTaskName(card.Name) {
			tasks = append(tasks, card)
		}
	}
	logger.Debugln("Number of tasks", len(tasks))
	return tasks, err
}

// FilterMemberTasks gets tasks by members
func (c *TrelloClient) FilterMemberTasks(tasks []*trello.Card) (err error) {
	if c == nil || c.CBoard == nil {
        return fmt.Errorf("no board specified, get board first")

    }
    for _, task := range tasks {
        if task.Members != nil && len(task.Members) > 0 {
			for _, member := range task.Members {
                if ValidateMember(member.ID) {
					stat, ok := c.MemberStatistics[member.ID]
					if ok {
						stat.NTask++
											
					}
                }
            }
		}
    }
    return err
}

// ValidateTasksName validates card name is task type or not
func ValidateTaskName(name string) bool {
	re := regexp.MustCompile(TASK_NAME_PATTERN)
	return re.MatchString(name)
}
