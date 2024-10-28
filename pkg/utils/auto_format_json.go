package utils

import (
	"bytes"
	"encoding/json"
)

func AutoFormatJSON(data interface{}) (string, error) {
	switch v := data.(type) {
	case []byte:
		// Если это []byte, проверяем, является ли это JSON, и форматируем
		formattedJSON, err := fromBytesToJSON(v)
		if err != nil {
			return "", err
		}
		return formattedJSON, nil
	default:
		// Если это не []byte, пытаемся сериализовать структуру в JSON
		formattedJSON, err := fromMapToJSON(data)
		if err != nil {
			return "", err
		}
		return formattedJSON, nil
	}
}

func fromBytesToJSON(body []byte) (string, error) {
	var prettyJSON bytes.Buffer
	if json.Valid(body) {
		// Форматируем JSON с отступами
		err := json.Indent(&prettyJSON, body, "", "  ")
		if err != nil {
			return prettyJSON.String(), err
		}
		return prettyJSON.String(), err
	}
	// Если тело не JSON, просто копируем его
	prettyJSON.Write(body)
	return prettyJSON.String(), nil
}

func fromMapToJSON(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
