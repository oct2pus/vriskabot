package logging

import (
	"fmt"
	"log"
	"path/filepath"
	"time"
)

var (
	// error logging
	Log         *log.Logger
	currentTime string
)

// logs errors
func CheckError(err error) bool {
	if err != nil {
		fmt.Println("error: ", err)
		Log.Println("error: ", err)
		return true
	}
	return false
}

func CreateLog() {
	// error logging
	path := filepath.Dir(executable)
	currentTime = time.Now().Format("2006-01-02@15h04m")
	file, err := os.Create(path + ".logs@" + currentTime + ".log")
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.Ldate|log.Ltime|log.Llongfile|log.LUTC)
}
