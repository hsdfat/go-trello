package trello_service

import (
	"fmt"
	"go-trello/logger"
	"go-trello/utils"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

func creatSheet(eachMemberData *MemberStats, time time.Time) {
	nameOfSheet := eachMemberData.Name
	//open file
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new sheet.
	index, err := f.NewSheet(nameOfSheet)
	if err != nil {
		logger.Errorln(err)
	}

	//get data to sheet of each member
	f.SetCellValue(nameOfSheet, "A1", "Date")
	f.SetCellValue(nameOfSheet, "B1", "Number Of Tasks")

	f.SetActiveSheet(index)
	if err := f.SaveAs(utils.NameOfFile); err != nil {
		fmt.Println(err)
	}
}

func SetCellValue(nameOfSheet string, dataDaily []string, totalTask int, numberOfSprint int, numberOfDayToCurrentDay int) {
	//open file
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Error(err)
		}
	}()

	// Create a new sheet.
	index, err := f.NewSheet(nameOfSheet)
	if err != nil {
		logger.Error(err)
	}
	//set size of coloum
	err_size_column := f.SetColWidth(nameOfSheet, "A", "L", 15)
	if err_size_column != nil {
		fmt.Println(err_size_column)
	}

	err_size_height := f.SetRowHeight(nameOfSheet, 1, 20)
	if err_size_height != nil {
		fmt.Println(err_size_height)
	}

	var j int = 65
	var k int = 65
	var countDay int = 1
	var countDayDaily int = 1
	//get data to sheet of each member
	f.SetCellValue(nameOfSheet, "A1", "Date")
	f.SetCellValue(nameOfSheet, "A2", "Tasks")
	f.SetCellValue(nameOfSheet, "A3", "Expected")
	f.SetCellValue(nameOfSheet, "A4", "Remaining Hours")
	f.SetCellValue(nameOfSheet, "A5", "Hours")
	f.SetCellValue(nameOfSheet, "B1", "StartDay")
	f.SetCellValue(nameOfSheet, "B2", totalTask)
	f.SetCellValue(nameOfSheet, "B3", totalTask)

	for i := 0; i < len(dataDaily); i += 4 {
		date := dataDaily[i]
		f.SetCellValue(nameOfSheet, string((j+2))+"1", date)
		f.SetCellValue(nameOfSheet, string((j+2))+"3", utils.RoundFloat(utils.GetYValue(-float64(totalTask)/float64(numberOfSprint), countDay, int32(totalTask)), 2))
		j += 1
		countDay += 1
	}

	logger.Info("#: ", len(dataDaily))

	for i := 0; i < numberOfDayToCurrentDay*4; i += 4 {
		//for i := 0; i < len(dataDaily); i += 4 {
		// date := dataDaily[i]
		logger.Info("number of day to current day: ", numberOfDayToCurrentDay)
		remainingTasks := dataDaily[i+1]
		remainingHours := dataDaily[i+2]
		remainingHoursLinear := dataDaily[i+3]
		f.SetCellValue(nameOfSheet, string((k+2))+"2", utils.ConvertStringToInt(remainingTasks))       //
		f.SetCellValue(nameOfSheet, string((k+2))+"4", utils.ConvertStringToInt(remainingHours))       //
		f.SetCellValue(nameOfSheet, string((k+2))+"5", utils.ConvertStringToInt(remainingHoursLinear)) //
		k += 1
		countDayDaily += 1
	}
	f.SetActiveSheet(index)
	if err := f.SaveAs(utils.NameOfFile); err != nil {
		logger.Error(err)
	}
}

func DrawLine(nameOfSheet string) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Errorln(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Errorln(err)
		}
	}()
	lineWidth := 1.2
	err_add_shape := f.AddShape(nameOfSheet,
		&excelize.Shape{
			Cell: "G6",
			Type: "line",
			Line: excelize.ShapeLine{Color: "FF3349", Width: &lineWidth},
			//Fill: excelize.Fill{Color: []string{"8EB9FF"}},

			// Paragraph: []excelize.RichTextRun{
			//     {
			//         Text: "Rectangle Shape",
			//         Font: &excelize.Font{
			//             Bold:      true,
			//             Italic:    true,
			//             Family:    "Times New Roman",
			//             Size:      18,
			//             Color:     "777777",
			//             Underline: "sng",
			//         },
			//     },
			// },
			// Width:  180,
			// Height: 40,
		},
	)
	if err_add_shape != nil {
		logger.Errorln(err_add_shape)
	}
	if err := f.SaveAs(utils.NameOfFile); err != nil {
		logger.Errorln(err)
	}
}

