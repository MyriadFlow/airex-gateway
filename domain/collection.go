package domain

import (
	"collection/errs"

	"collection/dto"
)

type Collection struct {
	Id                   string `gorm:"primaryKey"`
	Name                 string
	Symbol               string
	Description          string
	TotalSupply          uint32
	SellerFee            uint64
	MintPrice            float64
	GameResource         string
	LiveMintStart        string
	Sellers              []User `gorm:"many2many:collection_sellers;"`
	CreatorWalletAddress string `gorm:"type:not null"`
}

type CollectionSeller struct {
	CollectionId      string `gorm:"primaryKey"`
	UserWalletAddress string
	Share             string
}

type User struct {
	WalletAddress string       `gorm:"primaryKey" json:"walletAddress"`
	FlowIds       []dto.FlowId `gorm:"foreignkey:WalletAddress" json:"-"`
	Collections   []Collection `gorm:"foreignkey:CreatorWalletAddress" json:"-"`
}

type CollectionRepository interface {
	AddCollection(c Collection, add []dto.Address) *errs.AppError
}

// func (a Collection)ToNewCollectionResponseDto()
