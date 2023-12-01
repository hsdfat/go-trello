package trello_service

import (
	"fmt"
	"go-trello/logger"
	"strings"

	"github.com/adlio/trello"
)

// GetLists returns a list of boards
func (c *TrelloClient) GetLists() (lists []*trello.List, err error) {
	// Check board existence
	if c == nil || c.CBoard == nil {
		return nil, fmt.Errorf("no board specified, cannot get lists")
	}

	lists, err = c.CBoard.GetLists(trello.Defaults())

	if err != nil {
		return nil, err
	}

	for _, list := range lists {
		list.SetClient(c.Client)
	}
	c.Lists = lists
	return
}

// StatisticList returns done list id and skip lists
func (c *TrelloClient) StatisticList() (doneList string, skipList []string, err error) {
	// Check board existence
	if c == nil || c.CBoard == nil {
		return "", nil, fmt.Errorf("no board specified, cannot get lists")
	}
	logger.Debugln("Done list", doneList)

	lists := c.Lists
	if len(lists) == 0 {
		return "", nil, fmt.Errorf("no lists found")
	}

	for _, list := range lists {
		if strings.Contains(strings.ToLower(list.Name), DONE_LIST) {
			doneList = list.ID
		}
		if strings.Contains(strings.ToLower(list.Name), SPRINT_BACKLOG_LIST) || strings.Contains(strings.ToLower(list.Name), EPIC_LIST) {
			skipList = append(skipList, list.ID)
		}
	}
	c.DoneList = doneList
	c.SkipLists = skipList
	return
}