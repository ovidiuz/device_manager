package domain

type Measurement struct {
	Field       string `json:"_field"`
	Measurement string `json:"_measurement"`
	Start       string `json:"_start"`
	Stop        string `json:"_stop"`
	Time        string `json:"_time"`
	Value       int    `json:"_value"`
	Host        string `json:"host"`
	Result      string `json:"result"`
	Table       int    `json:"table"`
	Topic       string `json:"topic"`
}
