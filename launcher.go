package launcher

import (
	"fmt"
	"time"
)

const (
	MaxDuration     = 5 * time.Second
	DefaultDuration = MaxDuration
)

type Launcher interface {
	Launch() error
}

type LaunchProcess interface {
	Launcher
	Name() string
	TimeDuration() time.Duration
}

type Dependency interface {
	Add(launcher LaunchProcess)
}

type LaunchProcessor struct {
	Launcher
	Dependency
	dependencies []LaunchProcess
	alertSender  AlertSender
}

func New() *LaunchProcessor {
	return &LaunchProcessor{}
}

func (p *LaunchProcessor) SetAlertSender(a AlertSender) {
	if p.alertSender == nil {
		p.alertSender = a
	}
}

func (p *LaunchProcessor) Add(dependency LaunchProcess) {
	p.dependencies = append(p.dependencies, dependency)
}

func (p *LaunchProcessor) AddFunc(name string, f func() (err error)) {
	p.Add(&ProcessFunc{name: name, launchFunc: f})
}

func (p *LaunchProcessor) Launch() {
	for _, dep := range p.dependencies {
		p.launchDependency(dep)
	}
}

func (p *LaunchProcessor) launchDependency(l LaunchProcess) {
	for {
		if err := p.lockedLaunch(l); err != nil {
			p.error(fmt.Sprintf("launch %s: %s", l.Name(), err.Error()))
			time.Sleep(l.TimeDuration())
		} else {
			break
		}
	}
}

func (p *LaunchProcessor) lockedLaunch(l LaunchProcess) (err error) {
	done := make(chan struct{}, 1)
	go func() {
		err = l.Launch()
		defer func() { done <- struct{}{} }()
	}()
	tick := time.NewTicker(l.TimeDuration())
	go func() {
		name, i := l.Name(), 1
		for range tick.C {
			p.error(fmt.Sprintf("launch %s: timeout: iterate %d", name, i))
			i++
		}
	}()
	defer func() { <-done; tick.Stop() }()
	return
}

func (p *LaunchProcessor) error(msg string) {
	if p.alertSender != nil {
		p.alertSender.Send(msg)
	}
}
