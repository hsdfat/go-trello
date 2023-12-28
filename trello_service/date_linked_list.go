package trello_service

import (
	"fmt"
	"go-trello/logger"
	"go-trello/utils"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"

	"github.com/adlio/trello"
)

type DateLinkedList struct {
	head *DateNode
}

type DateNode struct {
	stat *DailyTrackingStats
	next *DateNode
}

func (list *DateLinkedList) CreateNode(stat *DailyTrackingStats) (node *DateNode) {
	return &DateNode{
		stat: stat,
		next: nil,
	}
}

func (list *DateLinkedList) AddNodeAtEnd(stat *DailyTrackingStats) {
	newNode := list.CreateNode(stat)
	if list.head == nil {
		list.head = newNode
		return
	}

	current := list.head
	for current.next != nil {
		current = current.next
	}
	current.next = newNode
}

func (list *DateLinkedList) PrintLinkList() {
	if list.head == nil {
		logger.Debugln("List is empty")
		return
	}
	current := list.head
	for current != nil {
		stat := current.stat
		if stat == nil {
			current = current.next
			continue
		}
		logger.Info(fmt.Sprintf("date [%s]: new task (done/progress/total): %d/%d/%d, new hour (done/progress/total): %d/%d/%d\t",
			stat.Date.Format("02-01-2006"), stat.NDoneTasks, stat.NProgressTasks, stat.NTasks, stat.NDoneHours, stat.NProgressHours, stat.NHours))
		current = current.next
	}
}

// calculate remaining hours of each day for SMF team
func (list *DateLinkedList) calculateRemainingTasksDailyList(numberOfMembers int, linear_hours int) []string {
	var remainingTasks int32 = 0
	var remainingHours int32 = 0
	remainingTasksData := []string{}
	if list.head == nil {
		logger.Debugln("List is empty")
		return nil
	}
	current := list.head
	for current != nil {
		stat := current.stat
		if stat == nil {
			current = current.next
			continue
		}
		remainingTasks += stat.NTasks - stat.NDoneTasks
		//time need to do remaining tasks
		remainingHours += stat.NHours - stat.NDoneHours
		remainingTasksData = append(remainingTasksData, stat.Date.Format("02-01-2006"))
		remainingTasksData = append(remainingTasksData, strconv.Itoa(int(remainingTasks)))
		remainingTasksData = append(remainingTasksData, strconv.Itoa(int(remainingHours)))
		remainingTasksData = append(remainingTasksData, strconv.Itoa(linear_hours))
		linear_hours -= 8 * numberOfMembers
		current = current.next
	}
	return remainingTasksData
}

func (list *DateLinkedList) PrintMemberStatTracking(id string) {
	if list.head == nil {
		logger.Debugln("List is empty")
		return
	}
	current := list.head
	for current != nil {
		stat := current.stat
		if stat == nil {
			current = current.next
			continue
		}
		memberStat, ok := stat.MemberStats[id]
		if !ok {
			current = current.next
			continue
		}

		logger.Info(fmt.Sprintf("member: [%s] date [%s]: new task (done/progress/extra task): %d/%d/%d, new hour (done/progress/extra task): %d/%d/%d\t", memberStat.FullName,
			stat.Date.Format("02-01-2006"), memberStat.NDoneTasks, memberStat.NProgressTasks, memberStat.NTasks, memberStat.NDoneHours, memberStat.NProgressHours, memberStat.NHours))
		current = current.next
	}
}

func (list *DateLinkedList) CountDaysInSprint() int {
	count := 0
	temp := list.head
	for temp != nil {
		temp = temp.next
		count += 1
	}
	return count
}

