package domain

import (
	"collection/errs"
	"collection/logger"

	"collection/dto"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
	// _ "github.com/go-sql-driver/mysql"
)

type UserRepositoryDb struct {
	client *gorm.DB
}

func (d UserRepositoryDb) AddUser(c Collection, add []dto.Address) *errs.AppError {
	collectionDb := d.client.Model(&Collection{})
	err := collectionDb.Create(&c).Error
	if  err != nil {
		logger.Error("Error While creating new account for collection " + err.Error())
		// return errs.NewUnexpectedError("Unexpected error from database")
		return nil
	}

	sellerDb := d.client.Model(&Seller{})
	for _, v := range add {
		newSeller := Seller{
			User_id: c.User_id,
			Address: v.Address,
			Share:   v.Share,
		}
		err = sellerDb.Create(&newSeller).Error
		if  err != nil {
			logger.Error("Error While creating new account" + err.Error())
			// return errs.NewUnexpectedError("Unexpected error from database")
			return nil
		}
	}
	return nil
}

func NewUserRepositoryDb(dbCLient *gorm.DB) UserRepositoryDb {
	return UserRepositoryDb{dbCLient}
}
