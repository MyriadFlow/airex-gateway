package authenticate

import (
	"collection/config/envconfig"
	"collection/domain"
	"collection/internal/pkg/paseto"
	"errors"
	"fmt"

	"github.com/streamingfast/solana-go"
)

type DefaultAuthenticateService struct {
	flowIdRepo domain.FlowIdRepositoryDb
}

type AuthenticateService interface {
	NewAuthenticateService(domain.FlowIdRepositoryDb) DefaultAuthenticateService
}

var ErrSignDenied = errors.New("signature denied")

// VerifySignAndGetPaseto verifies the signature for given flowID and returns paseto if it is valid
// Also deletes the flow id after approving signature
func (i *DefaultAuthenticateService) VerifySignAndGetPaseto(signatureHex string, flowId string) (string, error) {
	dataFlowId, err := i.flowIdRepo.GetFlowId(flowId)
	if err != nil {
		return "", fmt.Errorf("failed to get flow id from database: %w", err)
	}

	// Prepare expected signing data (msg)
	authEula := envconfig.EnvVars.AUTH_EULA
	signingData := fmt.Sprintf("%s%s", authEula, dataFlowId.FlowId)

	solanaSignature, err := solana.NewSignatureFromString(signatureHex)
	if err != nil {
		return "", fmt.Errorf("failed to get signature from hex signature: %w", err)
	}

	pubKey, err := solana.PublicKeyFromBase58(dataFlowId.WalletAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get pubkey from wallet address referred by flowid: %w", err)
	}
	signatureApproved := solanaSignature.Verify(pubKey, []byte(signingData))

	//If signature not approved then return error
	if !signatureApproved {
		return "", ErrSignDenied
	}

	pasetoService := paseto.PasetoService{
		DB: i.flowIdRepo.Client,
	}
	paseto, err := pasetoService.GetPasetoForUser(dataFlowId.WalletAddress)
	if err != nil {
		return "", fmt.Errorf("failed to generate paseto: %w", err)
	}

	err = i.flowIdRepo.DeleteFlowId(flowId)
	if err != nil {
		return "", fmt.Errorf("failed to delete flowid: %w", err)
	}
	return paseto, nil
}

func NewAuthenticateService(flowIdRepo domain.FlowIdRepositoryDb) DefaultAuthenticateService {
	return DefaultAuthenticateService{flowIdRepo}
}
