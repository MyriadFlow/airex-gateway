package dtoapis

type AuthenticateRequest struct {
	FlowId    string `json:"flowId" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}
type AuthenticatePayload struct {
	Token string `json:"token"`
}

type GetFlowIdPayload struct {
	Eula   string `json:"eula,omitempty"`
	FlowId string `json:"flowId"`
}
