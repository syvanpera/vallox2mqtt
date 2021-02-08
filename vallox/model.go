package vallox

type State struct {
	FanSpeed     int `json:"fan_speed"`     // 0x29
	OutsideTemp  int `json:"outside_temp"`  // 0x32
	InsideTemp   int `json:"inside_temp"`   // 0x34
	ExhaustTemp  int `json:"exhaust_temp"`  // 0x33
	IncomingTemp int `json:"incoming_temp"` // 0x35
}
