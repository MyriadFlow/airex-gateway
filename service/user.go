package service

import (
	"collection/config/envconfig"
	"collection/domain"
	"collection/dto"
	"collection/errs"
	"collection/logger"
	"crypto/rand"
	"encoding/json"

	// "fmt"
	"math/big"
	"strconv"

	// "fmt"
	"os"

	"go.uber.org/zap"
)

type UserService interface {
	NewCollection(dto.CollectionRequest) (*dto.JsonFile, *errs.AppError)
}

type DefaultUserService struct {
	repo domain.UserRepositoryDb
}

func (d DefaultUserService) NewCollection(req dto.CollectionRequest) (*dto.JsonFile, *errs.AppError) {
	err := req.ToValidate()
	if err != nil {
		return nil, err
	}

	a := domain.Collection{
		Collection_id:   req.Collection_id,
		Name:            req.Name,
		Symbol:          req.Symbol,
		Description:     req.Description,
		Total_supply:    req.Total_supply,
		Seller_fee:      req.Seller_fee,
		Mint_price:      req.Mint_price,
		Game_resource:   req.Game_resource,
		Live_mint_start: req.Live_mint_start,
	}

	seller := req.Seller
	d.repo.AddCollection(a, seller)

	//Asset File Making by id
	var firstAddress string
	for _, v := range seller {
		firstAddress = v.Address
		break
	}

	c := &dto.JsonFile{
		Price:                uint32(req.Mint_price),
		Number:               int(req.Total_supply),
		SolTreasuryAccount:   firstAddress,
		Creators:             seller,
		UploadMethod:         "bundlr",
		RetainAuthority:      true,
		IsMutable:            true,
		Symbol:               req.Symbol,
		SellerFeeBasisPoints: req.Seller_fee,
	}

	assetPath := envconfig.EnvVars.COLLECTION_PATH
	filename := req.Collection_id
	filepath := assetPath + "/" + filename
	handle := os.Mkdir(filepath, os.ModePerm)
	if handle != nil {
		logger.Error("config is not cretaing new directory", zap.Error(handle))
		return nil, errs.NewStatusInternalServerError("file is not opening")
	}
	con := filepath + "/" + "config.json"
	asset, error := os.OpenFile(con, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if error != nil {
		logger.Error("Assset File is not opening", zap.Error(error))
		return nil, errs.NewStatusInternalServerError("file is not opening")
	}

	b, _ := json.MarshalIndent(c, " ", " ")
	asset.Write(b)

	configFilePath := filepath + "/" + "Asset"
	error = os.Mkdir(configFilePath, os.ModePerm)
	if error != nil {
		logger.Error("config is not cretaing new directory", zap.Error(error))
		return nil, errs.NewStatusInternalServerError("file is not opening")
	}

	var max *big.Int = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(130), nil)
	// Generate cryptographically strong pseudo-random between [0, max)
	n, _ := rand.Int(rand.Reader, max)
	h := n.String()
	Animate := "ipfs://" + req.Game_resource + "/?uniqueParam=" + h

	for i := 0; i < int(req.Total_supply); i++ {

		e := dto.Config{
			Name:                    req.Name,
			Symbol:                  req.Symbol,
			Description:             req.Description,
			Seller_fee_basis_points: req.Seller_fee,
			Image:                   "ipfs://QmUehLuw9dBC1u8xqrWExsb62wGYdYoqb4vnwfeyEQiwNz",
			Animation_url:           Animate,
			Attribute: []dto.Attribute{
				{
					Trait_type: "Developers",
					Value:      "1337 Gamers",
				},
				{
					Trait_type: "Category",
					Value:      "Game",
				},
				{
					Trait_type: "License",
					Value:      "MIT License",
				},
			},
			Properties: dto.Properties{
				File: dto.File{
					Uri:  strconv.Itoa(i) + ".png",
					Type: "image/png",
				},
				Category: "image",
				Creators: req.Seller,
			},
		}
		jsoPath := configFilePath + "/" + strconv.Itoa(i) + ".json"
		jso, error := os.OpenFile(jsoPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if error != nil {
			logger.Error("json File is not opening", zap.Error(error))
			return nil, errs.NewStatusInternalServerError("file is not opening")
		}
		f, _ := json.MarshalIndent(e, " ", " ")
		jso.Write(f)
		imgPath := configFilePath + "/" + strconv.Itoa(i) + ".png"
		dec := dto.Load("dto/hdpng/pacman.png")
		dto.Save(imgPath, dec)
	}
	return c, err
}

func NewUserService(repository domain.UserRepositoryDb) DefaultUserService {
	return DefaultUserService{repository}
}
