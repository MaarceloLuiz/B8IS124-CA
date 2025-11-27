package geocalc

import "math"

//a = sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)
//c = 2 ⋅ atan2( √a, √(1−a) )
//d = R ⋅ c
func haversine(guessLat, guessLng, answerLat, answerLng float64) (float64, error) {
	const R = 6371e3 //radius of the Earth in meters

	guessLatRad := guessLat * (math.Pi / 180)
	answerLatRad := answerLat * (math.Pi / 180)
	deltaLat := (answerLat - guessLat) * (math.Pi / 180)
	deltaLng := (answerLng - guessLng) * (math.Pi / 180)

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(guessLatRad)*math.Cos(answerLatRad)*math.Pow(math.Sin(deltaLng/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c //distance in meters

	return d, nil
}
