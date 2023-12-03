package trello_service

import (
	"fmt"
	"go-trello/logger"
	"strconv"

	"github.com/adlio/trello"
	"github.com/spf13/viper"
	// excelize "github.com/xuri/excelize/v2"
	"github.com/xuri/excelize/v2"
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
func GetBoardInfo(id string) *TrelloClient {
	instance := GetInstance()
	board, err := instance.Client.GetBoard(id)
	if err != nil {
		logger.Errorln("1")
	}
	logger.Debugln("Get board:", board.Name)
	instance.Board = board

	list, err := instance.GetLists()
	if err != nil {
		logger.Errorln("2")
	}
	logger.Debugln("Get List", len(list))
	instance.StatisticList()

	_, err = instance.GetMembersInBoard()
	if err != nil {
		logger.Errorln("3")
	}

	logger.Debugln("Get members")

	// Get cards in board
	cards, err := instance.GetCardsInBoard(id)
	if err != nil {
		logger.Errorln("4")
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
	return instance
}

func Export_csv(id string) error {
	memberData := GetBoardInfo(id)
	//print
	fmt.Println(memberData.MemberStatistics)
	memberData.PrintMemberStatistics()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("SMF")
	if err != nil {
		fmt.Println(err)
		logger.Errorln("5")
	}

	f.SetCellValue("SMF", "A1", "Name")
	f.SetCellValue("SMF", "B1", "NDoneTasks")
	f.SetCellValue("SMF", "C1", "NProgressTasks")
	f.SetCellValue("SMF", "D1", "NTasks")
	f.SetCellValue("SMF", "E1", "NDoneHours")
	f.SetCellValue("SMF", "F1", "NProgressHours")
	f.SetCellValue("SMF", "G1", "NHours")
	i := 0
	for _, stat := range memberData.MemberStatistics {
		f.SetCellValue("SMF", "A"+strconv.Itoa((i+2)), stat.Name)
		f.SetCellValue("SMF", "B"+strconv.Itoa((i+2)), stat.NDoneTasks)
		f.SetCellValue("SMF", "C"+strconv.Itoa((i+2)), stat.NProgressTasks)
		f.SetCellValue("SMF", "D"+strconv.Itoa((i+2)), stat.NTasks)
		f.SetCellValue("SMF", "E"+strconv.Itoa((i+2)), stat.NDoneHours)
		f.SetCellValue("SMF", "F"+strconv.Itoa((i+2)), stat.NProgressHours)
		f.SetCellValue("SMF", "G"+strconv.Itoa((i+2)), stat.NHours)
		i += 1
	}
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
	return nil
}
