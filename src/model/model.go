package model

import "math/rand"

type Model struct {
	Nickname string  `json:"nickname"`
	PhotoUrl string  `json:"photo_url"`
	Text     string  `json:"text"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Value    float64 `json:"value"`
}

func NewModel() *Model {
	return &Model{
		Nickname: "pashaborisyk",
		PhotoUrl: "photoUrl",
		Text:     "Lets party",
		X:        rand.Float64(),
		Y:        rand.Float64(),
		Value:    rand.Float64(),
	}
}
