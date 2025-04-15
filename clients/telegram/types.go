package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
	// CallbackQuery *CallbackQuery   `json:"callback_query"`
}

// type CallbackQuery struct {
// 	ID   string          `json:"id"`
// 	From From            `json:"from"`
// 	Message *IncomingMessage `json:"message"`
// 	Data string          `json:"data"`
// }

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
