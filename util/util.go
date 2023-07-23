package util

import "go.mongodb.org/mongo-driver/bson/primitive"

func ObjectIdParser(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID([]byte{}), err
	}
	return oid, nil
}
