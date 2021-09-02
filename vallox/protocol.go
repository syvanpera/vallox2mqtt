package vallox

const (
	vxMsgLength = 6 // Always 6 bytes
	vxDomain    = 1 // Always 1
)

// Addresses for sender and receiver
const (
	AddressMainboards = iota + 0x10
	AddressMaster
)

const (
	AddressPanels = iota + 0x20
	AddressPanel1
	AddressPanel2
	AddressPanel3
	AddressPanel4
	AddressPanel5
	AddressPanel6
	AddressPanel7

	AddressLON
	AddressPanel8
)

// Variables
const (
	VariableFanSpeed    = 0x29
	VariableTempOutside = 0x32
	VariableTempExhaust = 0x33
	VariableTempInside  = 0x34
	VariableTempSupply  = 0x35
)

// 01H = speed 1
// 03H = speed 2
// 07H = speed 3
// 0FH = speed 4
// 1FH = speed 5
// 3FH = speed 6
// 7FH = speed 7
// FFH = speed 8
// TODO Maybe convert these to actual Maps?
var fanSpeedMap = [...]int{0x01, 0x03, 0x07, 0x0F, 0x1F, 0x3F, 0x7F, 0xFF}
var temperatureMap = [...]int{
	-74, -70, -66, -62, -59, -56, -54, -52,
	-50, -48, -47, -46, -44, -43, -42, -41,
	-40, -39, -38, -37, -36, -35, -34, -33,
	-33, -32, -31, -30, -30, -29, -28, -28,
	-27, -27, -26, -25, -25, -24, -24, -23,
	-23, -22, -22, -21, -21, -20, -20, -19,
	-19, -19, -18, -18, -17, -17, -16, -16,
	-16, -15, -15, -14, -14, -14, -13, -13,
	-12, -12, -12, -11, -11, -11, -10, -10,
	-9, -9, -9, -8, -8, -8, -7, -7,
	-7, -6, -6, -6, -5, -5, -5, -4,
	-4, -4, -3, -3, -3, -2, -2, -2,
	-1, -1, -1, -1, 0, 0, 0, 1,
	1, 1, 2, 2, 2, 3, 3, 3,
	4, 4, 4, 5, 5, 5, 5, 6,
	6, 6, 7, 7, 7, 8, 8, 8,
	9, 9, 9, 10, 10, 10, 11, 11,
	11, 12, 12, 12, 13, 13, 13, 14,
	14, 14, 15, 15, 15, 16, 16, 16,
	17, 17, 18, 18, 18, 19, 19, 19,
	20, 20, 21, 21, 21, 22, 22, 22,
	23, 23, 24, 24, 24, 25, 25, 26,
	26, 27, 27, 27, 28, 28, 29, 29,
	30, 30, 31, 31, 32, 32, 33, 33,
	34, 34, 35, 35, 36, 36, 37, 37,
	38, 38, 39, 40, 40, 41, 41, 42,
	43, 43, 44, 45, 45, 46, 47, 48,
	49, 49, 50, 51, 52, 53, 53, 54,
	55, 56, 57, 59, 60, 61, 62, 63,
	65, 66, 68, 69, 71, 73, 75, 77,
	79, 81, 82, 86, 90, 93, 97, 100,
	100, 100, 100, 100, 100, 100, 100, 100,
}

func convertFanSpeed(value byte) int {
	fanSpeed := 0

	for i := 0; i < len(fanSpeedMap); i++ {
		if fanSpeedMap[i] == int(value) {
			fanSpeed = i + 1
			break
		}
	}

	return fanSpeed
}

func convertTemperature(value byte) int {
	return temperatureMap[value]
}
