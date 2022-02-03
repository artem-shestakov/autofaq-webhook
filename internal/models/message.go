package models

type Messages struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Id             string   `json:"id"`
	Dt             string   `json:"dt"`
	SessionId      string   `json:"sessionId"`
	ConversationId string   `json:"conversationId"`
	Text           string   `json:"text"`
	Documents      []string `json:"documents"`
	Buttons        []Button `json:"buttons"`
	Operator       string   `json:"operator"`
	Files          []string `json:"files"`
	ChannelUserId  string   `json:"channelUserId"`
}

type Button struct {
	Text    string `json:"text"`
	Payload string `json:"payload"`
}
