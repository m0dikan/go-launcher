package launcher

type AlertSender interface {
	Send(msg string)
}
