package trello_service

import (
	"fmt"
	"go-trello/logger"
	"go-trello/utils"
	"strconv"
	"time"

	"github.com/adlio/trello"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

var c *TrelloClient // Using only one instance like singleton

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
		c.MemberStats = make(map[string]*MemberStats)
	}
	return c
}

// DeleteInstance deletes the instance of service
func DeleteInstance() {
	c = nil
}

// GetBoardInfo returns board information include board, members, actions of members
func GetBoardInfo(id string, startDay, endDay time.Time) *TrelloClient {
	instance := GetInstance()
	// instance.SetSprintStartDay(startDay)
	// instance.SetSprintEndDay(endDay)
	err := instance.SetSprintDuration(startDay, endDay)
	if err != nil {
		logger.Error(err)
	}
	board, err := instance.Client.GetBoard(id)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug("Get board:", board.Name)
	instance.Board = board

	list, err := instance.GetLists()
	if err != nil {
		logger.Error(err)
	}
	logger.Debug("Get List", len(list))
	instance.StatisticList()
	_, err = instance.GetMembersInBoard()
	if err != nil {
		logger.Error(err)
	}

	// Get cards in board
	// cards, err := instance.GetCardsInBoard(id)
	// if err != nil {
	// 	logger.Error(err)
	// }
	cards, err, number := instance.GetCardsInBoard(id)
	logger.Info("@#: ", number)

	tasks, err := instance.FilterTasks(cards)
	if err != nil {
		logger.Error(err)
	}
	instance.Tasks = tasks

	// Statistics members
	err = instance.StatisticTask(tasks)
	if err != nil {
		logger.Error(err)
	}
	instance.PrintMemberStatistics()

	// for memberId, _ := range instance.Members {
	// 	instance.DailyTrackingStats.PrintMemberStatTracking(memberId)
	// 	logger.Debug("______________________________", memberId)
	// 	instance.DailyTrackingStats.PrintLinkList()
	// }

	return instance
}

func ExportTotalMemberToCsv(memberData *TrelloClient) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet(utils.NameSMFTeam)
	if err != nil {
		logger.Error(err)
	}
	totalDoneTasks, totalProgressTasks, totalRemainingTasks, totalExtraTasks, totalTasks, totalDoneHours, totalProgressHours, totalExtraHours, totalHours := 0, 0, 0, 0, 0, 0, 0, 0, 0
	f.SetCellValue(utils.NameSMFTeam, "A1", "Name")
	f.SetCellValue(utils.NameSMFTeam, "B1", "Done Tasks")
	f.SetCellValue(utils.NameSMFTeam, "C1", "Progress Tasks")
	f.SetCellValue(utils.NameSMFTeam, "D1", "Remaining Tasks")
	f.SetCellValue(utils.NameSMFTeam, "E1", "Extra Tasks")
	f.SetCellValue(utils.NameSMFTeam, "F1", "Tasks")
	f.SetCellValue(utils.NameSMFTeam, "G1", "Done Hours")
	f.SetCellValue(utils.NameSMFTeam, "H1", "Progress Hours")
	f.SetCellValue(utils.NameSMFTeam, "I1", "Extra Hours")
	f.SetCellValue(utils.NameSMFTeam, "J1", "Hours")
	f.SetCellValue(utils.NameSMFTeam, "A13", "Total")
	i := 0
	for _, stat := range memberData.MemberStats {
		f.SetCellValue(utils.NameSMFTeam, "A"+strconv.Itoa((i+2)), stat.Name)
		f.SetCellValue(utils.NameSMFTeam, "B"+strconv.Itoa((i+2)), stat.NDoneTasks)
		totalDoneTasks += int(stat.NDoneTasks)
		f.SetCellValue(utils.NameSMFTeam, "C"+strconv.Itoa((i+2)), stat.NProgressTasks)
		totalProgressTasks += int(stat.NProgressTasks)
		f.SetCellValue(utils.NameSMFTeam, "D"+strconv.Itoa((i+2)), stat.NTasks-stat.NProgressTasks-stat.NDoneTasks)
		totalRemainingTasks += int(stat.NTasks - stat.NProgressTasks - stat.NDoneTasks)
		f.SetCellValue(utils.NameSMFTeam, "E"+strconv.Itoa((i+2)), stat.NExtraTasks)
		totalExtraTasks += int(stat.NExtraTasks)
		f.SetCellValue(utils.NameSMFTeam, "F"+strconv.Itoa((i+2)), stat.NTasks)
		totalTasks += int(stat.NTasks)
		f.SetCellValue(utils.NameSMFTeam, "G"+strconv.Itoa((i+2)), stat.NDoneHours)
		totalDoneHours += int(stat.NDoneHours)
		f.SetCellValue(utils.NameSMFTeam, "H"+strconv.Itoa((i+2)), stat.NProgressHours)
		totalProgressHours += int(stat.NProgressHours)
		f.SetCellValue(utils.NameSMFTeam, "I"+strconv.Itoa((i+2)), stat.NExtraHours)
		totalExtraHours += int(stat.NExtraHours)
		f.SetCellValue(utils.NameSMFTeam, "J"+strconv.Itoa((i+2)), stat.NHours)
		totalHours += int(stat.NHours)
		i += 1
	}

	//set total
	f.SetCellValue(utils.NameSMFTeam, "B"+strconv.Itoa((i+2)), totalDoneTasks)
	f.SetCellValue(utils.NameSMFTeam, "C"+strconv.Itoa((i+2)), totalProgressTasks)
	f.SetCellValue(utils.NameSMFTeam, "D"+strconv.Itoa((i+2)), totalRemainingTasks)
	f.SetCellValue(utils.NameSMFTeam, "E"+strconv.Itoa((i+2)), totalExtraTasks)
	f.SetCellValue(utils.NameSMFTeam, "F"+strconv.Itoa((i+2)), totalTasks)
	f.SetCellValue(utils.NameSMFTeam, "G"+strconv.Itoa((i+2)), totalDoneHours)
	f.SetCellValue(utils.NameSMFTeam, "H"+strconv.Itoa((i+2)), totalProgressHours)
	f.SetCellValue(utils.NameSMFTeam, "I"+strconv.Itoa((i+2)), totalExtraHours)
	f.SetCellValue(utils.NameSMFTeam, "J"+strconv.Itoa((i+2)), totalHours)

	//set size of coloum
	err_size_column := f.SetColWidth(utils.NameSMFTeam, "A", "J", 20)
	if err_size_column != nil {
		fmt.Println(err_size_column)
	}

	err_size_height := f.SetRowHeight(utils.NameSMFTeam, 1, 20)
	if err_size_height != nil {
		fmt.Println(err_size_height)
	}

	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs(utils.NameOfFile); err != nil {
		fmt.Println(err)
	}
	return nil
}

