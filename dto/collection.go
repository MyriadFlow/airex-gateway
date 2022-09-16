package dto

import (
	"collection/errs"
	"time"
)

type CollectionRequest struct {
	Id              string    `json:"user_id"`
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	Description     string    `json:"description"`
	Total_supply    uint32    `json:"total_supply"`
	Seller          []Address `json:"seller"`
	Seller_fee      uint64    `json:"seller_fee"`
	Mint_price      float64   `json:"mint_price"`
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
	Price                 float64     `json:"price"`
	Number                int         `json:"number"`
	Gatekeeper            interface{} `json:"gatekeeper"`
	Creators              []Address   `json:"creators"`
	SolTreasuryAccount    string      `json:"solTreasuryAccount"`
	SplTokenAccount       interface{} `json:"splTokenAccount"`
	SplToken              interface{} `json:"splToken"`
	GoLiveDate            time.Time   `json:"goLiveDate"`
	EndSettings           interface{} `json:"endSettings"`
	WhitelistMintSettings interface{} `json:"whitelistMintSettings"`
	HiddenSettings        interface{} `json:"hiddenSettings"`
	UploadMethod          string      `json:"uploadMethod"`
	RetainAuthority       bool        `json:"retainAuthority"`
	IsMutable             bool        `json:"isMutable"`
	Symbol                string      `json:"symbol"`
	SellerFeeBasisPoints  uint64      `json:"sellerFeeBasisPoints"`
	AwsConfig             interface{} `json:"awsConfig"`
	NftStorageAuthToken   string      `json:"nftStorageAuthToken"`
	ShdwStorageAccount    interface{} `json:"shdwStorageAccount"`
}

type Config struct {
	Name                 string      `json:"name"`
	Symbol               string      `json:"symbol"`
	Description          string      `json:"description"`
	SellerFeeBasisPoints uint64      `json:"seller_fee_basis_points"`
	Image                string      `json:"image"`
	ExternalURL          string      `json:"external_url"`
	Edition              int         `json:"edition"`
	Attributes           []Attribute `json:"attributes"`
	Properties           Properties  `json:"properties"`
	AnimationUrl         string      `json:"animation_url"`
}

type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
} 

type Properties struct {
	Files    []File    `json:"files"`
	Category string    `json:"category"`
	Creators []Address `json:"creators"`
}

type File struct {
	URI  string `json:"uri"`
	Type string `json:"type"`
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
