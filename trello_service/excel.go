package trello_service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go-trello/logger"
	"time"
)

func creatSheet(eachMemberData *MemberStats, time time.Time) {
	nameOfSheet := eachMemberData.Name
	//open file
	f, err := excelize.OpenFile("Book1.xlsx")
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
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func DrawLineChart(name_sheet string) {
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

	rowHead := string(int('B'))
	rowEnd := string(int('B') + 11)
	if err := f.AddChart(name_sheet, "A10", &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				//Name:       name_sheet,
				Categories: name_sheet + "!" + "$" + rowHead + "$1" + ":" + "$" + rowEnd + "$1",
				Values:     name_sheet + "!" + "$" + rowHead + "$2" + ":" + "$" + rowEnd + "$2",
				Line: excelize.ChartLine{
					Smooth: false,
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

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
