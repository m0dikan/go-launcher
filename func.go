package launcher

import "time"

type ProcessFunc struct {
	launchFunc   func() (err error)
	name         string
	timeDuration time.Duration
}

func (l *ProcessFunc) Launch() error {
	return l.launchFunc()
}

func (l *ProcessFunc) Name() string {
	return l.name
}

func (l *ProcessFunc) TimeDuration() time.Duration {
	if l.timeDuration == 0 {
		l.timeDuration = DefaultDuration
	}
	return l.timeDuration
}
