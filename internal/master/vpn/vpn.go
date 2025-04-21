package vpn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type VPN interface {
	GetConfig(userID int, ip string, port string, endpoint string) (string, error)
}

type VPNServer struct {
	Protocol string
	Host     string
	APIKey   string
}

func (v *VPNServer) GetConfig(userID int, ip string, port string, endpoint string) (string, error) {
	url := fmt.Sprintf("http://%s:%s%s", ip, port, endpoint)

	requestData := struct {
		UserID int `json:"user_id"`
	}{UserID: userID}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Создаем POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("oшибка создания запроса: %v", err)
	}

	// Устанавливаем заголовки (при необходимости)
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	return string(body), nil
}