func DrawLineChart(name_sheet string) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Errorln(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Errorln(err)
		}
	}()

	rowHead := string(int('B'))
	rowEnd := string(int('B') + 11)
	if err := f.AddChart(name_sheet, "A10", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				//Name:       name_sheet + "!$A$2",
				Categories: name_sheet + "!" + "$" + rowHead + "$1" + ":" + "$" + rowEnd + "$1",
				Values:     name_sheet + "!" + "$" + rowHead + "$2" + ":" + "$" + rowEnd + "$2",
				Line: excelize.ChartLine{
					Smooth: false,
					Width:  1.0,
				},
				Fill: excelize.Fill{
					//Type:    "",
					Pattern: 0,
					Color:   nil,
					Shading: 0,
				},
			},
			{
				//Name:       name_sheet + "!$A$3",
				//Categories: "Sheet1!$B$1:$D$1",
				Values: name_sheet + "!" + "$" + rowHead + "$3" + ":" + "$" + rowEnd + "$3",
				Line: excelize.ChartLine{
					Smooth: true,
					Width:  1.0,
				},
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Legend: excelize.ChartLegend{
			Position: "top",
		},
		Title: []excelize.RichTextRun{
			{
				Text: name_sheet,
			},
		},
		XAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      2,
			//TickLabelSkip:  1,
			//ReverseOrder:   false,
			//Secondary:      false,
			//Maximum:        nil,
			//Minimum:        nil,
			//Font:           excelize.Font{},
			//LogBase:        0,
			//NumFmt:         excelize.ChartNumFmt{},
			Title: []excelize.RichTextRun{
				{
					Text: "Date",
				},
			},
		},

		YAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      5,
			//TickLabelSkip:  1,
			//ReverseOrder:   false,
			//Secondary:      false,
			//Maximum:        nil,
			//Minimum:        nil,
			//Font:           excelize.Font{},
			//LogBase:        0,
			//NumFmt:         excelize.ChartNumFmt{},
			Title: []excelize.RichTextRun{
				{
					Text: "Remaining Tasks",
				},
			},
		},

		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowSerName:     true,
			ShowVal:         true,
		},
		Dimension: excelize.ChartDimension{
			Width:  650,
			Height: 500,
		},
		ShowBlanksAs: "span",
		HoleSize:     3,
	}); err != nil {
		fmt.Println(err)
		return
	}

	// lineWidth := 1.2
	// err_add_shape := f.AddShape(name_sheet,
	// &excelize.Shape{
	//     Cell: "G6",
	//     Type: "line",
	//     Line: excelize.ShapeLine{Color: "FF3349", Width: &lineWidth},
	//Fill: excelize.Fill{Color: []string{"8EB9FF"}},

	// Paragraph: []excelize.RichTextRun{
	//     {
	//         Text: "Rectangle Shape",
	//         Font: &excelize.Font{
	//             Bold:      true,
	//             Italic:    true,
	//             Family:    "Times New Roman",
	//             Size:      18,
	//             Color:     "777777",
	//             Underline: "sng",
	//         },
	//     },
	// },
	// Width:  180,
	// Height: 40,
	// },)
	// if err_add_shape != nil {
	// 	logger.Errorln(err_add_shape)
	// }

	if err := f.SaveAs(utils.NameOfFile); err != nil {
		fmt.Println(err)
	}
}

func DrawLineChartForTotal(name_sheet string) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Errorln(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Errorln(err)
		}
	}()

	rowHead := string(int('B'))
	rowEnd := string(int('B') + 7)
	if err := f.AddChart(name_sheet, "A17", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				//Name:       name_sheet + "!$A$2",
				Name:       name_sheet,
				Categories: name_sheet + "!" + "$" + rowHead + "$13" + ":" + "$" + rowEnd + "$13",
				Values:     name_sheet + "!" + "$" + rowHead + "$13" + ":" + "$" + rowEnd + "$13",
				Line: excelize.ChartLine{
					Smooth: false,
					Width:  1.0,
				},
				Fill: excelize.Fill{
					//Type:    "",
					Pattern: 0,
					Color:   nil,
					Shading: 0,
				},
			},
			{
				//Name:       name_sheet + "!$A$3",
				//Categories: "Sheet1!$B$1:$D$1",
				Values: name_sheet + "!" + "$" + rowHead + "$3" + ":" + "$" + rowEnd + "$3",
				Line: excelize.ChartLine{
					Smooth: true,
					Width:  1.0,
				},
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Legend: excelize.ChartLegend{
			Position: "top",
		},
		Title: []excelize.RichTextRun{
			{
				Text: name_sheet,
			},
		},
		XAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      2,
			//TickLabelSkip:  1,
			//ReverseOrder:   false,
			//Secondary:      false,
			//Maximum:        nil,
			//Minimum:        nil,
			//Font:           excelize.Font{},
			//LogBase:        0,
			//NumFmt:         excelize.ChartNumFmt{},
			Title: []excelize.RichTextRun{
				{
					Text: "Date",
				},
			},
		},

		YAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      5,
			//TickLabelSkip:  1,
			//ReverseOrder:   false,
			//Secondary:      false,
			//Maximum:        nil,
			//Minimum:        nil,
			//Font:           excelize.Font{},
			//LogBase:        0,
			//NumFmt:         excelize.ChartNumFmt{},
			Title: []excelize.RichTextRun{
				{
					Text: "Remaining Tasks",
				},
			},
		},

		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowSerName:     true,
			ShowVal:         true,
		},
		Dimension: excelize.ChartDimension{
			Width:  650,
			Height: 500,
		},
		ShowBlanksAs: "span",
		HoleSize:     3,
	}); err != nil {
		fmt.Println(err)
		return
	}

	if err := f.SaveAs(utils.NameOfFile); err != nil {
		fmt.Println(err)
	}
}

