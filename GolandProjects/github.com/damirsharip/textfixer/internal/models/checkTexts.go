package models

type TextInput struct {
	Text    string `form:"text" binding:"required"`
	Lang    string `form:"lang,omitempty"`
	Options int    `form:"options,omitempty"`
	Format  bool   `form:"format,omitempty"`
}

type TextsInput struct {
	Texts   []string `json:"texts" binding:"required"`
	Lang    string   `form:"lang,omitempty"`
	Options int      `form:"options,omitempty"`
	Format  bool     `form:"format,omitempty"`
}

type YandexResponse []struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

type CorrectorResponse struct {
	Texts []string `json:"texts"`
}
