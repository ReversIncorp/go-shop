package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func AppValidate() *validator.Validate {
	instance := validator.New() // Инициализация валидатора
	err := instance.RegisterValidation("password", validatePassword)
	if err != nil {
		logrus.Errorf("Failed to register validator: %v", err)
	}
	return instance // Возврат инициализированного экземпляра
}

// Кастомный валидатор для пароля.
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Минимальная длина 8 символов, максимальная 64 символа.
	if len(password) < 8 || len(password) > 64 {
		return false
	}

	hasLower := regexp.MustCompile("[a-z]").MatchString(password)
	hasUpper := regexp.MustCompile("[A-Z]").MatchString(password)
	hasDigit := regexp.MustCompile(`\\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=$begin:math:display$$end:math:display${};':"\\|,.<>\/?]+`).
		MatchString(password)

	// Проверяем наличие хотя бы одной строчной буквы.
	// Проверяем наличие хотя бы одной цифры.
	// Проверяем наличие хотя бы одного специального символа.
	// Пароль должен содержать заглавную, строчную букву, цифру и спец символ.
	return hasUpper && hasLower && hasDigit && hasSpecial
}
