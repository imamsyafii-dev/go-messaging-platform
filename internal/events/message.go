package events

type MessageCreated struct {
	EventType        string   `json:"event_type"`
	MessageID        string   `json:"message_id"`
	ConversationID   string   `json:"conversation_id"`
	SenderID         string   `json:"sender_id"`
	ClientMessageID  string   `json:"client_message_id"`
	MessageType      string   `json:"message_type"`
	Content          string   `json:"content,omitempty"`
	RecipientUserIDs []string `json:"recipient_user_ids"`
}
