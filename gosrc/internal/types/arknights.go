package types

type Data struct {
	Chars []Chars `json:"chars"`
	Pool  string  `json:"pool"`
	Ts    int64   `json:"ts"`
}
type Chars struct {
	IsNew  bool   `json:"isNew"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
}
