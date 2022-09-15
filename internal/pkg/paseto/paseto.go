// Package paseto provides methods to generate and verify paseto tokens
package paseto

import (
	"collection/config/envconfig"
	"fmt"

	pasetoclaims "collection/internal/pkg/paseto/paseto_claims"

	"github.com/vk-rv/pvx"
	"gorm.io/gorm"
)

// Returns paseto token for given wallet address
func GetPasetoForUser(db *gorm.DB, walletAddr string) (string, error) {
	pasetoExpiration := envconfig.EnvVars.PASETO_EXPIRATION
	signedBy := envconfig.EnvVars.SIGNED_BY
	customClaims := pasetoclaims.New(db, walletAddr, pasetoExpiration, signedBy)
	privateKey := envconfig.EnvVars.PASETO_PRIVATE_KEY
	symK := pvx.NewSymmetricKey([]byte(privateKey), pvx.Version4)
	pv4 := pvx.NewPV4Local()
	tokenString, err := pv4.Encrypt(symK, customClaims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyPaseto(pasetoToken string) error {
	pv4 := pvx.NewPV4Local()
	k := envconfig.EnvVars.PASETO_PRIVATE_KEY
	symK := pvx.NewSymmetricKey([]byte(k), pvx.Version4)
	var cc pasetoclaims.CustomClaims
	err := pv4.
		Decrypt(pasetoToken, symK).
		ScanClaims(&cc)
	if err != nil {
		err = fmt.Errorf("failed to scan claims: %w", err)
		return err
	}
	return nil
}
