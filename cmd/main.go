package main

import (
	"flag"
	"log"

	"github.com/tverytinov/tg-eng-bot/internal/client"
	"github.com/tverytinov/tg-eng-bot/internal/consumer"
	"github.com/tverytinov/tg-eng-bot/internal/service"

	"github.com/sirupsen/logrus"
)

var (
	tgHost         = "api.telegram.org"
	tgApiToken     = ""
	openAIApiToken = ""
	defaultLimit   = 100
)

func init() {
	flag.StringVar(&tgApiToken, "tg-api-token", "", "telegram api token")

	flag.StringVar(&openAIApiToken, "open-ai-api-token", "", "open ai api token")

	flag.Parse()

	if tgApiToken == "" || openAIApiToken == "" {
		log.Fatal("some tokens are invalid or missing")
	}
}

func main() {
	log := logrus.New()

	tgClient := client.NewClient(tgHost, tgApiToken)

	aiClient := client.NewAIClient(openAIApiToken)

	service := service.NewService(tgClient, aiClient)

	consumer := consumer.NewConsumer(log, service)

	log.Info("Server started")

	consumer.Start()
}
