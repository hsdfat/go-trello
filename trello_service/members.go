package trello_service

import (
	"fmt"

	"github.com/adlio/trello"
)

// GetMembersInBoard returns a list of members of the board
func (c *TrelloClient) GetMembersInBoard() (members []*trello.Member, err error) {
	if c.CBoard == nil {
		return nil, fmt.Errorf("no board specified, get board first")
	}
	c.Members, err = c.CBoard.GetMembers()
	if err != nil {
		return nil, err
	}
	return c.Members, nil
}

// GetMemberCard returns a list of actions of the member
func (c *TrelloClient) GetMemberActions(id string) (action []*trello.Action, err error) {
	if c.CBoard == nil {
        return nil, fmt.Errorf("no board specified, get board first")
    }
    
	// check member already exists in board
	if c.Members == nil || len(c.Members) == 0 {
		return nil, fmt.Errorf("no members specified, get board first")
	}

	for _, member := range c.Members {
		if member.ID == id {
			action := []*trello.Action{}
            err = c.Client.Get(fmt.Sprintf("member/%v/actions", member.ID), trello.Defaults(), &action)
            if err!= nil {
                return nil, err
            }

            return action, nil
        }
	}
	return nil, fmt.Errorf("no actions specified, get board first")
}

