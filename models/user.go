package models

import (
	"errors"
)

type User struct{
	Username *string `json:"username"`
	Password *string `json:"password"`

	Customerid *string `json:"customerid"`

	// time created
	// time updated
}

func (incomingUser *User) Validate() (error){
	if incomingUser.Username != nil {
		if *incomingUser.Username == "" { 	// Remark: In the first if-statement, we check the pointer to the username field. In the 2nd if-statement, we check its value.
			return errors.New("username cannot be the empty string")
		}
	}

	if incomingUser.Password != nil {
		if *incomingUser.Password == "" { 	// Remark: In the first if-statement, we check the pointer to the username field. In the 2nd if-statement, we check its value.
			return errors.New("password cannot be the empty string")
		}
	}

	if incomingUser.Customerid != nil {
		if *incomingUser.Customerid == "" { 	// Remark: In the first if-statement, we check the pointer to the username field. In the 2nd if-statement, we check its value.
			return errors.New("customerid cannot be the empty string")
		}
	}

	return nil
}