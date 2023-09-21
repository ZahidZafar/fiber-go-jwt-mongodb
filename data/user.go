package data

import (
	"greens-basket/utils"
)

type UserPrincipal struct {
	Subject string //`json:"subject" validate:"required"`
	Roles   string //`json:"roles"`
}

type AppUser struct {

	//Entity
	ID      string  `bson:"_id,omitempty" json:"id"`
	Phone   string  `bson:"ph" json:"phon" validate:"required"`
	Balance float32 `bson:"b" json:"balance"`
	Name    string  `bson:"n" json:"name"`
	Role    string  `bson:"r" json:"role"`

	Favorites []string `bson:"-" json:"favorites"`
}

func (au AppUser) Collection() string {
	return string(utils.Users)
}
