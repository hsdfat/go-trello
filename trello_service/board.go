package trello_service

import "go-trello/logger"

// GetBoard get board info by id
func (c *TrelloClient) GetBoard(id string) {
	var err error
	c.Board, err = c.Client.GetBoard(id)
	if err != nil {
		logger.Errorln("Get board err")
	}
	c.Board.SetClient(c.Client)
}