func (list *DateLinkedList) CountNumberToCurrentDayNew(starDayOfSprintInVn time.Time) int {
	count := 1
	temp := list.head
	var afterStartDay2 time.Time = starDayOfSprintInVn
	currentTimeInVn := utils.TimeLocal(time.Now())
	for temp != nil {
		afterStartDay2 = temp.stat.Date
		//afterStartDay = afterStartDay.AddDate(0, 0, 1)
		if afterStartDay2.Weekday() == 6 || afterStartDay2.Weekday() == 7 {
			continue
		}
		if utils.IsDateEqual(&afterStartDay2, &currentTimeInVn) {
			logger.Info("count:", count)
			return count
		}
		count += 1
		temp = temp.next
	}
	return count
}

func (list *DateLinkedList) CountNumberToCurrentDay(starDayOfSprintInVn time.Time) int {
	count := 1
	temp := list.head
	var afterStartDay time.Time = starDayOfSprintInVn
	currentTimeInVn := utils.TimeLocal(time.Now())
	for temp != nil {
		afterStartDay = afterStartDay.AddDate(0, 0, 1)
		if afterStartDay.Weekday() == 6 || afterStartDay.Weekday() == 7 {
			continue
		}
		if utils.IsDateEqual(&afterStartDay, &currentTimeInVn) {
			logger.Info("count:", count)
			return count
		}
		count += 1
		temp = temp.next
	}
	return count
}

func (list *DateLinkedList) ExportDataOfEachMemberToExcel(id string, totalTask int32, numberOfSprint int, totalHours int32, numberOfDayToCurrentDay int) {
	numberOfTasksNeedDone := totalTask
	numberOfRemainingHours := totalHours
	number := 1 // number to check with numberOfDaysToCurrentDay

	//export to excel
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Error(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Error(err)
		}
	}()

	if list.head == nil {
		logger.Debug("List is empty")
		return
	}

	current := list.head
	currentReal := list.head
	var i int = 65
	var j int = 65
	var countDay int = 1

	//get name and expected line to sheet of each member
	for current != nil {
		stat := current.stat
		if stat == nil {
			current = current.next
			continue
		}
		memberStat, ok := stat.MemberStats[id]

		if !ok {
			current = current.next
			continue
		}
		// Convert Name of Members to Company name
		ConvertNameOfMembersInLinkedList(memberStat)
		index, err := f.NewSheet(memberStat.Name)
		if err != nil {
			logger.Errorln(err)
		}
		f.SetCellValue(memberStat.Name, "A1", "Date")
		f.SetCellValue(memberStat.Name, "A2", "Tasks")
		f.SetCellValue(memberStat.Name, "A3", "Expected")
		f.SetCellValue(memberStat.Name, "A4", "Remaining Hours")
		f.SetCellValue(memberStat.Name, "A5", "Hours")
		f.SetCellValue(memberStat.Name, "B1", "StartDay")
		f.SetCellValue(memberStat.Name, "B2", int(totalTask))
		f.SetCellValue(memberStat.Name, "B3", int(totalTask))
		expected_task := utils.RoundFloat(utils.GetYValue(-float64(totalTask)/float64(numberOfSprint), countDay, totalTask), 2)
		f.SetCellValue(memberStat.Name, string((i+2))+"3", expected_task) //Expected task row
		f.SetActiveSheet(index)
		i += 1
		countDay += 1
		current = current.next
	}

	for currentReal != nil { //pointer to check real time
		stat := currentReal.stat
		if stat == nil {
			currentReal = currentReal.next
			continue
		}
		memberStat, ok := stat.MemberStats[id]
		//Get hour of each member
		remainingHourLinear := utils.FindHourOfEachMember(memberStat.Name)
		remainingHourLinear = remainingHourLinear - 8*(number-1)

		if !ok {
			currentReal = currentReal.next
			continue
		}
		numberOfTasksNeedDone = numberOfTasksNeedDone - memberStat.NDoneTasks
		numberOfRemainingHours = numberOfRemainingHours - memberStat.NDoneHours
		date := fmt.Sprintf("%s", stat.Date.Format("02-01-2006"))
		f.SetCellValue(memberStat.Name, string((j+2))+"1", date)                  //date
		f.SetCellValue(memberStat.Name, string((j+2))+"2", numberOfTasksNeedDone) //remaining tasks
		f.SetCellValue(memberStat.Name, string((j+2))+"4", numberOfRemainingHours)
		f.SetCellValue(memberStat.Name, string((j+2))+"5", remainingHourLinear)
		//remainingHourLinear -= 8

		//set size of coloum
		err_size_column := f.SetColWidth(memberStat.Name, "A", "L", 15)
		if err_size_column != nil {
			fmt.Println(err_size_column)
		}

		err_size_height := f.SetRowHeight(memberStat.Name, 1, 20)
		if err_size_height != nil {
			fmt.Println(err_size_height)
		}

		// f.SetActiveSheet(index)
		// countDay += 1
		j += 1
		if number == numberOfDayToCurrentDay {
			break
		}
		number += 1
		currentReal = currentReal.next
	}

	if err := f.SaveAs(utils.NameOfFile); err != nil {
		logger.Error(err)
	}
}

