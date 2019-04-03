package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Credentials is to be exported
type Credentials struct {
	SystemID         bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	UserID           string        //
	Name             string        //
	Password         string        //
	PasswordValidate string        //
	ApplicationID    string        //
	CentroID         string        //
	MobilePhone      string        //
	Expiry           string        //
	JWT              string        //
	KeyJWT           string        //
	ClaimSet         []Claim       //
	Status           string        // It is set to Active manually by Daniel 'Active' or Inactive.
	IsAdmin          string        //
	IsAnonymous      string        //
	ResetCode        string        //
}

// Claim is
type Claim struct {
	Type  string
	Value string
}
