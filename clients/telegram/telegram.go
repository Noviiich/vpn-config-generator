package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/Noviiich/vpn-config-generator/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
	adminID  int
}

const (
	getUpdatesMethod   = "getUpdates"
	sendMessageMethod  = "sendMessage"
	sendDocumentMethod = "sendDocument"
)

func New(host string, token string, adminID int) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
		adminID:  adminID,
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(ctx context.Context, offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(ctx, getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMarkdown(ctx context.Context, chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("parse_mode", "Markdown")

	_, err := c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send markdown", err)
	}

	return nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendDocument(ctx context.Context, chatID int, text string, fileName string) (err error) {
	defer func() { err = e.WrapIfErr("can't send document: %v", err) }()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляю поле chat_id
	if err := writer.WriteField("chat_id", strconv.Itoa(chatID)); err != nil {
		return err
	}

	part, err := writer.CreateFormFile("document", fileName)
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, strings.NewReader(text)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	_, err = c.doRequestDocument(ctx, sendDocumentMethod, body, writer.FormDataContentType())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) doRequestDocument(ctx context.Context, method string, body io.Reader, contentType string) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) DeleteApprovalButtons(ctx context.Context, messageID int) (err error) {
	defer func() { err = e.WrapIfErr("can't delete approval buttons", err) }()

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{},
	}

	markupJSON, err := json.Marshal(replyMarkup)
	if err != nil {
		return err
	}

	// Формируем параметры запроса
	formData := url.Values{}
	formData.Set("chat_id", strconv.Itoa(c.adminID))
	formData.Set("message_id", strconv.Itoa(messageID))
	formData.Set("reply_markup", string(markupJSON))

	// Отправляем запрос на редактирование сообщения
	_, err = c.doRequest(ctx, "editMessageReplyMarkup", formData)
	if err != nil {
		return err
	}

	return nil

}

func (c *Client) SendApprovalButtons(ctx context.Context, text string, userID int, tariff string) (err error) {
	defer func() { err = e.WrapIfErr("can't send approval buttons", err) }()

	approveData := fmt.Sprintf("approve_%d_%s", userID, tariff)
	rejectData := fmt.Sprintf("reject_%d_%s", userID, tariff)

	replyMarkup := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{
					"text":          "✅ Одобрить",
					"callback_data": approveData,
				},
				{
					"text":          "❌ Отклонить",
					"callback_data": rejectData,
				},
			},
		},
	}

	markupJSON, err := json.Marshal(replyMarkup)
	if err != nil {
		return err
	}

	formData := url.Values{}
	formData.Set("chat_id", strconv.Itoa(c.adminID))
	formData.Set("text", text)
	formData.Set("reply_markup", string(markupJSON))

	_, err = c.doRequest(ctx, sendMessageMethod, formData)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) NotifyUserSubscriptionApproved(ctx context.Context, chatID int) error {
	text := "✅ Ваша заявка на подписку одобрена! Теперь вы можете использовать VPN."
	return c.SendMessage(ctx, chatID, text)
}

func (c *Client) NotifyUserSubscriptionRejected(ctx context.Context, chatID int) error {
	text := "❌ К сожалению, ваша заявка на подписку отклонена. Пожалуйста, свяжитесь с администратором для получения дополнительной информации."
	return c.SendMessage(ctx, chatID, text)
}

func (c *Client) SendMessageWithKeyboard(ctx context.Context, chatID int, text string, buttons [][]string) error {
	keyboard := map[string]interface{}{
		"keyboard":        buttons,
		"resize_keyboard": true,
	}

	keyboardJSON, err := json.Marshal(keyboard)
	if err != nil {
		return e.Wrap("can't marshal keyboard", err)
	}

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("reply_markup", string(keyboardJSON))

	_, err = c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message with keyboard", err)
	}

	return nil
}

func (c *Client) SendMessageWithTariffButtons(ctx context.Context, chatID int, text string) error {
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{
					"text":          "Базовый - 50₽/мес",
					"callback_data": "tariff_basic",
				},
			},
			{
				{
					"text":          "Стандарт - 130₽/3 мес",
					"callback_data": "tariff_standard",
				},
			},
			{
				{
					"text":          "Премиум - 240₽/6 мес",
					"callback_data": "tariff_premium",
				},
			},
		},
	}

	keyboardJSON, err := json.Marshal(keyboard)
	if err != nil {
		return e.Wrap("can't marshal keyboard", err)
	}

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("reply_markup", string(keyboardJSON))

	_, err = c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message with tariff buttons", err)
	}

	return nil
}

func (c *Client) SendMessageWithProtocolButtons(ctx context.Context, chatID int, text string, buttons [][]string) error {
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{
					"text":          "WireGuard",
					"callback_data": "protocol_wireguard",
				},
			},
			{
				{
					"text":          "OpenVPN",
					"callback_data": "protocol_openvpn",
				},
			},
			// {
			// 	{
			// 		"text":          "IKEv2",
			// 		"callback_data": "protocol_ikev2",
			// 	},
			// },
		},
		"keyboard":        buttons,
		"resize_keyboard": true,
	}

	keyboardJSON, err := json.Marshal(keyboard)
	if err != nil {
		return e.Wrap("can't marshal keyboard", err)
	}

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("reply_markup", string(keyboardJSON))

	_, err = c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message with protocol buttons", err)
	}

	return nil
}
