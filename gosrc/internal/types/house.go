package types

type House struct {
	Name    string       `json:"name"`
	Amount  StringNumber `json:"amount"`
	Status  string       `json:"status"`
	PayTime string       `json:"payTime"`
	Qing    StringNumber `json:"qing"`
	Yang    StringNumber `json:"yang"`
}

type StringNumber string

func (s *StringNumber) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		*s = StringNumber(data[1 : len(data)-1])
	} else {
		*s = StringNumber(data)
	}
	return nil
}
