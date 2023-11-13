package types

type Rating int

const (
	_ Rating = iota
	Poor
	Average
	Good
	VeryGood
	Excellent
)

func (r Rating) String() string {
	switch r {
	case Poor:
		return "Poor"
	case Average:
		return "Average"
	case Good:
		return "Good"
	case VeryGood:
		return "Very Good"
	case Excellent:
		return "Excellent"
	default:
		return "Unknown Rating"
	}
}
