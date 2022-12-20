package model

type Donations struct {
	From        string `json:"from,omitempty"`
	Amount      int    `json:"amount,omitempty"`
	UserId      int    `json:"user_id,omitempty"`
	Email       string `json:"email,omitempty"`
	Message     string `json:"message,omitempty"`
	IsAnonymous bool   `json:"is_anonymous,omitempty"`
}
