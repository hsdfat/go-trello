package utils

import (
	"go-trello/logger"
	"math"
	"strconv"
	"time"
)

// name sheet
var (
	MemberActionDaily = "Daily"
	NameSMFTeam       = "SMF"
	Group			  = "Group"
	NameOfFile        = "SMF-Trello.xlsx"
)

// GetYValue returns the values a, b of Expected line function (y = ax + b)
func GetYValue(a float64, x int, b int32) float64 {
	return float64(x)*a + float64(b)
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func ConvertStringToInt(str string) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		logger.Errorln(err)
	}
	return result
}

// compare day, month, year
func IsDateEqual(time1 *time.Time, time2 *time.Time) bool {
	if (time1.Year() == time2.Year()) && time1.Month() == time2.Month() && time1.Day() == time2.Day() {
		return true
	} else {
		return false
	}
}

func TimeLocal(timeLocal time.Time) time.Time {
    loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
    if err != nil {
        logger.Error(err)
    }
    return timeLocal.In(loc)
}

//InSkipDays returns true if dataNeedCheck is in skipDays
func InSkipDays(skipDate []string, dateNeedCheck time.Time) bool {
	for i := 0 ; i < len(skipDate); i++ {
		//convert to date type
		skipDateType, err := time.Parse("02-01-2006",skipDate[i])
		logger.Info("skip date type: ", skipDateType)
		if err != nil {
			logger.Error(err)		
		}
		if IsDateEqual(&skipDateType, &dateNeedCheck) {
			return true
		}
	}
	return false
}