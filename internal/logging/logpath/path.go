// Package logpath содержит функции для работы с текстовыми логами приложения
package logpath

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// Logs - структура, представляющая файлы, в которые будут записываться логи.
type Logs struct {
	IncomeLog *os.File
	ErrorLog  *os.File
}

// DestinationLog создаёт файлы логов в указанной директории.
// incoming.txt хранит логи входящих http-запросов.
// logs.txt хранит внутренние логи сервиса.
func DestinationLog(path string) Logs {
	err := os.Mkdir(path, 0750)
	if err != nil && errors.Is(err, os.ErrExist) {
		fmt.Println("Directory already exists. Continue...")
	}

	incomeLog, err := os.OpenFile(path+"/incoming.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("failed to create dir: %v\n", err)
	}

	errorLog, err := os.OpenFile(path+"/logs.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("failed to create dir: %v\n", err)
	}

	return Logs{ErrorLog: errorLog, IncomeLog: incomeLog}
}

// CloseFiles закрывает дескрипторы файлов.
func (l Logs) CloseFiles() error {
	err := l.IncomeLog.Sync()
	if err != nil {
		if innErr := l.ErrorLog.Sync(); err != nil {
			return fmt.Errorf("can not close http log: %v\n can not close inner error log: %v\n", err, innErr)
		}
		return fmt.Errorf("can not close http log: %v\n", err)
	}

	err = l.ErrorLog.Sync()
	if err != nil {
		return fmt.Errorf("can not close inner error log: %v\n", err)
	}
	return nil
}
