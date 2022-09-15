package domain

import (
	"collection/errs"
	"collection/logger"
	"strconv"

	"collection/dto"

	_ "github.com/lib/pq"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryDb struct{
	client *sqlx.DB
}




func (d UserRepositoryDb)AddUser(c Collection,add []dto.Address)(*errs.AppError){
	sqlInsert := `INSERT INTO "Collection"(user_id,name,symbol,description,total_supply,seller_fee,mint_price,game_resource,live_mint_start)values($,$,$,$,$,$,$,$,$)`
	result,err := d.client.Exec(sqlInsert,c.User_id,c.Name,c.Symbol,c.Description,c.Total_supply,c.Seller_fee,c.Mint_price,c.Game_resource,c.Live_mint_start)
	if err!=nil{
		logger.Error("Error While creating new account for collection "+err.Error())
		// return errs.NewUnexpectedError("Unexpected error from database")
		return nil
	}

	id,err:=result.LastInsertId()
	if err!=nil{
		logger.Error("Error While getting last insert id"+err.Error())
		// return errs.NewUnexpectedError("Unexpected error from database")
		return nil
	}

	userId := strconv.FormatInt(id,10)
	addInsert := "INSERT INTO Seller(user_id,address,share)values($,$,$)"
	for _,v := range add{
		_,err := d.client.Exec(addInsert,userId,v.Address,v.Share)
		if err!=nil{
			logger.Error("Error While creating new account"+err.Error())
			// return errs.NewUnexpectedError("Unexpected error from database")
			return nil
		}
	}
	return nil
}

func NewUserRepositoryDb(dbCLient *sqlx.DB)UserRepositoryDb{
	return UserRepositoryDb{dbCLient}
}