package entities

type BotMessage struct {
	Command string `json:"command"`
	Value   string `json:"value"`
}
