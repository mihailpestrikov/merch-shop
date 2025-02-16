package dto

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []ReceivedCoin `json:"received"`
	Sent     []SentCoin     `json:"sent"`
}

type ReceivedCoin struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentCoin struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type MerchItemResponse struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type MerchItemsResponse struct {
	MerchItems []MerchItemResponse `json:"merchItems"`
}
