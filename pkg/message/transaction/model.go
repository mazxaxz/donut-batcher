package transaction

const MessageTypeTransaction = "transaction"

type Transaction struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Amount string `json:"amount"`
	// TODO
	Currency string `json:"currency"`
}
