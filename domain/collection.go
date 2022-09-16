package domain

import (
	"collection/errs"

	"collection/dto"
)

type Collection struct {
	Id            string `gorm:"primaryKey"`
	Name          string
	Symbol        string
	Description   string
	TotalSupply   uint32
	SellerFee     uint32
	MintPrice     float32
	GameResource  string
	LiveMintStart string
	Sellers       []User `gorm:"many2many:collection_sellers;"`
}

type CollectionSeller struct {
	CollectionId      string `gorm:"primaryKey"`
	UserWalletAddress string
	Share             string
}

type User struct {
	WalletAddress string       `gorm:"primaryKey" json:"walletAddress"`
	FlowIds       []dto.FlowId `gorm:"foreignkey:WalletAddress" json:"-"`
}

type CollectionRepository interface {
	AddCollection(c Collection, add []dto.Address) *errs.AppError
}

// func (a Collection)ToNewCollectionResponseDto()
