package domain

type Duration struct {
	Months       int32
	Days         int32
	Microseconds int64
}

func NewDuration(month int32, days int32, microseconds int64) Duration {
	return Duration{
		Months:       month,
		Days:         days,
		Microseconds: microseconds,
	}
}
