package vallox

type State struct {
	FanSpeed        int `json:"fan_speed"`         // 0x29
	DefaultFanSpeed int `json:"default_fan_speed"` // 0xA9
	OutsideTemp     int `json:"outside_temp"`      // 0x32
	InsideTemp      int `json:"inside_temp"`       // 0x34
	ExhaustTemp     int `json:"exhaust_temp"`      // 0x33
	IncomingTemp    int `json:"incoming_temp"`     // 0x35

	ServicePeriod int // 0xA6
	NextService   int // 0xAB

	// Indicator lights
	PowerOn   bool // 0xA3 bit 0
	HeatOn    bool // 0xA3 bit 3
	FaultOn   bool // 0xA3 bit 6
	ServiceOn bool // 0xA3 bit 7
}
