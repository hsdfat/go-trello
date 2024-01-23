package main

import (
	"fmt"
	"go-trello/logger"
	"go-trello/trello_service"
	"go-trello/utils"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
)

func init() {
	logger.Info("Update")
	currentTime := time.Now()
	timeFormat := fmt.Sprint(currentTime.Format("02-01-2006_15-04-05"))
	currentDay := fmt.Sprint(currentTime.Format("02-01-2006"))
	user, err := user.Current()
	if err != nil {
		logger.Error(err)
	}
	desktop := user.HomeDir + utils.DirSaveFile
	//creat folder to save file
	err = os.MkdirAll(desktop+"Report_Trello/Sprint"+strconv.Itoa(utils.Sprint)+"/"+currentDay, os.ModePerm)
	if err != nil {
		logger.Errorln(err)
	}
	*utils.PointerNameOfFile = desktop + "Report_Trello/Sprint" + strconv.Itoa(utils.Sprint) + "/" + currentDay + "/" + "SMF-Trello_" + timeFormat + ".xlsx"
	logger.Info("utils NameOfFile: ", utils.NameOfFile)
}

func main() {
	start := time.Now()
	fmt.Println("Hello, World!")
	binaryPath := "./config/"
	logger.SetLogLevel(5)
	viper.AddConfigPath(binaryPath)
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	log.Println("Api key:", viper.GetString("trello.apiKey"))
	viper.WatchConfig()

	trello_service.Start()
	boardId := viper.GetString("trello.boardId")
	startDay := viper.GetString("trello.startDay")
	startDayTime, err := time.Parse("02-01-2006", startDay)
	if err != nil {
		log.Panicln("Cannot parse start day: ", err)
	}
	endDay := viper.GetString("trello.endDay")
	endDayTime, err := time.Parse("02-01-2006", endDay)

	startDayTime = utils.TimeLocal(startDayTime)
	endDayTime = utils.TimeLocal(endDayTime)
	if err != nil {
		log.Panicln("Cannot parse end day: ", err)
	}

	ins := trello_service.GetBoardInfo(boardId, startDayTime, endDayTime)
	//ins.DailyTrackingStats.PrintMemberActions()
	trello_service.ConvertNameOfMembers(ins)
	trello_service.ExportTotalMemberToCsv(ins) //export data of SMF Sheet
	trello_service.DrawPieChartSMF(utils.NameSMFTeam)
	ins.DailyTrackingStats.ExportMemberActionsSprintToExcel()                                                    //export data of tracking action in SMF
	trello_service.ExportDataOfDailyToExcel(ins)                                                                 //data of Daily Sheet
	trello_service.DrawDailyLineChart(utils.MemberActionDaily, ins.DailyTrackingStats.CountDaysInSprint())       //draw line chart in sheet daily
	trello_service.DrawClusteredColumnChart(utils.MemberActionDaily, ins.DailyTrackingStats.CountDaysInSprint()) //draw column chart in Daily
	ins.DailyTrackingStats.ExportMemberActionsDailyToExcel()                                                     //export data of tracking action in Daily
	ins.DailyTrackingStats.ExportGroupActionsSprintToExcel(ins.Tasks)                                            //export data of tracking action in Group sheet
	trello_service.ExportDataOfMembersToExcel(ins)                                                               //Sheet: Data each member of team
	trello_service.DeleteSheet(utils.FileNeedDelete)

	//basic progress bar idea
	duration := time.Since(start)
	durationInt := int(duration / time.Second)
	logger.Info("Project took: ", durationInt)
	//bar := progressbar.Default(int64(durationInt))
	bar := progressbar.NewOptions(durationInt,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		// progressbar.OptionSetDescription("[cyan][1/3][reset] Processing ..."),
		progressbar.OptionSetDescription("Processing ..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for i := 0; i < durationInt; i++ {
		bar.Add(1)
		time.Sleep(100 * time.Millisecond)
	}
}
