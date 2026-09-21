package domain

type Duration struct {
	Month        int32
	Days         int32
	Microseconds int64
}

func NewDuration(month int32, days int32, microseconds int64) Duration {
	return Duration{
		Month:        month,
		Days:         days,
		Microseconds: microseconds,
	}
}
