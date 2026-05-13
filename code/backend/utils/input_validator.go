package utils

import (
	"fmt"
	"regexp"
)

var (
	// These regexes match characters that are NOT allowed
	// The ^ character negates the set of characters inside []. So it matches any character that is NOT in the set.
	// Examples:
	// [^a-zA-Z0-9_] matches any character that is not a letter, number, or underscore
	// [^a-zA-Z0-9!@#$%^&*()_+=\-\. ] matches any character that is not a letter, number, special character, or space.
	// The |^ | $ part matches spaces at the beginning or end of the string.
	// [^<>] matches any character that is not < or >
	usernameForbiddenRegex = regexp.MustCompile(`[^a-zA-Z0-9_]`)
	passwordForbiddenRegex = regexp.MustCompile(`(^ )|( $)|[^a-zA-Z0-9!@#$%^&*()_+=\-\. ]`)
	genericForbiddenRegex  = regexp.MustCompile(`[<>"'%;]`)
)

// ValidateUsername returns an error if the username contains invalid characters
func ValidateUsername(username string) error {
	if usernameForbiddenRegex.MatchString(username) {
		return fmt.Errorf("username contains invalid characters (only alphanumeric and underscores allowed)")
	}
	return nil
}

// ValidatePassword returns an error if the password contains invalid characters
func ValidatePassword(password string) error {
	if passwordForbiddenRegex.MatchString(password) {
		return fmt.Errorf(`password contains invalid characters (allowed: a-z, A-Z, 0-9, spaces, and !@#$%%^&*()_+=-.). Leading or trailing spaces are not allowed.`)
	}
	return nil
}

// ValidateGeneric returns an error if the input contains dangerous characters
func ValidateGeneric(input string, fieldName string) error {
	if genericForbiddenRegex.MatchString(input) {
		return fmt.Errorf("%s contains invalid characters (< > \" ' %% ; are not allowed)", fieldName)
	}
	return nil
}
