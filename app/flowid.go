package app

import (
	"collection/config/envconfig"
	"collection/dto/dtoapis"
	"collection/logger"
	"collection/service"
	"net/http"

	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/gin-gonic/gin"
	"github.com/streamingfast/solana-go"
	// "github.com/gorilla/mux"
)

type FlowIdHandler struct {
	service service.DefaultFlowIdService
}

func (u FlowIdHandler) GetFlowId(c *gin.Context) {
	walletAddress := c.Query("walletAddress")

	if walletAddress == "" {
		httpo.NewErrorResponse(http.StatusBadRequest, "wallet address (walletAddress) is required").
			Send(c, http.StatusBadRequest)
		return
	}
	_, err := solana.PublicKeyFromBase58(walletAddress)
	if err != nil {
		logger.Errorf("failed to get pubkey from wallet address (base58) %s: %s", walletAddress, err)
		httpo.NewErrorResponse(httpo.WalletAddressInvalid, "failed to parse wallet address (walletAddress)").Send(c, http.StatusBadRequest)
		return
	}

	flowId, err := u.service.CreateFlowId(walletAddress)
	if err != nil {
		logger.Errorf("failed to generate flow id: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "Unexpected error occured").Send(c, http.StatusInternalServerError)
		return
	}
	userAuthEULA := envconfig.EnvVars.AUTH_EULA
	payload := dtoapis.GetFlowIdPayload{
		FlowId: flowId,
		Eula:   userAuthEULA,
	}
	httpo.NewSuccessResponse(http.StatusOK, "Flowid successfully generated", payload).Send(c, http.StatusOK)
}
