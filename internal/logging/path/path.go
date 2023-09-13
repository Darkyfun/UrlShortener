package path

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Logs struct {
	IncomeLog *os.File
	ErrorLog  io.Writer
}

func DestinationLog(path string) *Logs {

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

	return &Logs{ErrorLog: errorLog, IncomeLog: incomeLog}
}
