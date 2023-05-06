package manager

type TableManager struct {
}

func NewTableManager() *TableManager {
	return &TableManager{}
}

func (tm *TableManager) Initialize() error {
	return nil
}

func (tm *TableManager) Start() error {
	return nil
}

func (tm *TableManager) Shutdown() error {
	return nil
}
