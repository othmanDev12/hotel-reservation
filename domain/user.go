package domain

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	bcryptCost       = 12
	minLengFirstName = 7
	minLengLastName  = 7
	minLengPassword  = 9
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p *CreateUserParams) ValidateUser(c *CreateUserParams) []string {
	var errors []string
	if len(c.FirstName) < minLengFirstName {
		errors = append(errors, fmt.Sprintf("length of first name must be greater than %d", len(c.FirstName)))
	}
	if len(c.LastName) < minLengLastName {
		errors = append(errors, fmt.Sprintf("length of last name must be grather than %d", len(c.LastName)))
	}
	if len(c.Password) < minLengPassword {
		errors = append(errors, fmt.Sprintf("length od password must be greater than %d", len(c.Password)))
	}
	if validEmail(c.Email) {
		errors = append(errors, fmt.Sprintf("email is invalid"))
	}
	return errors
}

func validEmail(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

type User struct {
	// json is like @JsonProperty in java , and for json  ignore we can use omitempty
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"password" json:"-"`
}

func NewCreateUser(userParams CreateUserParams) (*User, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(userParams.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         userParams.FirstName,
		LastName:          userParams.LastName,
		Email:             userParams.Email,
		EncryptedPassword: string(encryptedPass),
	}, nil
}
