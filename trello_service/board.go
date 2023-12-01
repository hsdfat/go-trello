package trello_service

import "go-trello/logger"

// GetBoard get board info by id
func (c *TrelloClient) GetBoard(id string) {
	var err error
	c.CBoard, err = c.Client.GetBoard(id)
	if err != nil {
		logger.Errorln("Get board err")
	}
	c.CBoard.SetClient(c.Client)
}
