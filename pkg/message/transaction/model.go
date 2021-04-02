package transaction

const MessageTypeTransaction = "transaction"

type Transaction struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Amount string `json:"amount"`
	// Currency represented in ISO 4217 standard
	Currency string `json:"currency"`
}
