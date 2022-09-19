package collection

import (
	"bytes"
	"collection/config/envconfig"
	"collection/domain"
	"collection/dto"
	"collection/errs"
	"collection/logger"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"

	// "fmt"
	"math/big"
	"strconv"

	// "fmt"
	"os"

	"go.uber.org/zap"
)

type CollectionService interface {
	NewCollection(dto.CollectionRequest) (*dto.JsonFile, *errs.AppError)
}

type DefaultCollectionService struct {
	repo domain.CollectionRepositoryDb
}

func (d DefaultCollectionService) NewCollection(req dto.CollectionRequest) (*dto.JsonFile, *errs.AppError) {
	appErr := req.ToValidate()
	if appErr != nil {
		return nil, appErr
	}

	a := domain.Collection{
		Id:            req.Id,
		Name:          req.Name,
		Symbol:        req.Symbol,
		Description:   req.Description,
		TotalSupply:   req.Total_supply,
		SellerFee:     req.Seller_fee,
		MintPrice:     req.Mint_price,
		GameResource:  req.Game_resource,
		LiveMintStart: req.Live_mint_start,
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
		Price:                req.Mint_price,
		Number:               int(req.Total_supply),
		SolTreasuryAccount:   firstAddress,
		Creators:             seller,
		UploadMethod:         "nft_storage",
		RetainAuthority:      true,
		IsMutable:            true,
		Symbol:               req.Symbol,
		SellerFeeBasisPoints: req.Seller_fee,
		NftStorageAuthToken:  envconfig.EnvVars.NFT_STORAGE,
	}

	assetPath := envconfig.EnvVars.COLLECTION_PATH
	filename := req.Id
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

	assetsFilePath := filepath + "/" + "assets"
	error = os.Mkdir(assetsFilePath, os.ModePerm)
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
			Name:                 req.Name,
			Symbol:               req.Symbol,
			Description:          req.Description,
			SellerFeeBasisPoints: req.Seller_fee,
			Image:                "ipfs://QmUehLuw9dBC1u8xqrWExsb62wGYdYoqb4vnwfeyEQiwNz",
			AnimationUrl:         Animate,
			Attributes: []dto.Attribute{
				{
					TraitType: "Developers",
					Value:     "1337 Gamers",
				},
				{
					TraitType: "Category",
					Value:     "Game",
				},
				{
					TraitType: "License",
					Value:     "MIT License",
				},
			},
			Properties: dto.Properties{
				Files: []dto.File{{
					URI:  strconv.Itoa(i) + ".png",
					Type: "image/png",
				}},
				Category: "image",
				Creators: req.Seller,
			},
		}
		jsoPath := assetsFilePath + "/" + strconv.Itoa(i) + ".json"
		jso, error := os.OpenFile(jsoPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if error != nil {
			logger.Error("json File is not opening", zap.Error(error))
			return nil, errs.NewStatusInternalServerError("file is not opening")
		}
		f, _ := json.MarshalIndent(e, " ", " ")
		jso.Write(f)
		imgPath := assetsFilePath + "/" + strconv.Itoa(i) + ".png"
		dec := dto.Load("dto/hdpng/pacman.png")
		dto.Save(imgPath, dec)

	}
	configFile := path.Join(filepath, "config.json")
	cacheFile := path.Join(filepath, "cache.json")
	//TODO handle error
	go SugarUploadAndDeploy(configFile, cacheFile, assetsFilePath)

	return c, appErr
}

func SugarUploadAndDeploy(configFile string, cacheFile string, assetsFilePath string) {
	err := SugarUpload(configFile, cacheFile, assetsFilePath)
	if err != nil {
		logger.Errorf("failed to upload using sugar :%s", err)
		return
	}
	SugarDeploy(configFile, cacheFile)
}

func SugarUpload(configFile string, cacheFile string, assetsFilePath string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("sugar", "upload", "--config", configFile, "--cache", cacheFile, assetsFilePath)
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("failed to upload artifacts: %s %s %w", stderr.String(), string(stdout), err)
	}
	return nil
}

func SugarDeploy(configFile string, cacheFile string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("sugar", "deploy", "--config", configFile, "--cache", cacheFile)
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("failed to deploy collection: %s %s %w", stderr.String(), string(stdout), err)
	}
	return nil
}
func NewCollectionService(repository domain.CollectionRepositoryDb) DefaultCollectionService {
	return DefaultCollectionService{repository}
}
