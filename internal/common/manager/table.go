package manager

import (
	"github.com/spaghettifunk/norman/internal/common/model"
)

type TableManager struct {
	OfflineTables   []*model.OfflineTable
	ReatltimeTables []*model.RealtimeTable
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

func (tm *TableManager) CreateTable(isOffline bool, config []byte) error {
	if isOffline {
		tb, err := model.NewOfflineTable(config)
		if err != nil {
			return err
		}
		tm.OfflineTables = append(tm.OfflineTables, tb)
	} else {
		tb, err := model.NewRealtimeTable(config)
		if err != nil {
			return err
		}
		tm.ReatltimeTables = append(tm.ReatltimeTables, tb)
	}
	return nil
}

func (tm *TableManager) Shutdown() error {
	return nil
}
