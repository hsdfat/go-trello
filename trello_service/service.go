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
	instance := GetInstance()
	board, err := instance.Client.GetBoard(id)
	if err != nil {
		return err
	}
	logger.Debugln("Get board:", board.Name)
	instance.CBoard = board
	members, err := instance.GetMembersInBoard()
	if err != nil {
		return err
	}
	instance.Members = members
	logger.Debugln("Get members")
	for _, member := range members {
		logger.Debugln(member.Email, member.FullName)
		action, err := instance.GetMemberActions(member.ID, trello.Defaults())
		if err != nil {
			return err
		}
		logger.Debugln("Get actions")
		for _, action := range action {
			logger.Debugln("Action:", action.Type)
		}
	}

    // Get cards in board
	cards, err := instance.GetCardsInBoard(id)
	if err!= nil {
        return err
    }
	instance.Cards = cards
	logger.Debugln("Get cards", len(cards))
    instance.FilterTasks(instance.Cards)

	list, err := instance.GetLists()
	if err!= nil {
        return err
    }
	logger.Debugln("Get List", len(list))
	instance.StaticList()
	
	return nil
}
