package launcher

import "time"

type ProcessFunc struct {
	launchFunc func() (err error)
	name       string
}

func (l *ProcessFunc) Launch() error {
	return l.launchFunc()
}

func (l *ProcessFunc) Name() string {
	return l.name
}

func (l *ProcessFunc) TimeDuration() time.Duration {
	return DefaultDuration
}
