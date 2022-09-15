package app

import (
	"collection/dto/dtoapis"
	"collection/internal/pkg/errorso"
	"collection/logger"
	"collection/service"
	"errors"
	"net/http"

	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/gin-gonic/gin"
	// "github.com/gorilla/mux"
)

func (u FlowIdHandler) Authenticate(c *gin.Context) {
	var req dtoapis.AuthenticateRequest
	err := c.BindJSON(&req)
	if err != nil {
		httpo.NewErrorResponse(http.StatusBadRequest, "failed to validate body").
			Send(c, http.StatusBadRequest)
		return
	}

	pasetoToken, err := u.service.VerifySignAndGetPaseto(req.Signature, req.FlowId)
	if err != nil {
		logger.Errorf("failed to get paseto: %s", err)

		// If signature denied
		if errors.Is(err, service.ErrSignDenied) {
			httpo.NewErrorResponse(httpo.SignatureDenied, "signature denied").
				Send(c, http.StatusUnauthorized)
			return
		}

		if errors.Is(err, errorso.ErrRecordNotFound) {
			httpo.NewErrorResponse(httpo.FlowIdNotFound, "flow id not found").
				Send(c, http.StatusNotFound)
			return
		}

		// If unexpected error
		httpo.NewErrorResponse(500, "failed to verify and get paseto").Send(c, 500)
		return
	} else {
		payload := dtoapis.AuthenticatePayload{
			Token: pasetoToken,
		}
		httpo.NewSuccessResponse(http.StatusOK, "Token generated successfully", payload).
			Send(c, http.StatusOK)
	}
}
