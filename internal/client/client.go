package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"voicer/internal/storage"
)

var apiKey = "" //enter your openrouter api-key

type Client struct {
	client *http.Client
}

func NewClient(cl *http.Client) *Client {
	return &Client{client: cl}
}

type ImageGenerationRequest struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	Modalities []string  `json:"modalities"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ImageGenerationResponse struct {
	Choices []struct {
		Message struct {
			Images []struct {
				ImageURL struct {
					URL string `json:"url"`
				} `json:"image_url"`
			} `json:"images"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *Client) GetFlag(s string) ([]byte, error) {
	content := fmt.Sprintf("Generate the official rectangular flag of %s, aspect ratio 3:2 without any decoration(like wave etc.). Extension of this file must be \".png\".", s)
	reqData := ImageGenerationRequest{
		Model: "x-ai/grok-imagine-image-quality",
		Messages: []Message{
			{Role: "user", Content: content},
		},
		Modalities: []string{"image"},
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}

	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}

	var imgResp ImageGenerationResponse
	err = json.Unmarshal(body, &imgResp)
	if err != nil {
		slog.Info(err.Error())
		return nil, err
	}
	var decoded []byte
	if len(imgResp.Choices) > 0 && len(imgResp.Choices[0].Message.Images) > 0 {
		imageUrl := imgResp.Choices[0].Message.Images[0].ImageURL.URL

		if strings.HasPrefix(imageUrl, "data:image/") {
			base64Data := strings.SplitN(imageUrl, ",", 2)[1]
			decoded, err = base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				slog.Info(err.Error())
				return nil, err
			}
			p, err := storage.Path()
			if err != nil {
				slog.Info(err.Error())
			}

			err = os.WriteFile(p+"/flags/"+s+".png", decoded, 0644)
			if err != nil {
				slog.Info(err.Error())
				return nil, err
			}
			slog.Info("Изображение сохранено в %s", p+"/flags/"+s+".png")
		} else {
			slog.Info("Получен неожиданный формат изображения:", imageUrl)
		}
	} else {
		slog.Info("В ответе нет изображений.")
	}

	return decoded, err
}
