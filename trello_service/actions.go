package trello_service

import (
	"fmt"

	"github.com/adlio/trello"
)

// GetActionsByCard get actions list by card but only UpdateCard action type
func (c *TrelloClient) GetActionsByCard() (actions []*trello.Action, err error) {
	// Check board
	if c == nil || c.CBoard == nil {
		return nil, fmt.Errorf("no board specified, cannot get actions")
	}
	// Check card
	if c.Cards == nil || len(c.Cards) == 0 {
		return nil, fmt.Errorf("no card specified, cannot get actions")
	}

	for _, card := range c.Cards {
		card.GetListChangeActions()
	}

	return
}
