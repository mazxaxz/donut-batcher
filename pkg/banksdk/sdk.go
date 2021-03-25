package banksdk

import (
	"context"
	"fmt"
)

type Client interface {
	Send(_ context.Context, userID, amount string) error
}

type clientContext struct {
}

func New() Client {
	c := clientContext{}
	return &c
}

func (c *clientContext) Send(_ context.Context, userID, amount string) error {
	fmt.Println("### sending money to the bank...")
	fmt.Println(fmt.Sprintf("UserID: %s, Amount: %s sent!", userID, amount))
	return nil
}
