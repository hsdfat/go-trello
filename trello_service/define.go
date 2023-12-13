package trello_service

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/adlio/trello"
	"github.com/spf13/viper"
)

type TrelloClient struct {
	Client *trello.Client
	Board  *trello.Board

	Members  map[string]*trello.Member
	Lists    map[string]*trello.List
	Actions  map[string]*trello.Action
	Cards    map[string]*trello.Card
	Labels   map[string]*trello.Label
	Caretory map[string]string
	Tasks 	 []*Task

	DoneList       string
	SkipLists      []string
	SprintStartDay *time.Time
	SprintEndDay   *time.Time

	MemberStats        map[string]*MemberStats
	DailyTrackingStats *DateLinkedList
}

type MemberActions struct {
	Time          time.Time
	ListBefore    string
	ListAfter     string
	NameOfMember  string
	ContentOfTask string
	ActionTypes []string
}

type MemberStats struct {
	Name           string
	Email          string
	FullName       string
	TotalTasks     []*Task
	NTasks         int32
	NHours         int32
	NProgressTasks int32
	NProgressHours int32
	NDoneTasks     int32
	NDoneHours     int32
	NExtraTasks    int32 // tasks
	NExtraHours    int32 // hours
	Actions        *trello.Action
}

type DailyTrackingStats struct {
	DateBefore time.Time
	Date       time.Time
	Mutex      *sync.Mutex

	MemberStats    map[string]*MemberStats
	MemberActions  map[string]*MemberActions
	NTasks         int32
	NHours         int32
	NProgressTasks int32
	NProgressHours int32
	NDoneTasks     int32
	NDoneHours     int32
}

func (d *DailyTrackingStats) CheckInCurrentTime(t time.Time) bool {
	return false
}

type cardResult struct {
	Error       error
	Date        string
	Complete    bool
	Points      float64
	TrelloError bool
}

type Board struct {
	ID              string `gorm:"primary_key"`
	Name            string
	DateStart       time.Time
	DateEnd         time.Time
	Cards           uint
	Points          float64
	CardsCompleted  uint
	PointsCompleted float64
	CardProgress    []CardProgress
}

type Task struct {
	Card         *trello.Card
	CreationTime *time.Time
	IsDone       bool
	IsInProgress bool
	IsExtra      bool
	Hour         int32
}

// CardProgress represents the progress of a card.
type CardProgress struct {
	BoardID string
	Date    time.Time
	Points  float64
}

var client *trello.Client

// Start starts watching boards that are active. Refreshes according
// to the refresh rate set in the configuration.
func Start() {
	client = trello.NewClient(
		viper.GetString("trello.apiKey"),
		viper.GetString("trello.token"),
	)
	// go runBoards()
	// ch := gocron.Start()
	// refreshRate := uint64(viper.GetInt64("trello.refreshRate"))
	// gocron.Every(refreshRate).Minutes().Do(runBoards)
	// <-ch
}

func runBoards() {
	boardId := viper.GetString("trello.boardId")
	Run(boardId)
}

func GetBoard(id string) error {
	board, err := client.GetBoard(id)
	if err != nil {
		log.Println("Get Board err:", err)
		return err
	}
	log.Println(board.ID)
	return nil
}

// Run fetches and saves the points of a given board.
func Run(boardID string) {
	log.Printf("Checking board ID %s", boardID)
	board, err := client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		log.Printf("Couldn't fetch board: %s", err)
		return
	}
	log.Printf("Board name: %s", board.Name)
	lastListID, err := getDoneList(board)
	if err != nil {
		log.Printf("Couldn't fetch last list: %s", err)
	}
	resultChannel := make(chan *cardResult)
	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Printf("Couldn't fetch cards: %s", err)
	}
	for _, card := range cards {
		go determineCardComplete(card, lastListID, resultChannel)
	}
	boardEntity := Board{
		ID:   boardID,
		Name: board.Name,
	}
	var pointsPerDay = make(map[string]float64)
	for i := 0; i < len(cards); i++ {
		response := <-resultChannel
		if response.Error != nil {
			log.Fatalln(response.Error)
		}
		if response.Complete {
			boardEntity.CardsCompleted++
			boardEntity.PointsCompleted += response.Points
			if _, ok := pointsPerDay[response.Date]; ok {
				pointsPerDay[response.Date] = response.Points + pointsPerDay[response.Date]
			} else {
				pointsPerDay[response.Date] = response.Points
			}
		}
		boardEntity.Cards++
		boardEntity.Points += response.Points
	}
	log.Printf("Cards progress: %d/%d", boardEntity.CardsCompleted, boardEntity.Cards)
	log.Printf("Total points: %f/%f", boardEntity.PointsCompleted, boardEntity.Points)
}

func getDoneList(board *trello.Board) (string, error) {
	var listID string
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return "", err
	}
	for _, list := range lists {
		if strings.Contains(list.Name, "Done") {
			listID = list.ID
			break
		}
	}
	return listID, nil
}

func determineCardComplete(card *trello.Card, listID string, res chan *cardResult) {
	points := getPoints(card)
	if card.IDList != listID {
		res <- &cardResult{
			Complete: false,
			Points:   points,
		}
		return
	}
	actions, err := card.GetActions(trello.Defaults())
	if err != nil {
		res <- &cardResult{
			Error: err,
		}
		return
	}
	date := card.DateLastActivity
	for _, action := range actions {
		if action.Data.ListAfter != nil && action.Data.ListBefore != nil &&
			action.Data.ListAfter.ID != action.Data.ListBefore.ID && action.Data.ListAfter.ID == listID {
			date = &action.Date
			break
		}
	}
	res <- &cardResult{
		Complete: true,
		Date:     date.Format("2006-01-02"),
		Points:   points,
	}
}

func getPoints(card *trello.Card) float64 {
	r := regexp.MustCompile(`\(([0-9]*\.[0-9]+|[0-9]+)\)`)
	matches := r.FindStringSubmatch(card.Name)
	if len(matches) != 2 {
		return 0
	}
	points, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		log.Fatalln(err)
	}
	return points
}
