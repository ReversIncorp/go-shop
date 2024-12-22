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

	// Проверяем наличие хотя бы одной заглавной буквы
	hasUppercase, hasLowercase, hasDigit, hasSpecialChar := regexp.MustCompile(`[A-Z]`).MatchString(password),
		regexp.MustCompile(`[a-z]`).MatchString(password),
		regexp.MustCompile(`[0-9]`).MatchString(password),
		///TODO: переделать тут валидацию
		regexp.MustCompile("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[!@#$%^&*()_\\+\\-=\n$begin:math:display"+
			"$$end:math:display${};':\"\\\\|,.<>\\/?]).{8,}$").MatchString(password)
	// Проверяем наличие хотя бы одной строчной буквы.
	// Проверяем наличие хотя бы одной цифры.
	// Проверяем наличие хотя бы одного специального символа.
	// Пароль должен содержать заглавную, строчную букву, цифру и спец символ.
	return hasUppercase && hasLowercase && hasDigit && hasSpecialChar
}
