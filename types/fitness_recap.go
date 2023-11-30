package types

import "time"

//TODO: Think about just having runs and weigh ins to make it easier to have multiple of each in a day.

// FitnessRecap stores weight and distance as integers that need to be turned into floats then divided.
// We do it this way because we are going to be adding lots of them together and don't want to deal with
// the inaccuracies of floating points while doing that.
type FitnessRecap struct {
	HundredthsOfAMile int
	TenthsOfAPound    int
	Date              time.Time
}

func (f *FitnessRecap) GetDistance() float64 {
	return float64(f.HundredthsOfAMile) / 100
}

func (f *FitnessRecap) GetWeight() float64 {
	return float64(f.TenthsOfAPound) / 10
}
