package models

type WebhookResponse struct {
	Status bool          `json:"ok"`
	Result WebhookResult `json:"result"`
}

type WebhookResult struct {
	MessageID int    `json:"message_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`

	From struct {
		ID        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	} `json:"from"`

	ReplyToMessage struct {
	} `json:"reply_to_message"`

	Chat struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
}