func (list *DateLinkedList) TrackingTaskCreationByDate(task *Task, wg *sync.WaitGroup) {
	defer wg.Done()
	if list.head == nil {
		logger.Debugln("List is empty")
		return
	}
	if task.CreationTime == nil {
		logger.Debugln("Task is not has creation time")
		return
	}
	current := list.head
	for current != nil {
		stat := current.stat
		if endOfDay(stat.Date).After(*task.CreationTime) {
			atomic.AddInt32(&stat.NTasks, 1)
			atomic.AddInt32(&stat.NHours, task.Hour)
			return
		}
		current = current.next
	}
}

func (list *DateLinkedList) PrintMemberActions() {
	if list.head == nil {
		logger.Debug("List is empty")
		return
	}
	current := list.head
	for current != nil {
		stat := current.stat
		for key, value := range stat.MemberActions {
			logger.Debug("-------------------------------------------------------")
			logger.Debug("Key: ", key)
			logger.Debug("Time: ", value.Time)
			logger.Debug("Before: ", value.ListBefore)
			logger.Debug("After: ", value.ListAfter)
			logger.Debug("Name: ", value.NameOfMember)
			logger.Debug("Task: ", value.ContentOfTask)
			logger.Debug("Action types: ", value.ActionTypes)
		}
		current = current.next
	}
}

func (list *DateLinkedList) GetMemberActionsDaily() []*MemberActions {
	var memberActions []*MemberActions
	today := utils.TimeLocal(time.Now())

	yesterday := today.AddDate(0, 0, -1)

	weekday := utils.TimeLocal(time.Now()).Weekday()
	// if is Monday
	if weekday == 1 {
		yesterday = today.AddDate(0, 0, -3)
	}

	// if yesterday is skipdate
	skipDate := viper.GetStringSlice("trello.skipDays")
	for utils.InSkipDays(skipDate, yesterday) {
		yesterday = yesterday.AddDate(0, 0, -1)
	}

	if list.head == nil {
		logger.Debug("List is empty")
	}
	current := list.head
	for current != nil {
		stat := current.stat
		for _, infoAction := range stat.MemberActions {
			infoAction.Time = utils.TimeLocal(infoAction.Time) //reset value of infoAction.Time to local time
			if utils.IsDateEqual(&infoAction.Time, &yesterday) {
				memberActions = append(memberActions, infoAction)
			}
		}
		current = current.next
	}
	return memberActions
}

func (list *DateLinkedList) GetMemberActionsSprint() []*MemberActions {
	var memberActions []*MemberActions
	if list.head == nil {
		logger.Debug("List is empty")
	}
	current := list.head
	for current != nil {
		stat := current.stat
		for _, infoAction := range stat.MemberActions {
			memberActions = append(memberActions, infoAction)
		}
		current = current.next
	}
	return memberActions
}

