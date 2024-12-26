package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode/utf8"
)

// ValidateEmail 이메일 주소 유효성 검사
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePassword 비밀번호 유효성 검사
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)

	if !hasNumber || !hasLetter {
		return fmt.Errorf("password must contain both letters and numbers")
	}

	return nil
}

// ValidateUsername 사용자 이름 유효성 검사
func ValidateUsername(username string) error {
	length := utf8.RuneCountInString(username)
	if length < 3 || length > 30 {
		return fmt.Errorf("username must be between 3 and 30 characters")
	}

	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	if err != nil {
		return fmt.Errorf("invalid username format")
	}
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

// ValidatePhoneNumber 전화번호 유효성 검사
func ValidatePhoneNumber(phone string) error {
	matched, err := regexp.MatchString(`^\+?[0-9]{10,15}$`, phone)
	if err != nil {
		return fmt.Errorf("invalid phone number format")
	}
	if !matched {
		return fmt.Errorf("invalid phone number")
	}
	return nil
}
