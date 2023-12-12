package main

import (
	"fmt"
	"go-trello/logger"
	"go-trello/trello_service"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func main() {
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
	//err = trello_service.GetBoardInfo(boardId)
	//if err != nil {
	//	log.Println(err)
	//}

	startDay := viper.GetString("trello.startDay")
	startDayTime, err := time.Parse("02-01-2006", startDay)
	if err != nil {
		log.Panicln("Cannot parse start day: ", err)
	}
	endDay := viper.GetString("trello.endDay")
	endDayTime, err := time.Parse("02-01-2006", endDay)
	if err != nil {
		log.Panicln("Cannot parse end day: ", err)
	}

	ins := trello_service.GetBoardInfo(boardId, startDayTime, endDayTime)
	ins.DailyTrackingStats.PrintMemberActions()
	trello_service.ExportTotalMemberToCsv(ins)
	trello_service.DrawPieChart()
	trello_service.ExportDataOfMembersToExcel(ins)
	trello_service.ExportDataOfDailyToExcel(ins)
	trello_service.DrawDailyLineChart("Daily")
	trello_service.DrawClusteredColumnChart("Daily")
	//trello_service.DrawLineChartForTotal("SMF")

}
