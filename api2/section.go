package api2

type Section struct {
	Categories []*Category `json:"categories"`
	Channels   []*Channel  `json:"channels"`
	Radios     []*Radio    `json:"radios"`
}
