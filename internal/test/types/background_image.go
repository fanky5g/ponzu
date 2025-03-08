package types

type BackgroundImage struct {
	File     string `json:"file" reference:"Upload"`
	Alt      string `json:"alt"`
	Position string `json:"position"`
}
