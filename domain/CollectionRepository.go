package domain

import (
	"collection/errs"
	"collection/logger"
	"fmt"

	"collection/dto"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
	// _ "github.com/go-sql-driver/mysql"
)

type CollectionRepositoryDb struct {
	client *gorm.DB
}

func (d CollectionRepositoryDb) AddCollection(c Collection, add []dto.Address) *errs.AppError {
	collectionDb := d.client
	err := collectionDb.Create(&c).Error
	if err != nil {
		logger.Error("Error While creating new account for collection " + err.Error())
		// return errs.NewUnexpectedError("Unexpected error from database")
		return nil
	}
	association := collectionDb.Model(&c).Association("Sellers")
	if err = association.Error; err != nil {
		logger.Error("Error While association with sellers for collections " + err.Error())
		// return errs.NewUnexpectedError("Unexpected error from database")
		return nil
	}
	for _, v := range add {
		user := User{
			WalletAddress: v.Address,
		}
		_ = CollectionSeller{
			CollectionId:      c.Id,
			UserWalletAddress: v.Address,
			Share:             fmt.Sprint(v.Share),
		}
		err = association.Append(&user)

		if err != nil {
			logger.Error("Error While creating new account" + err.Error())
			// return errs.NewUnexpectedError("Unexpected error from database")
			return nil
		}
	}

	return nil
}

func NewCollectionRepositoryDb(dbCLient *gorm.DB) CollectionRepositoryDb {
	return CollectionRepositoryDb{dbCLient}
}
