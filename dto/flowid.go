package dto

type FlowId struct {
	WalletAddress string
	FlowId        string `gorm:"primary_key"`
}
