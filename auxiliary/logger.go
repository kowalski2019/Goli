package auxiliary

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const pathToLogFile = "./logs/log.txt"

func getCurrentDate() string {
	return time.Now().Format("01-02-2006 15:04:05")
}

func WriteLog(level string, message string, error error) {
	var file, err = os.OpenFile(pathToLogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("[ERROR]: Could not open or create log file! Due to: %s", err)
	}

	// Print log to stdout.
	if error != nil {
		if strings.ToUpper(level) == "ERROR" {
			_, err := file.WriteString(fmt.Sprintf("%s\t[%s]: %s, FATAL! \nGo`s opinion: %s\n", getCurrentDate(), strings.ToUpper(level), message, error))
			if err != nil {
				log.Fatalf("[ERROR]: Could not write to log file!")
			}
			log.Fatalf("%s\t[%s]: %s \nGo`s opinion: %s\n", getCurrentDate(), strings.ToUpper(level), message, error)
		} else {
			_, err := file.WriteString(fmt.Sprintf("%s\t[%s]: %s \nGo`s opinion: %s\n", getCurrentDate(), strings.ToUpper(level), message, error))
			if err != nil {
				log.Fatalf("[ERROR]: Could not write to log file!")
			}
			fmt.Printf("%s\t[%s]: %s \nGo`s opinion: %s\n", getCurrentDate(), strings.ToUpper(level), message, error)
		}
	} else {
		if strings.ToUpper(level) == "ERROR" {
			_, err := file.WriteString(fmt.Sprintf("%s\t[%s]: %s \n", getCurrentDate(), strings.ToUpper(level), message))
			if err != nil {
				log.Fatalf("[ERROR]: Could not write to log file!")
			}
			log.Fatalf("%s\t[%s]: %s \n", getCurrentDate(), strings.ToUpper(level), message)
		} else {
			_, err := file.WriteString(fmt.Sprintf("%s\t[%s]: %s \n", getCurrentDate(), strings.ToUpper(level), message))
			if err != nil {
				log.Fatalf("[ERROR]: Could not write to log file!")
			}
			fmt.Printf("%s\t[%s]: %s \n", getCurrentDate(), strings.ToUpper(level), message)
		}
	}

	defer file.Close()
}
