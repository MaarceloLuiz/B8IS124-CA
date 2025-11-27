package geocalc

import (
	"math"

	"github.com/sirupsen/logrus"
)

func getCompass(guessLat, guessLng, answerLat, answerLng float64) string {
	bearing := initialBearing(guessLat, guessLng, answerLat, answerLng)
	degrees := radiansToDegrees(bearing)
	normalizedDegrees := normalizeDegrees(degrees)
	compass := degreesToCompass(normalizedDegrees)
	if compass == "" {
		logrus.Fatalf("Invalid compass direction for degrees: %f", normalizedDegrees)
	}

	return compass
}

func degreesToCompass(degrees float64) string {
	directions := []string{
		"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW",
	}
	idx := int((degrees+11.25)/22.5) % 16
	return directions[idx]
}

// When calculate the initial bearing using the atan2 function, the result (after converting to degrees)
// can be any value between -180° and +180°. However, compass bearings are always expressed as 0° to 360°.
// normalizeDegrees adjusts an angle to ensure it falls within the range of 0° to 360°.
func normalizeDegrees(degrees float64) float64 {
	return math.Mod(degrees+360, 360)
}

func radiansToDegrees(bearing float64) float64 {
	return bearing * (180.0 / math.Pi)
}

// θ = atan2( sin Δλ ⋅ cos φ2 , cos φ1 ⋅ sin φ2 − sin φ1 ⋅ cos φ2 ⋅ cos Δλ )
func initialBearing(guessLat, guessLng, answerLat, answerLng float64) float64 {
	guessLatRad := guessLat * (math.Pi / 180)
	guessLngRad := guessLng * (math.Pi / 180)
	answerLatRad := answerLat * (math.Pi / 180)
	answerLngRad := answerLng * (math.Pi / 180)
	deltaLng := answerLngRad - guessLngRad

	y := math.Sin(deltaLng) * math.Cos(answerLatRad)
	x := math.Cos(guessLatRad)*math.Sin(answerLatRad) - math.Sin(guessLatRad)*math.Cos(answerLatRad)*math.Cos(answerLngRad-guessLngRad)
	bearing := math.Atan2(y, x)

	return bearing
}
