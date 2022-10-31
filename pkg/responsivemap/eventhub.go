package rmap

type IEvent interface {
	Name() string
	Msg() interface{}
}

type IListener interface {
	Name() string
	Notify(IEvent)
}

func NewEventHub() *EventHub {
	return &EventHub{
		Listeners: NewSubTreeRoot(),
	}
}

type EventHub struct {
	Name      string
	Listeners *SubTreeRoot
}

func (h *EventHub) Off(l IListener) {
	h.Listeners.Remove(l.Name(), l)
}

// On 注册一个监听
func (h *EventHub) On(l IListener) {
	h.Listeners.Add(l.Name(), l)
}

func (h *EventHub) Process(msg IEvent) {
	h.Listeners.Get(msg.Name())

	for _, l := range h.Listeners.Get(msg.Name()) {
		l.(IListener).Notify(msg)
	}
}
