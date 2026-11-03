package domain

type XraySyncChange struct {
	InboundTag string
	Email      string
}

type XraySyncReport struct {
	DryRun  bool
	Added   []XraySyncChange
	Removed []XraySyncChange
	Errors  []error
}

func (r XraySyncReport) IsEmpty() bool {
	return len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Errors) == 0
}
