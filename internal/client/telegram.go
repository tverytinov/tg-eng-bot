package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tverytinov/tg-eng-bot/internal/model"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type TgClient struct {
	tgHost     string
	tgBasePath string
	client     http.Client
}

func NewClient(host, token string) *TgClient {
	return &TgClient{
		tgHost:     host,
		tgBasePath: newBasePath(token),
		client:     http.Client{},
	}
}

func (c *TgClient) Update(limit, offset int) ([]model.Update, error) {
	query := fmt.Sprintf(`https://%v/%v/%v?limit=%v&offset=%v`, c.tgHost, c.tgBasePath, getUpdatesMethod, limit, offset)

	req, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		return nil, fmt.Errorf("error http.NewRequest(): %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error c.client.Do(): %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var res model.Updates

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("error json.NewDecoder().Decode(): %w", err)
	}

	return res.Result, nil
}

func (c *TgClient) SendMessage(chatID int, text string) error {
	query := fmt.Sprintf(`https://%v/%v/%v`, c.tgHost, c.tgBasePath, sendMessageMethod)

	params := model.SendMessage{
		ChatID: strconv.Itoa(chatID),
		Text:   text,
	}

	bytedParams, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("error json.Marshal(): %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, query, bytes.NewBuffer(bytedParams))
	if err != nil {
		return fmt.Errorf("error http.NewRequest(): %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error c.client.Do(): %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("%s", resp.Status)
}

func newBasePath(token string) string {
	return "bot" + token
}
