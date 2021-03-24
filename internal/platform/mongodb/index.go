package mongodb

import "context"

type Indexer interface {
	Index(ctx context.Context) error
}