func DrawDailyLineChart(name_sheet string) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Errorln(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Errorln(err)
		}
	}()

	rowHead := string(int('B'))
	rowEnd := string(int('B') + 11)
	if err := f.AddChart(name_sheet, "A10", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				//Name:       name_sheet + "!$A$2",
				Categories: name_sheet + "!" + "$" + rowHead + "$1" + ":" + "$" + rowEnd + "$1",
				Values:     name_sheet + "!" + "$" + rowHead + "$2" + ":" + "$" + rowEnd + "$2",
				Line: excelize.ChartLine{
					Smooth: false,
					Width:  1.0,
				},
				Fill: excelize.Fill{
					//Type:    "",
					Pattern: 0,
					Color:   nil,
					Shading: 0,
				},
			},
			{
				//Name:       name_sheet + "!$A$3",
				//Categories: "Sheet1!$B$1:$D$1",
				Values: name_sheet + "!" + "$" + rowHead + "$3" + ":" + "$" + rowEnd + "$3",
				Line: excelize.ChartLine{
					Smooth: true,
					Width:  1.0,
				},
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Legend: excelize.ChartLegend{
			Position: "top",
		},
		Title: []excelize.RichTextRun{
			{
				Text: name_sheet,
			},
		},
		XAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      2,
			Title: []excelize.RichTextRun{
				{
					Text: "Date",
				},
			},
		},

		YAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      5,
			Title: []excelize.RichTextRun{
				{
					Text: "Remaining Tasks",
				},
			},
		},

		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowSerName:     true,
			ShowVal:         true,
		},
		Dimension: excelize.ChartDimension{
			Width:  650,
			Height: 500,
		},
		ShowBlanksAs: "span",
		HoleSize:     3,
	}); err != nil {
		fmt.Println(err)
		return
	}

	if err := f.SaveAs(utils.NameOfFile); err != nil {
		fmt.Println(err)
	}
}

func DrawRemainingHours(name_sheet string) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Errorln(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Errorln(err)
		}
	}()

	rowHead := string(int('B'))
	rowEnd := string(int('B') + 11)
	if err := f.AddChart(name_sheet, "A10", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				//Name:       name_sheet + "!$A$2",
				Categories: name_sheet + "!" + "$" + rowHead + "$1" + ":" + "$" + rowEnd + "$1",
				Values:     name_sheet + "!" + "$" + rowHead + "$2" + ":" + "$" + rowEnd + "$2",
				Line: excelize.ChartLine{
					Smooth: false,
					Width:  1.0,
				},
				Fill: excelize.Fill{
					//Type:    "",
					Pattern: 0,
					Color:   nil,
					Shading: 0,
				},
			},
			{
				//Name:       name_sheet + "!$A$3",
				//Categories: "Sheet1!$B$1:$D$1",
				Values: name_sheet + "!" + "$" + rowHead + "$3" + ":" + "$" + rowEnd + "$3",
				Line: excelize.ChartLine{
					Smooth: true,
					Width:  1.0,
				},
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Legend: excelize.ChartLegend{
			Position: "top",
		},
		Title: []excelize.RichTextRun{
			{
				Text: name_sheet,
			},
		},
		XAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      2,
			Title: []excelize.RichTextRun{
				{
					Text: "Date",
				},
			},
		},

		YAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      5,
			Title: []excelize.RichTextRun{
				{
					Text: "Remaining Tasks",
				},
			},
		},

		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowSerName:     true,
			ShowVal:         true,
		},
		Dimension: excelize.ChartDimension{
			Width:  650,
			Height: 500,
		},
		ShowBlanksAs: "span",
		HoleSize:     3,
	}); err != nil {
		fmt.Println(err)
		return
	}

	if err := f.SaveAs(utils.NameOfFile); err != nil {
		fmt.Println(err)
	}
}

