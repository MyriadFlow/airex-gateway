package authenticate

import (
	"collection/domain"
	"collection/dto/dtoapis"
	"collection/internal/pkg/errorso"
	"collection/logger"
	"errors"
	"log"
	"net/http"

	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "github.com/gorilla/mux"
)

type AuthenticateHandler struct {
	service DefaultAuthenticateService
}

// ApplyRoutes applies router to gin Router
func ApplyRoutes(dbClient *gorm.DB, r *gin.RouterGroup) {
	flowIdRepo := domain.NewFlowIdRepositoryDb(dbClient)

	authenticateHandler := AuthenticateHandler{NewAuthenticateService(flowIdRepo)}
	g := r.Group("/authenticate")
	{
		g.POST("", authenticateHandler.Authenticate)
	}
}

func (u AuthenticateHandler) Authenticate(c *gin.Context) {
	//TODO check network
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
		if errors.Is(err, ErrSignDenied) {
			httpo.NewErrorResponse(httpo.SignatureDenied, "signature denied").
				Send(c, http.StatusUnauthorized)
			return
		}

		if errors.Is(err, errorso.ErrRecordNotFound) {
			httpo.NewErrorResponse(httpo.FlowIdNotFound, "flow id not found").
				Send(c, http.StatusNotFound)
			return
		}

		log.Printf("failed to verify and get paseto: %s", err)

		logger.Errorf("failed to verify and get paseto: %s", err)
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
