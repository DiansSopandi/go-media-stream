package pkg

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/DiansSopandi/media_stream/dto"
	"github.com/go-playground/validator/v10"
)

func ValidateCreateUserRequest(req *dto.UserCreateRequest) error {
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !isValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	// if len(req.Password) < 8 {
	// 	return fmt.Errorf("password must be at least 8 characters")
	// }

	// Enhanced password validation
	if err := validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

func ValidateRegisterUserRequest(req *dto.UserRegisterRequest) error {

	validate := validator.New()

	err := validate.Struct(req)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("%s %s", e.Field(), validationMessage(e))
		}
	}

	if req.PasswordConfirm == "" {
		return fmt.Errorf("password confirmation is required")
	}

	if req.Password != req.PasswordConfirm {
		return fmt.Errorf("password and password confirmation do not match")
	}

	if len(req.Roles) == 0 {
		return fmt.Errorf("at least one role is required")
	}

	return nil
}

func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters long", e.Param())
	case "eqfield":
		return fmt.Sprintf("must be equal to %s", e.Param())
	default:
		return fmt.Sprintf("is invalid (%s)", e.Tag())
	}
}

func isValidEmail(email string) bool {
	// Simple email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidationError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "validation failed") ||
		strings.Contains(errMsg, "not found") ||
		strings.Contains(errMsg, "required") ||
		strings.Contains(errMsg, "invalid")
}

func isConflictError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "already exists") ||
		strings.Contains(errMsg, "duplicate")
}

func validatePassword(password string) error {
	// Check minimum length
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Check maximum length (optional, prevents extremely long passwords)
	if len(password) > 128 {
		return fmt.Errorf("password must not exceed 128 characters")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	// Check each character for required types
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Validate all requirements are met
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	// Check for alphanumeric combination (at least one letter and one number)
	hasLetter := hasUpper || hasLower
	if !hasLetter || !hasNumber {
		return fmt.Errorf("password must contain a combination of letters and numbers")
	}

	// Optional: Check for common weak patterns
	if err := checkWeakPatterns(password); err != nil {
		return err
	}

	return nil
}

func checkWeakPatterns(password string) error {
	// Check for sequential characters (123, abc, etc.)
	// sequential := regexp.MustCompile(`(?i)(123|234|345|456|567|678|789|890|abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)`)
	// if sequential.MatchString(password) {
	// 	return fmt.Errorf("password should not contain sequential characters")
	// }
	// Check for sequential numbers (longer sequences)
	sequentialNumbers := regexp.MustCompile(`(012|123|234|345|456|567|678|789|890|901)`)
	if sequentialNumbers.MatchString(password) {
		return fmt.Errorf("password should not contain sequential numbers")
	}

	// Check for sequential letters (case insensitive)
	sequentialLetters := regexp.MustCompile(`(?i)(abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)`)
	if sequentialLetters.MatchString(password) {
		return fmt.Errorf("password should not contain sequential letters")
	}

	// Check for keyboard patterns
	keyboardPatterns := regexp.MustCompile(`(?i)(qwerty|asdf|zxcv|qwer|asdf|zxcv|1234|4321)`)
	if keyboardPatterns.MatchString(password) {
		return fmt.Errorf("password should not contain keyboard patterns")
	}

	// Check for repeated characters (3 or more)
	if hasRepeatedChars(password, 3) {
		return fmt.Errorf("password should not contain more than 2 consecutive identical characters")
	}

	// Check for common weak passwords
	commonPasswords := []string{
		"password", "12345678", "qwerty", "abc123", "letmein",
		"password123", "admin123", "welcome123", "monkey",
		"dragon", "sunshine", "iloveyou", "trustno1",
	}

	passwordLower := strings.ToLower(password)
	for _, common := range commonPasswords {
		if strings.Contains(passwordLower, strings.ToLower(common)) {
			return fmt.Errorf("password contains common weak pattern: %s", common)
		}
	}

	// Check for all numbers
	allNumbers := regexp.MustCompile(`^\d+$`)
	if allNumbers.MatchString(password) {
		return fmt.Errorf("password should not contain only numbers")
	}

	// Check for all letters
	allLetters := regexp.MustCompile(`^[a-zA-Z]+$`)
	if allLetters.MatchString(password) {
		return fmt.Errorf("password should contain numbers or special characters")
	}

	// Check for repeated characters (aaa, 111, etc.)
	// repeated := regexp.MustCompile(`(.)\1{2,}`)
	// if repeated.MatchString(password) {
	// 	return fmt.Errorf("password should not contain more than 2 consecutive identical characters")
	// }

	for _, common := range commonPasswords {
		if matched, _ := regexp.MatchString(`(?i)`+regexp.QuoteMeta(common), password); matched {
			return fmt.Errorf("password is too common, please choose a stronger password")
		}
	}

	return nil
}

func validatePasswordWithRegex(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Combined regex for all requirements
	// (?=.*[a-z]) - at least one lowercase
	// (?=.*[A-Z]) - at least one uppercase
	// (?=.*\d) - at least one digit
	// (?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]) - at least one special char
	// .{8,} - minimum 8 characters
	passwordRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]).{8,}$`)

	if !passwordRegex.MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}

// Custom function to check repeated characters
func hasRepeatedChars(s string, maxRepeat int) bool {
	if len(s) < maxRepeat {
		return false
	}

	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
			if count >= maxRepeat {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

// Alternative: More comprehensive repeated character check
func hasRepeatedCharsAdvanced(password string, maxRepeat int) bool {
	runes := []rune(password)
	if len(runes) < maxRepeat {
		return false
	}

	count := 1
	for i := 1; i < len(runes); i++ {
		if runes[i] == runes[i-1] {
			count++
			if count >= maxRepeat {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}