func DrawClusteredColumnChart(name_sheet string) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Errorln(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Errorln(err)
		}
	}()

	rowHead := string(int('B'))
	rowEnd := string(int('B') + 11)
	if err := f.AddChart(name_sheet, "G10", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				//Name:       name_sheet + "!$A$5",
				Categories: name_sheet + "!" + "$" + rowHead + "$1" + ":" + "$" + rowEnd + "$1",
				Values:     name_sheet + "!" + "$" + rowHead + "$5" + ":" + "$" + rowEnd + "$5",
			},
			{
				//Name:       name_sheet + "!$A$4",
				Categories: name_sheet + "!" + "$" + rowHead + "$1" + ":" + "$" + rowEnd + "$1",
				Values:     name_sheet + "!" + "$" + rowHead + "$4" + ":" + "$" + rowEnd + "$4",
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Legend: excelize.ChartLegend{
			Position: "left",
		},
		Title: []excelize.RichTextRun{
			{
				Text: name_sheet,
			},
		},
		XAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      2,
			Title: []excelize.RichTextRun{
				{
					Text: "Date",
				},
			},
		},

		YAxis: excelize.ChartAxis{
			None:           false,
			MajorGridLines: false,
			MinorGridLines: true,
			MajorUnit:      5,
			Title: []excelize.RichTextRun{
				{
					Text: "Remaining Hours",
				},
			},
		},

		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowSerName:     true,
			ShowVal:         true,
		},
		Dimension: excelize.ChartDimension{
			Width:  650,
			Height: 500,
		},
		ShowBlanksAs: "span",
		HoleSize:     3,
	}); err != nil {
		fmt.Println(err)
		return
	}

	if err := f.SaveAs(utils.NameOfFile); err != nil {
		fmt.Println(err)
	}
}

func DrawPieChartSMF(nameSheet string) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Error(err)
		}
	}()

	if err := f.AddChart(nameSheet, "A15", &excelize.Chart{
		Type: excelize.Pie,
		Series: []excelize.ChartSeries{
			{
				Name:       "Amount",
				Categories: nameSheet + "!$B$1:$D$1",
				Values:     nameSheet + "!$B$13:$D$13",
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Title: []excelize.RichTextRun{
			{
				Text: nameSheet,
			},
		},
		PlotArea: excelize.ChartPlotArea{
			ShowPercent: true,
		},
	}); err != nil {
		fmt.Println(err)
		return
	}
	if err := f.SaveAs(utils.NameOfFile); err != nil {
		fmt.Println(err)
	}
}

func SetMemberActionsDaily(memberActionDaily string, memberActions []*MemberActions) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Error(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error(err)
		}
	}()
	columnSizeErr := f.SetColWidth(memberActionDaily, "N", "S", 15)
	if columnSizeErr != nil {
		logger.Error(columnSizeErr)
	}
	f.SetCellValue(memberActionDaily, "N10", "Time")
	f.SetCellValue(memberActionDaily, "O10", "List Before")
	f.SetCellValue(memberActionDaily, "P10", "List After")
	f.SetCellValue(memberActionDaily, "Q10", "Name")
	f.SetCellValue(memberActionDaily, "R10", "Task")
	f.SetCellValue(memberActionDaily, "S10", "Action Types")
	row := 11
	for _, memberAction := range memberActions {
		// logger.Info("**************************************************")
		// logger.Info(memberAction.Time)
		// logger.Info(memberAction.ListBefore)
		// logger.Info(memberAction.ListAfter)
		// logger.Info(memberAction.NameOfMember)
		// logger.Info(memberAction.ContentOfTask)
		// logger.Info(memberAction.ActionTypes)
		f.SetCellValue(memberActionDaily, "N"+strconv.Itoa(row), memberAction.Time)
		f.SetCellValue(memberActionDaily, "O"+strconv.Itoa(row), memberAction.ListBefore)
		f.SetCellValue(memberActionDaily, "P"+strconv.Itoa(row), memberAction.ListAfter)
		f.SetCellValue(memberActionDaily, "Q"+strconv.Itoa(row), memberAction.NameOfMember)
		f.SetCellValue(memberActionDaily, "R"+strconv.Itoa(row), memberAction.ContentOfTask)
		f.SetCellValue(memberActionDaily, "S"+strconv.Itoa(row), memberAction.ActionTypes)
		row += 1
	}
	if err := f.SaveAs(utils.NameOfFile); err != nil {
		logger.Error(err)
	}
}

