package flowid

import (
	"collection/domain"
	"collection/internal/pkg/errorso"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrSignDenied = errors.New("signature denied")

type DefaultFlowIdService struct {
	flowIdRepo domain.FlowIdRepositoryDb
	userRepo   domain.UserRepositoryDb
}

type FlowIdService interface {
	CreateFlowId(string) (string, error)
}

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

func NewFlowIdService(flowIdRepo domain.FlowIdRepositoryDb, userRepo domain.UserRepositoryDb) DefaultFlowIdService {
	return DefaultFlowIdService{flowIdRepo, userRepo}
}
