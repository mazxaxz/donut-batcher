package batch

import "go.mongodb.org/mongo-driver/bson/primitive"

type BatchResult struct {
	ID     primitive.ObjectID
	Status Status
}