func SetMemberActionsSprint(nameOfSheet string, memberActions []*MemberActions) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Error(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error(err)
		}
	}()
	columnSizeErr := f.SetColWidth(nameOfSheet, "J", "O", 15)
	if columnSizeErr != nil {
		logger.Error(columnSizeErr)
	}
	f.SetCellValue(nameOfSheet, "K1", "Time")
	f.SetCellValue(nameOfSheet, "L1", "List Before")
	f.SetCellValue(nameOfSheet, "M1", "List After")
	f.SetCellValue(nameOfSheet, "N1", "Name")
	f.SetCellValue(nameOfSheet, "O1", "Task")
	f.SetCellValue(nameOfSheet, "P1", "Action Types")
	row := 2
	for _, memberAction := range memberActions {
		// logger.Info("--------------------------------")
		// logger.Info("memberAction.Time: ", memberAction.Time)
		// logger.Info("memberAction.ListBefore: ", memberAction.ListBefore)
		// logger.Info("memberAction.ListAfter: ", memberAction.ListAfter)
		// logger.Info("memberAction.NameOfMember: ", memberAction.NameOfMember)
		// logger.Info("memberAction.ContentOfTask: ", memberAction.ContentOfTask)
		// logger.Info("memberAction.ActionTypes: ", memberAction.ActionTypes)
		// logger.Info("memberAction.TypeOfTask: ", memberAction.TypeOfTask)
		
		f.SetCellValue(nameOfSheet, "K"+strconv.Itoa(row), memberAction.Time)
		f.SetCellValue(nameOfSheet, "L"+strconv.Itoa(row), memberAction.ListBefore)
		f.SetCellValue(nameOfSheet, "M"+strconv.Itoa(row), memberAction.ListAfter)
		f.SetCellValue(nameOfSheet, "N"+strconv.Itoa(row), memberAction.NameOfMember)
		f.SetCellValue(nameOfSheet, "O"+strconv.Itoa(row), memberAction.ContentOfTask)
		f.SetCellValue(nameOfSheet, "P"+strconv.Itoa(row), memberAction.ActionTypes)
		row += 1
	}
	if err := f.SaveAs(utils.NameOfFile); err != nil {
		logger.Error(err)
	}
}

func SetGroupActionsSprint(nameOfSheet string, tasks []*Task) {
	f, err := excelize.OpenFile(utils.NameOfFile)
	if err != nil {
		logger.Error(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet(nameOfSheet)
	if err != nil {
		logger.Errorln(err)
	}
	f.SetCellValue(nameOfSheet, "A1", "Type of task")
	f.SetCellValue(nameOfSheet, "B1", "Card name")
	f.SetCellValue(nameOfSheet, "C1", "Status")
	f.SetCellValue(nameOfSheet, "D1", "Owner")
	row := 2
	for _, task := range tasks {
		logger.Info("************************")
		statusOfTask := GetStatusOfTaskInGroupSheet(task)
		logger.Info("Kind of task: ", task.TypeOfTask)
		logger.Info("Name of member: ", task.Card.Name)
		logger.Info("memberAction.Current Status: ", statusOfTask)
		logger.Info("member ", ConvertNameOfMember(task.Members.Username))
		
		f.SetCellValue(nameOfSheet, "A"+strconv.Itoa(row), task.TypeOfTask)
		f.SetCellValue(nameOfSheet, "B"+strconv.Itoa(row), task.Card.Name)
		f.SetCellValue(nameOfSheet, "C"+strconv.Itoa(row), statusOfTask)
		f.SetCellValue(nameOfSheet, "D"+strconv.Itoa(row), ConvertNameOfMember(task.Members.Username))
		// f.SetCellValue(nameOfSheet, "E"+strconv.Itoa(row), task.ContentOfTask)
		row += 1
	}
	f.SetActiveSheet(index)
	if err := f.SaveAs(utils.NameOfFile); err != nil {
		logger.Error(err)
	}
}

func DeleteSheet(nameOfSheet string) {
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
	errDeleteSheet := f.DeleteSheet(nameOfSheet)
	if errDeleteSheet != nil {
		logger.Errorln(errDeleteSheet)
	}
	// f.SetActiveSheet(index)
	if err := f.SaveAs(utils.NameOfFile); err != nil {
		logger.Error(err)
	}
}
