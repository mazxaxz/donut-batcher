package banksdk

import (
	"context"
	"fmt"
)

type Clienter interface {
	Send(_ context.Context, userID, amount, currency string) error
}

type clientContext struct{}

func New() Clienter {
	c := clientContext{}
	return &c
}

func (c *clientContext) Send(_ context.Context, userID, amount, currency string) error {
	fmt.Println("### sending money to the bank...")
	fmt.Println(fmt.Sprintf("UserID: %s, Amount: %s, Currency: %s sent!", userID, amount, currency))
	return nil
}
