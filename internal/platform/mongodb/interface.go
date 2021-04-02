package mongodb

import "go.mongodb.org/mongo-driver/bson"

type SingleResulter interface {
	Decode(v interface{}) error
	DecodeBytes() (bson.Raw, error)
	Err() error
}
