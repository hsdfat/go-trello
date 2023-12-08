package trello_service

import (
	"fmt"
	"go-trello/logger"
	"go-trello/utils"
	"sync"
	"sync/atomic"
	"time"

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
		logger.Debugln(fmt.Sprintf("date [%s]: new task (done/progress/total): %d/%d/%d, new hour (done/progress/total): %d/%d/%d\t",
			stat.Date.Format("02-01-2006"), stat.NDoneTasks, stat.NProgressTasks, stat.NTasks, stat.NDoneHours, stat.NProgressHours, stat.NHours))
		current = current.next
	}
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

		logger.Debugln(fmt.Sprintf("member: [%s] date [%s]: new task (done/progress/extra task): %d/%d/%d, new hour (done/progress/extra task): %d/%d/%d\t", memberStat.FullName,
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

func (list *DateLinkedList) ExportDataOfEachMemberToExcel(id string, totalTask int32, numberOfSprint int, totalHours int32) {
	numberOfTasksNeedDone := totalTask
	numberOfRemainingHours := totalHours

	//export to excel
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		logger.Errorln(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Errorln(err)
		}
	}()

	if list.head == nil {
		logger.Debugln("List is empty")
		return
	}

	current := list.head
	var i int = 64
	var countDay int = 1
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
		//get data to sheet of each member
		f.SetCellValue(memberStat.Name, "A1", "Date")
		f.SetCellValue(memberStat.Name, "A2", "Tasks")
		f.SetCellValue(memberStat.Name, "A3", "Expected")
		f.SetCellValue(memberStat.Name, "A4", "Hours")

		// Create a new sheet.
		index, err := f.NewSheet(memberStat.Name)
		if err != nil {
			logger.Errorln(err)
		}
		date := fmt.Sprintf("%s", stat.Date.Format("02-01-2006"))

		f.SetCellValue(memberStat.Name, string((i+2))+"1", date)
		f.SetCellValue(memberStat.Name, string((i+2))+"2", numberOfTasksNeedDone)

		expected_task := utils.RoundFloat(utils.GetYValue(-float64(totalTask)/float64(numberOfSprint), countDay, totalTask), 2)
		fmt.Println("************", expected_task)
		f.SetCellValue(
			memberStat.Name,
			string((i+2))+"3",
			expected_task,
		)
		countDay += 1
		i += 1
		numberOfTasksNeedDone = numberOfTasksNeedDone + memberStat.NTasks - memberStat.NDoneTasks
		numberOfRemainingHours = numberOfRemainingHours - memberStat.NDoneHours
		//DrawLineChart(f, memberStat.FullName)
		//name_sheet := memberStat.Name
		//fmt.Println("@@@: ", name_sheet)

		//set size of coloum
		err_size_column := f.SetColWidth(memberStat.Name, "A", "L", 15)
		if err_size_column != nil {
			fmt.Println(err_size_column)
		}

		err_size_height := f.SetRowHeight(memberStat.Name, 1, 20)
		if err_size_height != nil {
			fmt.Println(err_size_height)
		}

		f.SetActiveSheet(index)
		current = current.next
	}

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func (list *DateLinkedList) ExportDataOfSMFTeamToExcel(nameOfSheet string, id string, totalTask int32, numberOfSprint int, totalHours int32) {
	numberOfTasksNeedDone := totalTask
	numberOfRemainingHours := totalHours

	//export to excel
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		logger.Errorln(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Errorln(err)
		}
	}()

	if list.head == nil {
		logger.Debugln("List is empty")
		return
	}

	current := list.head
	var i int = 64
	var countDay int = 1
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
		//get data to sheet of each member
		f.SetCellValue(nameOfSheet, "A1", "Date")
		f.SetCellValue(nameOfSheet, "A2", "Tasks")
		f.SetCellValue(nameOfSheet, "A3", "Expected")
		f.SetCellValue(nameOfSheet, "A4", "Hours")

		// Create a new sheet.
		index, err := f.NewSheet(nameOfSheet)
		if err != nil {
			logger.Errorln(err)
		}
		date := fmt.Sprintf("%s", stat.Date.Format("02-01-2006"))
		fmt.Println("$$: ", string((i+2))+"1")
		f.SetCellValue(memberStat.Name, string((i+2))+"1", date)
		f.SetCellValue(memberStat.Name, string((i+2))+"2", numberOfTasksNeedDone)

		expected_task := utils.RoundFloat(utils.GetYValue(-float64(totalTask)/float64(numberOfSprint), countDay, totalTask), 2)
		fmt.Println("************", expected_task)
		f.SetCellValue(
			memberStat.Name,
			string((i+2))+"3",
			expected_task,
		)
		countDay += 1
		i += 1
		numberOfTasksNeedDone = numberOfTasksNeedDone + memberStat.NTasks - memberStat.NDoneTasks
		numberOfRemainingHours = numberOfRemainingHours - memberStat.NDoneHours
		//DrawLineChart(f, memberStat.FullName)
		name_sheet := memberStat.Name
		fmt.Println("@@@: ", name_sheet)

		//set size of coloum
		err_size_column := f.SetColWidth(memberStat.Name, "A", "L", 15)
		if err_size_column != nil {
			fmt.Println(err_size_column)
		}

		err_size_height := f.SetRowHeight(memberStat.Name, 1, 20)
		if err_size_height != nil {
			fmt.Println(err_size_height)
		}

		f.SetActiveSheet(index)
		current = current.next
	}

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
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

func (list *DateLinkedList) TrackingAction(task *Task, action *trello.Action, wg *sync.WaitGroup) {
	defer wg.Done()

	ins := GetInstance()
	if list.head == nil {
		logger.Debugln("List is empty")
		return
	}
	current := list.head
	for current != nil {
		stat := current.stat
		if endOfDay(stat.Date).After(action.Date) {
			if action.Data != nil {
				var taskDone, taskUndone, taskInProgress, taskNotInProgress bool
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
						taskNotInProgress = true
						atomic.AddInt32(&stat.NProgressTasks, -1)
						atomic.AddInt32(&stat.NProgressHours, -1*task.Hour)
					}
				}
				if action.Data.Card != nil {
					cardId := action.Data.Card.ID
					card, ok := ins.Cards[cardId]
					if ok {
						// get memberId
						memberIdList := card.IDMembers
						for _, id := range memberIdList {
							memberStat, ok := stat.MemberStats[id]
							if ok {
								if taskDone {
									atomic.AddInt32(&memberStat.NDoneTasks, 1)
									atomic.AddInt32(&memberStat.NDoneHours, task.Hour)
								}
								if taskInProgress {
									atomic.AddInt32(&memberStat.NProgressTasks, 1)
									atomic.AddInt32(&memberStat.NProgressHours, task.Hour)
								}
								if taskUndone {
									atomic.AddInt32(&memberStat.NDoneTasks, -1)
									atomic.AddInt32(&memberStat.NDoneHours, -1*task.Hour)
								}
								if taskNotInProgress {
									atomic.AddInt32(&memberStat.NProgressTasks, -1)
									atomic.AddInt32(&memberStat.NProgressHours, -1*task.Hour)
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
