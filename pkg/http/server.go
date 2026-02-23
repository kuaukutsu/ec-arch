package http

type Server interface {
	Start()
	Notify() <-chan error
	Shutdown() error
}
