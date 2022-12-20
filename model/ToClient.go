package model

type ToClient struct {
	Donation Donations `json:"donation"`
	AudioUrl string    `json:"audio_url"`
	Audio    []byte    `json:"audio"`
}
