package domain

import (
	"collection/errs"

	"collection/dto"
)

type Collection struct {
	User_id         string  `db:"user_id"`
	Name            string  `db:"name"`
	Symbol          string  `db:"symbol"`
	Description     string  `db:"description"`
	Total_supply    uint32  `db:"total_supply"`
	Seller_fee      uint32  `db:"seller_fee"`
	Mint_price      float32 `db:"mint_price"`
	Game_resource   string  `db:"game_resource"`
	Live_mint_start string  `db:"live_mint_start"`
}

type Seller struct {
	User_id string `db:"user_id"`
	Address string `db:"address"`
	Share   int    `db:"share"`
}

type UserRepository interface{
	AddUser(c Collection,add []dto.Address)(*errs.AppError)
}
// func (a Collection)ToNewCollectionResponseDto()