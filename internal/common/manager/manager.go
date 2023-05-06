package manager

type SystemManager interface {
	Initialize() error
	Start() error
	Shutdown() error
}
