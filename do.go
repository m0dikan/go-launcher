package launcher

func Do(name string, a AlertSender, f func() error) {
	p := &LaunchProcessor{}
	p.SetAlertSender(a)
	p.AddFunc(name, f)
	p.Launch()
}
