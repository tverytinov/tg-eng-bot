package consumer

import (
	"fmt"

	"github.com/tverytinov/tg-eng-bot/internal/model"
	"github.com/tverytinov/tg-eng-bot/internal/service"

	"github.com/sirupsen/logrus"
)

type Consumer struct {
	log          *logrus.Logger
	service      *service.Service
	defaultLimit int
}

func NewConsumer(log *logrus.Logger, serv *service.Service) *Consumer {
	return &Consumer{
		log:     log,
		service: serv,
	}
}

func (c *Consumer) Start() {
	for {
		updates, err := c.service.Fetch(c.defaultLimit)
		if err != nil {
			c.log.Errorf("error c.service.Fetch(): %v", err.Error())

			continue
		}

		if len(updates) == 0 {

			continue
		}

		if err := c.Iterator(updates); err != nil {
			c.log.Errorf("error c.Iterator(): %v", err.Error())

			continue
		}

	}
}

func (c *Consumer) Iterator(updates []model.Update) error {
	for _, upd := range updates {
		if err := c.service.Process(upd); err != nil {
			return fmt.Errorf("error c.service.Process(): %w", err)
		}
	}

	return nil
}
