// Package pasetoclaims provides claim declaration for token generation and verification
package pasetoclaims

import (
	"os/user"
	"time"

	"github.com/vk-rv/pvx"
	"gorm.io/gorm"
)

type CustomClaimsWrapper struct {
	Db *gorm.DB
	Cc CustomClaims
}

// CustomClaims defines claims for paseto containing wallet address, signed by and RegisteredClaims
type CustomClaims struct {
	WalletAddress string `json:"walletAddress"`
	SignedBy      string `json:"signedBy"`
	pvx.RegisteredClaims
}

// Valid checks if the claims are valid agaist RegisteredClaims and checks if wallet address
// exist in database
func (c CustomClaimsWrapper) Valid() error {
	db := c.Db
	if err := c.Cc.RegisteredClaims.Valid(); err != nil {
		return err
	}
	err := db.Model(&user.User{}).Where("wallet_address = ?", c.Cc.WalletAddress).First(&user.User{}).Error
	return err
}

// New returns CustomClaims with wallet address, signed by and expiration
func New(db *gorm.DB, walletAddress string, expiration time.Duration, signedBy string) CustomClaimsWrapper {
	expirationTime := time.Now().Add(expiration)
	return CustomClaimsWrapper{
		Db: db,
		Cc: CustomClaims{
			walletAddress,
			signedBy,
			pvx.RegisteredClaims{
				Expiration: &expirationTime,
			},
		},
	}
}
