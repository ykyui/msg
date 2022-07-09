package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

func (u UserInfo) SignupVerify() bool {
	if len(u.Username) < 5 {
		return false
	}
	if len(u.Password) < 5 {
		return false
	}
	return true
}

func (u *UserInfo) Select() error {
	collection := getCollection(u)
	return collection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"_id": u.ID}, {"username": u.Username}}}).Decode(&u)
}
func (u *UserInfo) Insert() (string, error) {
	collection := getCollection(u)
	id, err := collection.InsertOne(context.TODO(), u)
	return id.InsertedID.(primitive.ObjectID).String(), err

}
func (u *UserInfo) Update() error {
	collection := getCollection(u)
	_, err := collection.UpdateByID(context.Background(), u.ID, bson.M{"$set": u})
	return err
}
func (u *UserInfo) Delete() error {
	return nil
}
