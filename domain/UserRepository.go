package domain

import (
	"collection/internal/pkg/errorso"
	"collection/logger"
	"fmt"

	"gorm.io/gorm"
)

type UserRepositoryDb struct {
	Client *gorm.DB
}

// Add adds user with given wallet address to database
func (i *UserRepositoryDb) Add(walletAddress string) error {
	db := i.Client
	newUser := User{
		WalletAddress: walletAddress,
	}
	err := db.Create(&newUser).Error
	return err
}

// Get returns user with given wallet address from database
func (i *UserRepositoryDb) Get(walletAddr string) (*User, error) {
	db := i.Client
	var user User
	res := db.Find(&user, User{
		WalletAddress: walletAddr,
	})

	if err := res.Error; err != nil {
		err = fmt.Errorf("failed to get user from database: %w", err)
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, errorso.ErrRecordNotFound
	}

	return &user, nil
}

// Get returns collection for that user with provided collection id
func (i *UserRepositoryDb) GetCollection(walletAddr string, collectionId string) (*Collection, error) {
	db := i.Client
	var user User
	res := db.Find(&user, User{
		WalletAddress: walletAddr,
	})
	logger.Info(walletAddr)

	if err := res.Error; err != nil {
		err = fmt.Errorf("failed to get user from database: %w", err)
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, errorso.ErrRecordNotFound
	}
	var col Collection
	err := db.Model(&col).Where("creator_wallet_address = ? AND id = ?", walletAddr, collectionId).First(&col).Error
	if err != nil {
		err = fmt.Errorf("failed to get collection from database: %w", err)
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, errorso.ErrRecordNotFound
	}

	return &col, nil
}

func NewUserRepositoryDb(dbCLient *gorm.DB) UserRepositoryDb {
	return UserRepositoryDb{dbCLient}
}
