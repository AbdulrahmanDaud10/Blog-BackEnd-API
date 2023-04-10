package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserName  string    `gorm:"size:255;not null;unique" json:"user_name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Hash function takes in password as a string and GenerateFromPassword returns the bcrypt hash of the password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword function takes in hashedPassword and password and then CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent. Returns nil on success, or an error on failure.
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Function BeforeSave checks that the hashed password exists first before saving it, returns errror if its does not exist
func (user *User) BeforeSave() error {
	hashedPassword, err := Hash(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return nil
}

func (user *User) Prepare() {
	user.ID = 0
	user.UserName = html.EscapeString(strings.TrimSpace(user.UserName))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

// Function Validate converts all the strings to lowercases and then verifies is not empty and prompts messages if the Users details are empty or in an invalid format
func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {

	case "update":
		if user.UserName == "" {
			return errors.New("Username Required")
		}
		if user.Password == "" {
			return errors.New("Password Required")
		}
		if user.Email == "" {
			return errors.New("Email Required")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil

	case "login":
		if user.Password == "" {
			return errors.New("Password Required")
		}
		if user.Email == "" {
			return errors.New("Email Required")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid EmailS")
		}

		return nil

	default:
		if user.UserName == "" {
			return errors.New("Required Nickname")
		}
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}

}

// Function SaveUser saves the user to the database and ensures its not empty
func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&user).Error

	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// Function FindAllUsers queries through the database to retrieve the users but to a llmit of 100 users
func (user *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}

	return &users, err
}

// FUnction FindUserByID queries for a user using a specific ID from the users column
// TODO: Figure out the GORM error handling method incase you
func (user *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return &User{}, err
	}

	//TODO: Find an error handling method to ensure that the ID actually exist
	//if gorm.ErrInvalidValue(err) {
	//	return &User{}, errors.New("User not found")
	//}

	return user, err
}

// Function UpdateUser bring up to date a User info
func (user *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error) {
	err := user.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  user.Password,
			"user_Name": user.UserName,
			"email":     user.Email,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// Function DeleteUser drops a user from the User table and returns the affected row
func (user *User) DeleteUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
