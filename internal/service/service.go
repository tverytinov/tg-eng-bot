package service

import (
	"fmt"
	"regexp"

	"github.com/tverytinov/tg-eng-bot/internal/client"
	"github.com/tverytinov/tg-eng-bot/internal/model"
)

type Service struct {
	tgClient *client.TgClient
	aiClient *client.AIClient
	offset   int
}

func NewService(tgClient *client.TgClient, aiClient *client.AIClient) *Service {
	return &Service{
		tgClient: tgClient,
		aiClient: aiClient,
	}
}

func (s *Service) Fetch(limit int) ([]model.Update, error) {
	updates, err := s.tgClient.Update(limit, s.offset)
	if err != nil {
		return nil, fmt.Errorf("error s.tgClient.Update(): %w", err)
	}

	if len(updates) > 0 {
		s.offset = updates[len(updates)-1].ID + 1
	}

	for i := 0; i <= len(updates)-1; i++ {
		answer, err := s.aiClient.QuestionAI(prepareQuestion(updates[i].Message.Text))
		if err != nil {
			return nil, fmt.Errorf("error s.aiClient.QuestionAI(): %w", err)
		}

		updates[i].Message.Text = answer
	}

	return updates, nil
}

func (s *Service) Process(update model.Update) error {
	if err := s.tgClient.SendMessage(update.Message.Chat.ID, update.Message.Text); err != nil {
		return fmt.Errorf("error s.tgClient.SendMessage(): %w", err)
	}

	return nil
}

func prepareQuestion(text string) string {
	queryEng := `
	1.Как это переводится "%s" с английского на русский?
	2.Какой смысл у этого слова?
	3.Продемонстрируй 5 примеров на английском с этим словом, а ниже укажи перевод на русском.
	`
	queryRus := `
	1.Как это переводится "%s" с русского на английский?
	2.Продемонстрируй 5 примеров на английском с этим словом, а ниже укажи перевод на русском.
	`

	latinRegex := regexp.MustCompile(`^[a-zA-Z]+(?:\s[a-zA-Z]+)*$`)
	cyrillicRegex := regexp.MustCompile(`^[а-яА-Я]+(?:\s[а-яА-Я]+)*$`)

	if latinRegex.MatchString(text) {
		return fmt.Sprintf(queryEng, text)
	}

	if cyrillicRegex.MatchString(text) {
		return fmt.Sprintf(queryRus, text)
	}

	return `Incorrect input`
}
