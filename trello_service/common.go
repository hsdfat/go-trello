package trello_service

import (
	"fmt"
	"go-trello/logger"
	"strconv"
	"time"
)

// GetCreationTime gets the creation time of card/action's id
func GetCreationTime(id string) (*time.Time, error) {
	// TODO: Convert id to time
	// Get first 8 characters of id
	if len(id) < 8 {
		return nil, fmt.Errorf("invalid id: %v", id)
	}
	hexString := id[:8]
	// Convert id to unix timestamp: hex->decimal format
	result, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v, cannot convert from hex to unix number timestamp", id)
	}
	createTime := time.Unix(result, 0)
	return &createTime, nil
}

// SetSprintStartDay set the start day of sprint
func (c *TrelloClient) SetSprintStartDay(t time.Time) {
	c.SprintStartDay = &t
}

// SetSprintEndDay set the start day of sprint
func (c *TrelloClient) SetSprintEndDay(t time.Time) {
	c.SprintEndDay = &t
}

// SetSprintDuration set start day end day of sprint
func (c *TrelloClient) SetSprintDuration(start, end time.Time) error {
	if !end.After(start) {
		return fmt.Errorf("end time must be after start time")
	}
	c.SetSprintStartDay(start)
	c.SetSprintEndDay(end)
	// Init Statitics Tracking Day Map

	c.DailyTrackingStats = &DateLinkedList{}
	date := start
	var dateBefore time.Time
	var dateStat *DailyTrackingStats
	logger.Debugln("Init Statitics Tracking Day", end.Add(23*time.Hour+3590*time.Second).After(date))
	for end.Add(23*time.Hour + 3590*time.Second).After(date) {
		if date.Weekday() == time.Sunday || date.Weekday() == time.Saturday {
			date = date.Add(24 * time.Hour)
			continue
		}
		dateStat = &DailyTrackingStats{
			Date:        date,
			DateBefore:  dateBefore,
			MemberStats: make(map[string]*MemberStats),
		}

		c.DailyTrackingStats.AddNodeAtEnd(dateStat)

		dateBefore = date
		date = date.Add(24 * time.Hour)
	}
	return nil
}

func (c *TrelloClient) GetCreateCurrentDate(t *Task) (date string, err error) {
	// if c == nil || c.Board == nil {
	// 	return "", fmt.Errorf("no b0oard specified, cannot get create current date")
	// }
	// if t.CreationTime == nil {
	// 	return "", fmt.Errorf("no creation time specified, cannot get")
	// }
	// for date, dailyStat := range c.DailyTrackingStats {
	// 	if dailyStat.Date.After(*t.CreationTime) && dailyStat.DateBefore.Before(*t.CreationTime) {
	// 		return date, nil
	// 	}
	// }

	return "", fmt.Errorf("")
}

func (c *TrelloClient) CheckTaskInSprint(t *Task) bool {
	if t.CreationTime == nil {
		return false
	}
	if c.SprintEndDay != nil && !c.SprintEndDay.After(*t.CreationTime) {
		return false
	}

	return true
}
