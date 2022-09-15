package dto

import (
	"collection/errs"
	"time"

	"gopkg.in/guregu/null.v4"
)

type CollectionRequest struct {
	Collection_id         string    `json:"user_id"`
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	Description     string    `json:"description"`
	Total_supply    uint32    `json:"total_supply"`
	Seller          []Address `json:"seller"`
	Seller_fee      uint32    `json:"seller_fee"`
	Mint_price      float32   `json:"mint_price"`
	Game_resource   string    `json:"game_resource"`
	Live_mint_start string    `json:"live_mint_start"`
}

type Address struct {
	Address string `json:"address"`
	Share   int    `json:"share"`
}

type CollectionResponse struct {
	Id     string
	Config *JsonFile
}

type JsonFile struct {
	Price                 uint32
	Number                int
	Gatekeeper            null.String
	Creators              []Address
	SolTreasuryAccount    string
	SplTokenAccount       null.String
	SplToken              null.String
	GoLiveDate            *time.Time
	EndSettings           null.String
	WhitelistMintSettings null.String
	HiddenSettings        null.String
	UploadMethod          string
	RetainAuthority       bool
	IsMutable             bool
	Symbol                string
	SellerFeeBasisPoints  uint32
	AwsConfig             null.String
	NftStorageAuthToken   null.String
	ShdwStorageAccount    null.String
}

type Config struct {
	Name                    string
	Symbol                  string
	Description             string
	Seller_fee_basis_points uint32
	Image                   string
	Animation_url           string
	Attribute               []Attribute
	Properties              Properties
}

type Attribute struct {
	Trait_type string
	Value      string
}

type Properties struct {
	File     File
	Category string
	Creators []Address
}

type File struct {
	Uri  string
	Type string
}

func (c CollectionRequest) ToValidate() *errs.AppError {
	if c.Name == "" || c.Description == "" || c.Game_resource == "" || c.Live_mint_start == "" || c.Symbol == "" || c.Total_supply == 0 || c.Seller == nil || c.Seller_fee == 0 || c.Mint_price == 0 {
		return errs.NewValidationError("Fields Empty")
	}
	if c.Mint_price < 0 {
		return errs.NewValidationError("Price should be Greater than Zero")
	}
	var share int
	var index int
	for i, v := range c.Seller {
		share += v.Share
		index += i
	}

	if index > 4 {
		return errs.NewValidationError("address should not be more than 4.")
	}
	if share > 100 || share < 100 {
		return errs.NewValidationError("Share should be 100")
	}
	return nil

}
