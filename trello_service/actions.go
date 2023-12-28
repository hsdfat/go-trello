package trello_service

import (
	"fmt"
	"sync"
	"github.com/adlio/trello"
)

// GetActionsByCard get actions list by card but only ListChangeActions action type
func (c *TrelloClient) GetActionsByCard(task *Task, wgParent *sync.WaitGroup) (err error) {
	defer wgParent.Done()
	// Check board
	if c == nil || c.Board == nil {
		return fmt.Errorf("no board specified, cannot get actions")
	}
	// Check card
	if task.Card == nil {
		return fmt.Errorf("no card specified, cannot get actions")
	}

	actions, err := task.Card.GetActions(trello.Arguments{"filter": "updateCard:idList,updateCard:closed,createCard"})
	if err != nil {
		return err
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(actions))
	//logger.Debugln("len actions: ", len(actions))
	for _, action := range actions {
		go c.DailyTrackingStats.TrackingAction(task, action, wg)
	}
	wg.Wait()

	return
}
