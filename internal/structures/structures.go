package structures

import "math"

const (
	Q    float64 = 0.283 // related to the gain
	Gain float64 = math.Pi
	Pt   float64 = 0.9         // Power of transmitter
	P_n  float64 = 0.000000001 // variance of noise at the receiver
)

type Coordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Environment struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type RISCHANNL [][]float64

type Updates struct {
	Ris Coordinates
	Rx  Coordinates
	Tx  Coordinates
	Los bool
}
