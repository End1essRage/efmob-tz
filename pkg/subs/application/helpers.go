package application

import (
	"time"

	"github.com/end1essrage/efmob-tz/pkg/subs/domain"
)

func Periods(sf, st, ef, et *time.Time) (*domain.Period, *domain.Period, error) {
	var startPeriod *domain.Period
	if sf != nil || st != nil {
		var err error
		startPeriod, err = domain.NewPeriod(sf, st)
		if err != nil {
			return nil, nil, err
		}
	}

	var endPeriod *domain.Period
	if ef != nil || et != nil {
		var err error
		endPeriod, err = domain.NewPeriod(ef, et)
		if err != nil {
			return nil, nil, err
		}
	}

	return startPeriod, endPeriod, nil
}
