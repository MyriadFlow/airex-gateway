package service

import (
	"collection/config/envconfig"
	"collection/domain"
	"collection/internal/pkg/errorso"
	"collection/internal/pkg/paseto"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/streamingfast/solana-go"
)

type DefaultFlowIdService struct {
	flowIdRepo domain.FlowIdRepositoryDb
	userRepo   domain.UserRepositoryDb
}

var ErrSignDenied = errors.New("signature denied")

// Create and insert flow Id into the database and return it
func (i *DefaultFlowIdService) CreateFlowId(walletAddress string) (string, error) {

	//Check if user exist
	_, err := i.userRepo.Get(walletAddress)
	if err != nil {
		if errors.Is(err, errorso.ErrRecordNotFound) {
			//If doesn't exist then add that
			err = i.userRepo.Add(walletAddress)
			if err != nil {
				return "", fmt.Errorf("failed to add user: %w", err)
			}
		} else {
			return "", fmt.Errorf("failed to check if user exist: %w", err)
		}
	}

	flowIdString := uuid.NewString()
	err = i.flowIdRepo.AddFlowId(walletAddress, flowIdString)
	if err != nil {
		return "", fmt.Errorf("failed to add flowId into database: %w", err)
	}

	return flowIdString, nil
}

// VerifySignAndGetPaseto verifies the signature for given flowID and returns paseto if it is valid
// Also deletes the flow id after approving signature
func (i *DefaultFlowIdService) VerifySignAndGetPaseto(signatureHex string, flowId string) (string, error) {
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

	paseto, err := paseto.GetPasetoForUser(i.flowIdRepo.Client, dataFlowId.WalletAddress)
	if err != nil {
		return "", fmt.Errorf("failed to generate paseto: %w", err)
	}

	err = i.flowIdRepo.DeleteFlowId(flowId)
	if err != nil {
		return "", fmt.Errorf("failed to delete flowid: %w", err)
	}
	return paseto, nil
}

func NewFlowIdService(flowIdRepo domain.FlowIdRepositoryDb, userRepo domain.UserRepositoryDb) DefaultFlowIdService {
	return DefaultFlowIdService{flowIdRepo, userRepo}
}
