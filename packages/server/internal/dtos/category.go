package dtos

type Color struct {
	Background string `json:"background"`
	Foregound  string `json:"foreground"`
}

type Category struct {
	Id       int    `json:"id"`
	Label    string `json:"label"`
	Color    Color  `json:"color"`
	IconName string `json:"iconName"`
}
