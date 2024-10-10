package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"regexp"
)

var instance *validator.Validate

func AppValidate() *validator.Validate {
	if instance == nil {
		instance = validator.New() // Инициализация валидатора
		err := instance.RegisterValidation("password", validatePassword)
		if err != nil {
			logrus.Errorf("Failed to register validator: %v", err)
		}
	}
	return instance // Возврат инициализированного экземпляра
}

// Кастомный валидатор для пароля
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Минимальная длина 8 символов, максимальная 64 символа
	if len(password) < 8 || len(password) > 64 {
		return false
	}

	// Проверяем наличие хотя бы одной заглавной буквы
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Проверяем наличие хотя бы одной строчной буквы
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Проверяем наличие хотя бы одной цифры
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// Проверяем наличие хотя бы одного специального символа
	hasSpecialChar := regexp.MustCompile(`[!@#\$%\^&\*$begin:math:text$$end:math:text$_\+\-=$begin:math:display$$end:math:display$\{\};:'"\|\\,.<>/?]+`).MatchString(password)

	// Пароль должен содержать заглавную, строчную букву, цифру и спец. символ
	return hasUppercase && hasLowercase && hasDigit && hasSpecialChar
}
