package launcher

import (
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type MockedObjectProcess struct {
	mock.Mock
	attempt int
	times   int
}

func (m *MockedObjectProcess) Launch() error {
	if m.attempt >= m.times {
		return nil
	} else {
		m.attempt++
		return m.Mock.Called().Error(0)
	}
}

func (m *MockedObjectProcess) Name() string {
	return m.Mock.Called().String(0)
}

func (m *MockedObjectProcess) TimeDuration() time.Duration {
	return m.Mock.Called().Get(0).(time.Duration)
}

func TestLaunchProcessor_Launch(t *testing.T) {
	proc := func(procName string, duration time.Duration, result error, times int) LaunchProcess {
		m := new(MockedObjectProcess)
		m.times = times
		m.On("Name").Return(procName)
		m.On("TimeDuration").Return(duration)
		m.On("Launch").Return(result).Times(times)
		return LaunchProcess(m)
	}
	type fields struct {
		Launcher     Launcher
		Dependency   Dependency
		dependencies []LaunchProcess
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "4 process",
			fields: fields{
				Launcher:   nil,
				Dependency: nil,
				dependencies: []LaunchProcess{
					proc("postgres master", time.Microsecond, nil, 1),
					proc("postgres slave", time.Millisecond*10, errors.New("test"), 3),
					proc("sync", time.Microsecond, nil, 1),
					proc("rabbit", time.Microsecond, nil, 1),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &LaunchProcessor{
				Launcher:     tt.fields.Launcher,
				Dependency:   tt.fields.Dependency,
				dependencies: tt.fields.dependencies,
			}
			p.Launch()
		})
	}
}

type exampleLaunchProcessMock struct {
	err             error
	name            string
	duration        time.Duration
	durationProcess time.Duration
}

func (s *exampleLaunchProcessMock) Launch() error {
	time.Sleep(s.durationProcess)
	return s.err
}

func (s *exampleLaunchProcessMock) Name() string {
	return s.name
}

func (s *exampleLaunchProcessMock) TimeDuration() time.Duration {
	return s.duration
}

func TestLaunchProcessorTimeouts_Launch(t *testing.T) {
	type args struct {
		l LaunchProcess
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "with error",
			args: args{
				&exampleLaunchProcessMock{
					err:             errors.New("error1"),
					name:            "process1",
					duration:        time.Second,
					durationProcess: time.Microsecond,
				},
			},
			wantErr: true,
		},
		{
			name: "no error",
			args: args{
				&exampleLaunchProcessMock{
					err:             nil,
					name:            "process2",
					duration:        time.Second,
					durationProcess: time.Microsecond,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := tt.args.l.Launch(); (err != nil) != tt.wantErr {
			t.Errorf("test (%s) is fail", tt.name)
		}
	}
}

type BetaAlertSender struct {
	AlertSender
}

func (a *BetaAlertSender) Send(msg ...interface{}) {
	logrus.Error(msg...)
}

type BetaLauncher struct{}

func (l *BetaLauncher) Launch() error {
	time.Sleep(time.Second * 2)
	return nil
}
func (l *BetaLauncher) Name() string                { return "beta launcher" }
func (l *BetaLauncher) TimeDuration() time.Duration { return time.Second * 1 }

func TestMain(t *testing.T) {
	l := New()
	l.AddFunc("proc 1", func() (err error) {
		return nil
	})
	l.Add(&BetaLauncher{})

	l.SetAlertSender(&BetaAlertSender{})
	l.Launch()

	time.Sleep(time.Second)
	runmain()
}