func ExportDataOfMembersToExcel(memberData *TrelloClient) {
	numberOfSprint := memberData.DailyTrackingStats.CountDaysInSprint()
	startDay := viper.GetString("trello.startDay")
	startDayOfSprint, err := time.Parse("02-01-2006", startDay)
	if err != nil {
		logger.Error("Cannot parse start day: ", err)
	}
	startDayOfSprintInVn := utils.TimeLocal(startDayOfSprint)
	numberOfDayToCurrentDay := memberData.DailyTrackingStats.CountNumberToCurrentDayNew(startDayOfSprintInVn) // number of days from start day to current day

	for memberId, _ := range memberData.Members {
		totalTasks := memberData.MemberStats[memberId].NTasks
		totalHours := memberData.MemberStats[memberId].NHours
		memberData.DailyTrackingStats.ExportDataOfEachMemberToExcel(memberId, totalTasks, numberOfSprint, totalHours, numberOfDayToCurrentDay)
		DrawLineChart(memberData.MemberStats[memberId].Name)
	}
}

func ExportDataOfDailyToExcel(memberData *TrelloClient) {
	numberOfMembers := len(memberData.Members)
	numberOfSprint := memberData.DailyTrackingStats.CountDaysInSprint()
	memberData.DailyTrackingStats.PrintLinkList()

	//initTotalTime := 8 * numberOfMembers * numberOfSprint
	initTotalTime := viper.GetInt("smfHour.smf")

	startDay := viper.GetString("trello.startDay")
	startDayOfSprint, err := time.Parse("02-01-2006", startDay)
	if err != nil {
		logger.Error("Cannot parse start day: ", err)
	}
	startDayOfSprintInVn := utils.TimeLocal(startDayOfSprint)
	numberOfDayToCurrentDay := memberData.DailyTrackingStats.CountNumberToCurrentDayNew(startDayOfSprintInVn) // number of days from start day to current day
	logger.Info("CountNumberToCurrentDay2: ", memberData.DailyTrackingStats.CountNumberToCurrentDayNew(startDayOfSprintInVn))
	logger.Info("Number of: ", numberOfDayToCurrentDay)
	dataDailyList := memberData.DailyTrackingStats.calculateRemainingTasksDailyList(numberOfMembers, initTotalTime)
	//SetCellValue return values of excel
	SetCellValue(utils.MemberActionDaily, dataDailyList, int(memberData.DailyTrackingStats.head.stat.NTasks), numberOfSprint, numberOfDayToCurrentDay) // total tasts here
}

func SortMembersActionsDailyUseName(memberActions []*MemberActions) {
	length := len(memberActions)
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-i-1; j++ {
			if memberActions[j].NameOfMember > memberActions[j+1].NameOfMember {
				memberActions[j], memberActions[j+1] = memberActions[j+1], memberActions[j]
			}
		}
	}
}

// SortTasksUseTypeOfTask returns list sorted of type of tasks, ex: Tool hieu nang
func SortTasksUseTypeOfTask(tasks []*Task) {
	length := len(tasks)
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-i-1; j++ {
			if tasks[j].TypeOfTask > tasks[j+1].TypeOfTask {
				tasks[j], tasks[j+1] = tasks[j+1], tasks[j]
			}
		}
	}
}

func SortTasksUseCardName(tasks []*Task) {
	length := len(tasks)
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-i-1; j++ {
			if tasks[j].Card.Name > tasks[j+1].Card.Name {
				tasks[j], tasks[j+1] = tasks[j+1], tasks[j]
			}
		}
	}
}

func SortMembersActionsDailyUseTypeOfTask(tasks []*Task) {
	length := len(tasks)
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-i-1; j++ {
			if tasks[j].TypeOfTask > tasks[j+1].TypeOfTask {
				tasks[j], tasks[j+1] = tasks[j+1], tasks[j]
			}
		}
	}
}

func SortMembersActionsDailyUseTime(memberActions []*MemberActions) {
	length := len(memberActions)
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-i-1; j++ {
			if memberActions[j].Time.After(memberActions[j+1].Time) {
				memberActions[j], memberActions[j+1] = memberActions[j+1], memberActions[j]
			}
		}
	}
}

func GetStatusOfTaskInGroupSheet(tasks *Task) string {
	if tasks == nil {
		logger.Error("Nil task")
	}
	if tasks.IsDone {
		return "Done"
	} else if tasks.IsInProgress {
		return "Inprogress"
	} else {
		return "Sprint Backlog or Pending"
	}
}
