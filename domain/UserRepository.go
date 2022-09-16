package domain

import (
	"collection/internal/pkg/errorso"
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
	err := db.Model(&newUser).Create(&newUser).Error
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

func NewUserRepositoryDb(dbCLient *gorm.DB) UserRepositoryDb {
	return UserRepositoryDb{dbCLient}
}