func (list *DateLinkedList) ExportMemberActionsDailyToExcel() {
	memberActions := list.GetMemberActionsDaily()
	if memberActions == nil {
		logger.Info("Not actions in this day: ", utils.TimeLocal(time.Now()))
	}
	SortMembersActionsDailyUseName(memberActions)
	SortMembersActionsDailyUseTime(memberActions)
	SetMemberActionsDaily(utils.MemberActionDaily, memberActions)
}

func (list *DateLinkedList) ExportMemberActionsSprintToExcel() {
	memberActions := list.GetMemberActionsSprint()
	if memberActions == nil {
		logger.Info("Not actions in this sprint: ", time.Now())
	}
	SortMembersActionsDailyUseName(memberActions)
	SortMembersActionsDailyUseTime(memberActions)
	SetMemberActionsSprint(utils.NameSMFTeam, memberActions)
}

// ExportGroupActionsSprintToExcel
func (list *DateLinkedList) ExportGroupActionsSprintToExcel(tasks []*Task) {
	// memberActions := list.GetMemberActionsSprint()
	// if memberActions == nil {
	// 	logger.Info("Not actions in this sprint: ", time.Now())
	// }
	SortTasksUseTypeOfTask(tasks)
	//SortTasksUseCardName(tasks)
	SetGroupActionsSprint(utils.Group, tasks)
}

func (list *DateLinkedList) TrackingAction(task *Task, action *trello.Action, wg *sync.WaitGroup) {
	defer wg.Done()

	ins := GetInstance()

	//add skip date list

	//skipDate := viper.GetStringSlice("trello.skipDays")
	// for i := 0; i < len(skipDate); i++ {
	// 	skipDateTime, err := time.Parse("02-01-2006", skipDate[i])
	// 	if err != nil {
	// 		logger.Error("Cannot parse start day: ", err)
	// 	}
	// 	//skipDateTimeVn := utils.TimeLocal(skipDateTime)
	// 	skipDateTimeVn := skipDateTime
	// 	logger.Info("skipDateTimeVn: ", skipDateTimeVn)
	// }

	if list.head == nil {
		logger.Debugln("List is empty")
		return
	}
	current := list.head
	for current != nil {
		stat := current.stat
		if endOfDay(stat.Date).After(action.Date) {
			logger.Info("action.Date need check: ", utils.TimeLocal(action.Date))
			// if action.Data != nil && slices.Contains(skipDate, utils.TimeLocal(action.Date).String()) {
			//if action.Data != nil && !utils.InSkipDays(skipDate, utils.TimeLocal(action.Date)) {
			if action.Data != nil {
				var taskDone, taskUndone, taskInProgress, taskNotInProgress bool

				if action.Type == "createCard" {
					if action.Data.List.ID == ins.DoneList {
						taskDone = true
						atomic.AddInt32(&stat.NDoneTasks, 1)
						atomic.AddInt32(&stat.NDoneHours, task.Hour)
					}
				}
				logger.Info("Action Date: ", action.Date)
				if action.Data.ListAfter != nil {
					if action.Data.ListAfter.ID == ins.DoneList {
						taskDone = true
						atomic.AddInt32(&stat.NDoneTasks, 1)
						atomic.AddInt32(&stat.NDoneHours, task.Hour)
					} else if !CheckListContains(ins.SkipLists, action.Data.ListAfter.ID) {
						taskInProgress = true
						atomic.AddInt32(&stat.NProgressTasks, 1)
						atomic.AddInt32(&stat.NProgressHours, task.Hour)

					}
				}
				if action.Data.ListBefore != nil {
					if action.Data.ListBefore.ID == ins.DoneList {
						taskUndone = true
						atomic.AddInt32(&stat.NDoneTasks, -1)
						atomic.AddInt32(&stat.NDoneHours, -1*task.Hour)
					} else if !CheckListContains(ins.SkipLists, action.Data.ListBefore.ID) {
						if !taskInProgress {
							taskNotInProgress = true
						}
						atomic.AddInt32(&stat.NProgressTasks, -1)
						atomic.AddInt32(&stat.NProgressHours, -1*task.Hour)
					}
				}
				if action.Data.Card != nil {
					cardId := action.Data.Card.ID
					card, ok := ins.Cards[cardId]
					_, _, _, typeOfTask := ValidateTaskName(card.Name)
					if ok {
						// get memberId
						memberIdList := card.IDMembers
						for _, id := range memberIdList {
							memberStat, ok := stat.MemberStats[id]
							if ok {
								var needSaved = false
								memberAction := MemberActions{
									Time:          utils.TimeLocal(action.Date),
									NameOfMember:  memberStat.FullName,
									ContentOfTask: action.Data.Card.Name,
									TypeOfTask:    typeOfTask, // ex: Tool hieu nang
								}
								if action.Data.ListBefore != nil {
									memberAction.ListBefore = action.Data.ListBefore.Name
								}
								if action.Data.ListAfter != nil {
									memberAction.ListAfter = action.Data.ListAfter.Name
								}
								if taskDone {
									needSaved = true
									atomic.AddInt32(&memberStat.NDoneTasks, 1)
									atomic.AddInt32(&memberStat.NDoneHours, task.Hour)
									memberAction.ActionTypes = append(memberAction.ActionTypes, "Done")
									// memberAction_done := MemberActions{stat.Date, action.Data.ListBefore.Name, action.Data.ListAfter.Name, memberStat.FullName, action.Data.Card.Name}
									//logger.Debug("#", memberAction_done)
								}
								if taskInProgress {
									needSaved = true
									atomic.AddInt32(&memberStat.NProgressTasks, 1)
									atomic.AddInt32(&memberStat.NProgressHours, task.Hour)
									memberAction.ActionTypes = append(memberAction.ActionTypes, "InProgress")
								}
								if taskUndone {
									needSaved = true
									atomic.AddInt32(&memberStat.NDoneTasks, -1)
									atomic.AddInt32(&memberStat.NDoneHours, -1*task.Hour)
									memberAction.ActionTypes = append(memberAction.ActionTypes, "Undone")
								}
								if taskNotInProgress {
									needSaved = true
									atomic.AddInt32(&memberStat.NProgressTasks, -1)
									atomic.AddInt32(&memberStat.NProgressHours, -1*task.Hour)
									memberAction.ActionTypes = append(memberAction.ActionTypes, "NotInProgress")
								}
								if needSaved {
									stat.Mutex.Lock()
									stat.MemberActions[action.Date.String()] = &memberAction
									stat.Mutex.Unlock()
								}
							}
						}
					}
				}
			}
			return
		}
		current = current.next
	}
}

