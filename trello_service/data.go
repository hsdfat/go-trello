package trello_service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go-trello/logger"
	"go-trello/utils"
)

func (list *DateLinkedList) ExportDataOfSMFTeamToExcel(id string, totalTask int32, numberOfSprint int, totalHours int32) {
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
		fmt.Println("$$: ", string((i+2))+"1")
		f.SetCellValue(memberStat.Name, string((i+2))+"1", date)
		f.SetCellValue(memberStat.Name, string((i+2))+"2", numberOfTasksNeedDone)
		fmt.Println("!!numberOfSprint: ", numberOfSprint)
		fmt.Println("!!:countDay ", countDay)
		fmt.Println("!!:totalTask ", totalTask)
		f.SetCellValue(memberStat.Name, string((i+2))+"3", utils.GetYValue(-float64(totalTask)/float64(numberOfSprint), countDay, totalTask))
		//f.SetCellValue(memberStat.Name, string((i+2))+"4", )
		//f.SetCellValue(memberStat.Name, strconv.Itoa((i+2))+"1", date)
		//f.SetCellValue(memberStat.Name, strconv.Itoa((i+2))+"2", numberOfTasksNeedDone)
		countDay += 1
		i += 1
		numberOfTasksNeedDone = numberOfTasksNeedDone + memberStat.NTasks - memberStat.NDoneTasks
		numberOfRemainingHours = numberOfRemainingHours - memberStat.NDoneHours
		//DrawLineChart(f, memberStat.FullName)
		name_sheet := memberStat.Name
		fmt.Println("@@@: ", name_sheet)

		f.SetActiveSheet(index)
		current = current.next
	}

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
