package main

import (
	"fmt"
	"go-trello/trello_service"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("Hello, World!")
	binaryPath := "./config/"
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
	trello_service.Export_csv(boardId)
}
