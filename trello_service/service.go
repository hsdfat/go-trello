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
		// Init instance
		c.Members = make(map[string]*trello.Member)
		c.Cards = make(map[string]*trello.Card)
		c.Labels = make(map[string]*trello.Label)
		c.Actions = make(map[string]*trello.Action)
		c.Lists = make(map[string]*trello.List)
		c.Caretory = make(map[string]string)
		c.MemberStatistics = make(map[string]*MemberStatistics)
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
	instance.Board = board

	list, err := instance.GetLists()
	if err != nil {
		return err
	}
	logger.Debugln("Get List", len(list))
	instance.StatisticList()

	_, err = instance.GetMembersInBoard()
	if err != nil {
		return err
	}

	logger.Debugln("Get members")

	// Get cards in board
	cards, err := instance.GetCardsInBoard(id)
	if err != nil {
		return err
	}
	logger.Debugln("Get cards", len(cards))
	tasks, err := instance.FilterTasks(cards)
	if err != nil {
		logger.Errorln(err)
	}

	// Statistics members
	err = instance.StatisticTask(tasks)
	if err != nil {
		logger.Errorln(err)
	}
	instance.PrintMemberStatistics()
	return nil
}
