// Package pasetoclaims provides claim declaration for token generation and verification
package pasetoclaims

import (
	"os/user"
	"time"

	"github.com/vk-rv/pvx"
	"gorm.io/gorm"
)

// CustomClaims defines claims for paseto containing wallet address, signed by and RegisteredClaims
type CustomClaims struct {
	WalletAddress string `json:"walletAddress"`
	SignedBy      string `json:"signedBy"`
	pvx.RegisteredClaims
	Db *gorm.DB `json:"-"`
}

// Valid checks if the claims are valid agaist RegisteredClaims and checks if wallet address
// exist in database
func (c CustomClaims) Valid() error {
	db := c.Db
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}
	err := db.Where("wallet_address = ?", c.WalletAddress).First(&user.User{}).Error
	return err
}

// New returns CustomClaims with wallet address, signed by and expiration
func New(db *gorm.DB, walletAddress string, expiration time.Duration, signedBy string) CustomClaims {
	expirationTime := time.Now().Add(expiration)
	return CustomClaims{
		walletAddress,
		signedBy,
		pvx.RegisteredClaims{
			Expiration: &expirationTime,
		},
		db,
	}
}
