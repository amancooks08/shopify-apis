package utils

import (
	"errors"
	"regexp"
	"shopify-apis/domain"
)

func ValidateUser(user domain.User) (bool, error) {
	if !validateMobileNumber(user.MobileNumber) {
		return false, errors.New("invalid mobile number")
	}

	if !vaildateName(user.FirstName) {
		return false, errors.New("invalid first name")
	}

	if !vaildateName(user.LastName) {
		return false, errors.New("invalid last name")
	}

	if !validateEmail(user.Email) {
		return false, errors.New("invalid email")
	}
	
	return true, nil

}
func validateMobileNumber(contact string) bool {
	// Define a regular expression pattern for contact numbers("+ country-code mobile-number")
	re := regexp.MustCompile(`^\+\d{1,4}\d{5,15}$`)
	return re.MatchString(contact)
}

func vaildateName(name string) bool {
	re := regexp.MustCompile(`^[A-Za-z-' ]+$`)
	return re.MatchString(name)
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}