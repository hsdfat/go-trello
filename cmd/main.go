package main

import (
	"fmt"
	"go-trello/logger"
	"go-trello/trello_service"
	"go-trello/utils"
	"log"
	"strings"
	"time"
	"github.com/spf13/viper"
)

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
	trello_service.ExportTotalMemberToCsv(ins)
	trello_service.DrawPieChartSMF(utils.NameSMFTeam)
	ins.DailyTrackingStats.ExportMemberActionsSprintToExcel()         //export data of tracking action in SMF
	trello_service.ExportDataOfDailyToExcel(ins)                      //data of Sheet Daily
	trello_service.DrawDailyLineChart(utils.MemberActionDaily)        //draw line chart in sheet daily
	trello_service.DrawClusteredColumnChart(utils.MemberActionDaily)  //draw column chart in Daily
	ins.DailyTrackingStats.ExportMemberActionsDailyToExcel()          //export data of tracking action in Daily
	ins.DailyTrackingStats.ExportGroupActionsSprintToExcel(ins.Tasks) //export data of tracking action in Group sheet
	trello_service.ExportDataOfMembersToExcel(ins)                    //Sheet: Data each member of team
	trello_service.DeleteSheet(utils.FileNeedDelete)
	duration := time.Since(start)
	logger.Info("Project took: ", duration)
}
