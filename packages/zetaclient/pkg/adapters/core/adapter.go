package core

//go:generate mockery --name Adapter

type Adapter interface {
	MonitorCore()
	GetPromCounter(name string) (prom.Counter, error)
}
