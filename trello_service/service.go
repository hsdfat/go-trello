package trello_service

import (
	"go-trello/logger"

	"github.com/adlio/trello"
	"github.com/spf13/viper"
)

var c *TrelloClient

// GetInstance returns singleton trello client instance
func GetInstance() *TrelloClient {
	if c == nil {
		c = &TrelloClient{}
		c.Client = trello.NewClient(
			viper.GetString("trello.apiKey"),
			viper.GetString("trello.token"),
		)
	}
	return c
}

// GetBoardInfo returns board information include board, members, actions of members
func GetBoardInfo(id string) error {
	board, err := GetInstance().Client.GetBoard(id)
	if err != nil {
		return err
	}
	logger.Debugln("Get board:", board.Name)
	GetInstance().CBoard = board
	members, err := GetInstance().GetMembersInBoard()
	if err != nil {
		return err
	}
	GetInstance().Members = members
	logger.Debugln("Get members")
	for _, member := range members {
		logger.Debugln(member.Email, member.FullName)
		action, err := GetInstance().GetMemberActions(member.ID)
		if err != nil {
			return err
		}
		logger.Debugln("Get actions")
		for _, action := range action {
			logger.Debugln("Action:", action.Type)
		}
	}

	return nil
}
