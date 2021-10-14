package launcher

type AlertSender interface {
	Send(...interface{})
}
