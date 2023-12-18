package trello_service

import (
	"fmt"
	"go-trello/logger"

	"github.com/adlio/trello"
)

var nameOfMembers = map[string]string{
    "nguyenanhtuan010598": "tuanna137",
    "phatlechi": "phatlc",
    "bobochacha07": "binhtd7",
	"kuroka1":"dongnt18",
	"dokieuminhdkm": "minhdk",
	"hiunguyn484":"hieund152",
	"namnguyen08489158":"namnp16",
	"thinhnx5": "thinhnx5",
	"xuanbachtrng": "bachtx",
	"hungto18":"hungto18",
	"maingoctrinh":"trinhmn",
}

// GetMembersInBoard returns a list of members of the board
func (c *TrelloClient) GetMembersInBoard() (members []*trello.Member, err error) {
	if c.Board == nil {
		return nil, fmt.Errorf("no board specified, get board first")
	}
	members, err = c.Board.GetMembers()
	if err != nil {
		return nil, err
	}

	logger.Debug("Init members statistics")
	for _, m := range members {
		m.SetClient(c.Client)
		c.MemberStats[m.ID] = &MemberStats{
			Email:    m.Email,
			Name:     m.Username,
			FullName: m.FullName,
			// TypeOfTask: m.TypeOfTask,
		}
		c.Members[m.ID] = m
	}
	c.DailyTrackingStats.InitMembersDailyTracking(c.MemberStats)

	return members, nil
}

// GetMemberCard returns a list of actions of the member
func (c *TrelloClient) GetMemberActions(id string, args trello.Arguments) (action []*trello.Action, err error) {
	if c.Board == nil {
		return nil, fmt.Errorf("no board specified, get board first")
	}

	// check member already exists in board
	if c.Members == nil || len(c.Members) == 0 {
		return nil, fmt.Errorf("no members specified, get board first")
	}

	if ValidateMember(id) {
		action := []*trello.Action{}
		err = c.Client.Get(fmt.Sprintf("member/%v/actions", id), args, &action)
		if err != nil {
			return nil, err
		}

		return action, nil
	}
	return nil, fmt.Errorf("no valid member, check board first")
}

// ValidateMember checks if member exists in board
func ValidateMember(id string) bool {
	if c.Members == nil || len(c.Members) == 0 {
		return false
	}

	for _, member := range c.Members {
		if member.ID == id {
			return true
		}
	}
	return false
}

// FilterCardByMember gets cards list by member
func FilterCardByMember(id string) (cards []*trello.Card, err error) {
	if c.Board == nil {
		return nil, fmt.Errorf("no specified board, check board first")
	}

	if !ValidateMember(id) {
		return nil, fmt.Errorf("no valid member")
	}

	return
}

func ConvertNameOfMembers(memberData *TrelloClient) {
	for nameMemberInEmail, nameMemberInCompany := range nameOfMembers {
		for memberId, _ := range memberData.Members {
			nameEachMemberInCompany := memberData.MemberStats[memberId].Name
			if nameMemberInEmail == nameEachMemberInCompany {
				memberData.MemberStats[memberId].Name = nameMemberInCompany
			}	
		}
	}
}


func ConvertNameOfMembersInLinkedList(memberStat *MemberStats) {
	for nameMemberInEmail, nameMemberInCompany := range nameOfMembers {
		if nameMemberInEmail == memberStat.Name {
			memberStat.Name = nameMemberInCompany
		}	
	}
}