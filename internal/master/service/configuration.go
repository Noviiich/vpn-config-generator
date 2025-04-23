package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func (s *VPNService) GetConfig(ctx context.Context, userID int) (config string, err error) {
	baseURL := "http://localhost:8080/config"
	params := url.Values{}
	params.Add("user_id", strconv.Itoa(userID))
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Создаем GET-запрос
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var response struct {
		Config string `json:"config"`
	}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return response.Config, nil
}
