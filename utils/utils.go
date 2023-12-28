package utils

import (
	"go-trello/logger"
	"math"
	"strconv"
	"time"
	"github.com/spf13/viper"
)

// name sheet
var (
	MemberActionDaily = "Daily"
	NameSMFTeam       = "SMF"
	Group			  = "Group"
	NameOfFile        = "SMF-Trello.xlsx"
	FileNeedDelete    = "Sheet1"
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
		if err != nil {
			logger.Error(err)		
		}
		if IsDateEqual(&skipDateType, &dateNeedCheck) {
			return true
		}
	}
	return false
}

func FindHourOfEachMember(name string) int {
	var initLinearTime int
	if name == "dongnt18" {
		initLinearTime = viper.GetInt("smfHourMembers.dongnt18")
	} else if name == "minhdk" {
		initLinearTime = viper.GetInt("smfHourMembers.minhdk")
	}  else if name == "binhtd7" {
		initLinearTime = viper.GetInt("smfHourMembers.binhtd7")
	}  else if name == "tuanna137" {
		initLinearTime = viper.GetInt("smfHourMembers.tuanna137")
	}  else if name == "namnp16" {
		initLinearTime = viper.GetInt("smfHourMembers.namnp16")
	}  else if name == "hungtq40" {
		initLinearTime = viper.GetInt("smfHourMembers.hungtq40")
	}  else if name == "phatlc" {
		initLinearTime = viper.GetInt("smfHourMembers.phatlc")
	}  else if name == "trinhmn" {
		initLinearTime = viper.GetInt("smfHourMembers.trinhmn")
	}  else if name == "thinhnx5" {
		initLinearTime = viper.GetInt("smfHourMembers.thinhnx5")
	}  else if name == "hieund152" {
		initLinearTime = viper.GetInt("smfHourMembers.hieund152")
	}  else if name == "bachtx" {
		initLinearTime = viper.GetInt("smfHourMembers.bachtx")
	} 
	return initLinearTime
}