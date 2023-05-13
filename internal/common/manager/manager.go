package manager

type SystemManager interface {
	Initialize() error
	Execute(config []byte) error
	Shutdown() error
}
