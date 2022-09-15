package domain

import (
	"collection/dto"
	"collection/internal/pkg/errorso"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
	// _ "github.com/go-sql-driver/mysql"
)

type FlowIdRepositoryDb struct {
	client *gorm.DB
}

func (i *FlowIdRepositoryDb) GetFlowId(flowId string) (*dto.FlowId, error) {
	db := i.client
	var userFlowId dto.FlowId
	res := db.Find(&userFlowId, &dto.FlowId{
		FlowId: flowId,
	})

	if err := res.Error; err != nil {
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, errorso.ErrRecordNotFound
	}
	return &userFlowId, nil
}

// Adds flow id into database for given wallet Address
func (i *FlowIdRepositoryDb) AddFlowId(walletAddr string, flowId string) error {
	db := i.client
	err := db.Create(&dto.FlowId{
		WalletAddress: walletAddr,
		FlowId:        flowId,
	}).Error

	return err
}

func (i *FlowIdRepositoryDb) DeleteFlowId(flowId string) error {
	db := i.client
	err := db.Delete(&dto.FlowId{
		FlowId: flowId,
	}).Error

	return err
}

func NewFlowIdRepositoryDb(dbCLient *gorm.DB) FlowIdRepositoryDb {
	return FlowIdRepositoryDb{dbCLient}
}
