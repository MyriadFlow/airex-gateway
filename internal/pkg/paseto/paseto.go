// Package paseto provides methods to generate and verify paseto tokens
package paseto

import (
	"collection/config/envconfig"
	"fmt"

	pasetoclaims "collection/internal/pkg/paseto/paseto_claims"

	"github.com/vk-rv/pvx"
	"gorm.io/gorm"
)

type PasetoService struct {
	DB *gorm.DB
}

// Returns paseto token for given wallet address
func (s *PasetoService) GetPasetoForUser(walletAddr string) (string, error) {
	pasetoExpiration := envconfig.EnvVars.PASETO_EXPIRATION
	signedBy := envconfig.EnvVars.SIGNED_BY
	customClaims := pasetoclaims.New(s.DB, walletAddr, pasetoExpiration, signedBy)
	privateKey := envconfig.EnvVars.PASETO_PRIVATE_KEY
	symK := pvx.NewSymmetricKey([]byte(privateKey), pvx.Version4)
	pv4 := pvx.NewPV4Local()
	tokenString, err := pv4.Encrypt(symK, &customClaims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *PasetoService) VerifyPaseto(pasetoToken string) (*pasetoclaims.CustomClaims, error) {
	pv4 := pvx.NewPV4Local()
	k := envconfig.EnvVars.PASETO_PRIVATE_KEY
	symK := pvx.NewSymmetricKey([]byte(k), pvx.Version4)
	var cc pasetoclaims.CustomClaims
	cc.Db = s.DB
	err := pv4.
		Decrypt(pasetoToken, symK).
		ScanClaims(&cc)
	if err != nil {
		err = fmt.Errorf("failed to scan claims: %w", err)
		return nil, err
	}
	return &cc, nil
}
