package types

type ButtonLink struct {
	Type string       `json:"type"`
	Text string       `json:"text"`
	Link LinkWithType `json:"link"`
}

type LinkWithType struct {
	Href string `json:"href"`
	Type string `json:"type"`
}
