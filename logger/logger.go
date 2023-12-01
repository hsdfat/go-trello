package logger

import "log"

func Errorln(v ...any) {
	log.Println("[Error]: ", v)
}

func Debugln(v ...any) {
	log.Println("[Debug]: ", v)
}

func Fatalln(v ...any) {
	log.Fatalln("[Error]: ", v)
}