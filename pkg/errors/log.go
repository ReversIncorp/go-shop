package errors

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
)

func LogErrorWithTracer(err error) {
	if wrappedErr := tracerr.Wrap(err); wrappedErr != nil {
		logrus.Errorf("%s, stack trace: ", wrappedErr.Error())
		PrintFilteredSourceColor(wrappedErr, "marketplace")
	} else {
		logrus.Error(err)
	}
}

func PrintFilteredSourceColor(err error, module string) {
	// Получаем стектрейс из ошибки
	tracedErr := tracerr.Wrap(err)
	stack := tracerr.StackTrace(tracedErr)

	// Фильтруем стектрейс
	var filteredStack []string
	for _, frame := range stack {
		frameStr := frame.String() // Конвертация Frame в строку
		if strings.Contains(frameStr, module) {
			filteredStack = append(filteredStack, frameStr)
		}
	}

	for _, frame := range filteredStack {
		fmt.Println(frame)
	}
}
