package types

type Dispatcher interface {
	DispatchMessageLocal(m *Message)
}
