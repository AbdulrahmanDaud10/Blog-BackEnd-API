package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "user_name") {
		return errors.New("Username is already taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email is already taken")
	}

	if strings.Contains(err, "title") {
		return errors.New("Title is already taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Password is incorrect")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email is already taken")
	}

	return errors.New("Incorrrect details")
}