func (list *DateLinkedList) InitMembersDailyTracking(memberMap map[string]*MemberStats) {
	if list.head == nil {
		logger.Debugln("List is empty")
		return
	}
	current := list.head
	for current != nil {
		stat := current.stat
		if stat == nil {
			current = current.next
			continue
		}
		for id, member := range memberMap {
			memberStats := *member
			stat.MemberStats[id] = &memberStats
		}
		current = current.next
	}
}

func (list *DateLinkedList) SaveMemberActionsEachDay(date *time.Time, listBefore string, listAfter string, nameOfMember string, contentOfTask string) []string {
	memberActionsEachDay := []string{}
	if list.head == nil {
		logger.Debug("List of members actions each day is empty")
		return nil
	}
	current := list.head
	for current != nil {
		stat := current.stat
		if stat == nil {
			current = current.next
			continue
		}
		memberActionsEachDay = append(memberActionsEachDay, stat.Date.Format("02-01-2006"))
		memberActionsEachDay = append(memberActionsEachDay, listBefore)
		memberActionsEachDay = append(memberActionsEachDay, listAfter)
		memberActionsEachDay = append(memberActionsEachDay, nameOfMember)
		memberActionsEachDay = append(memberActionsEachDay, contentOfTask)
		current = current.next
	}
	return memberActionsEachDay
}

func endOfDay(t time.Time) time.Time {
	endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
	return endOfDay
}

func CheckListContains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
