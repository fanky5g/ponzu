package types

type Author struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Image string `json:"image" reference:"Upload"`
}
