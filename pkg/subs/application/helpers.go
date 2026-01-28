package application

import (
	"time"

	"github.com/end1essrage/efmob-tz/pkg/subs/domain"
)

func period(f, t *time.Time) (*domain.Period, error) {
	from := time.Now()
	to := time.Now()
	if f != nil {
		from = *f
	}
	if t != nil {
		to = *t
	}

	return domain.NewPeriod(from, to)
}

func Periods(sf, st, ef, et *time.Time) (*domain.Period, *domain.Period, error) {
	var startPeriod *domain.Period
	if sf != nil && st != nil {
		var err error
		startPeriod, err = period(sf, st)
		if err != nil {
			return nil, nil, err
		}
	}

	var endPeriod *domain.Period
	if ef != nil && et != nil {
		var err error
		endPeriod, err = period(ef, et)
		if err != nil {
			return nil, nil, err
		}
	}

	return startPeriod, endPeriod, nil
}
