package types

type Banner struct {
	Text       string          `json:"text"`
	Background BackgroundImage `json:"background"`
	Cta        []ButtonLink    `json:"cta"`
}
