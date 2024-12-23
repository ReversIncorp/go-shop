package errorHandling

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
)

func LogErrorWithTracer(err error) {
	if wrappedErr := tracerr.Wrap(err); wrappedErr != nil {
		logrus.Errorf("%s, stack trace: ", wrappedErr.Error())
		printFilteredSourceColor(wrappedErr, "marketplace")
	} else {
		logrus.Error(err)
	}
}

func FatalErrorWithTracer(message string, err error) {
	redColor := "\033[31m"
	resetColor := "\033[0m"
	if wrappedErr := tracerr.Wrap(err); wrappedErr != nil {
		logrus.Printf("%s[ERROR]%s %s:\n", redColor, resetColor, message)
		printFilteredSourceColor(wrappedErr, "marketplace")
		logrus.Fatalf("%s[ERROR DETAILS]%s: %v", redColor, resetColor, err) // Завершаем выполнение
	} else {
		logrus.Printf("%s[ERROR]%s %s:\n", redColor, resetColor, message)
		logrus.Fatalf("%s[ERROR DETAILS]%s: %v", redColor, resetColor, err) // Завершаем выполнение
	}
}

func printFilteredSourceColor(err error, module string) {
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
