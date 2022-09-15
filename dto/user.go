package dto

// CustodialUser custodial user model with wallet address and one to many relation with FlowId
type User struct {
	WalletAddress string   `json:"-" gorm:"primaryKey;not null"`
	FlowId        []FlowId `gorm:"foreignkey:WalletAddress" json:"-"`
}
