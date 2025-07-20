package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type TodayTasks struct {
	Tasks         []string
	weekDay       string
	month         string
	weekDayNumber int
}

func Init() *TodayTasks {
	weekday := time.Now().Weekday().String()
	month := time.Now().Month().String()
	dataNumber := time.Now().Day()
	return &TodayTasks{
		Tasks:         make([]string, 0),
		weekDay:       weekday,
		month:         month,
		weekDayNumber: dataNumber,
	}

}

func (task *TodayTasks) ParseFile(filepath string) {
	flag := false
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error opening the file ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	today := fmt.Sprintf("%s – %s %d", task.weekDay, task.month, task.weekDayNumber)
	regex := regexp.MustCompile(`- \[.*\]`)

	println(today)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, today) {
			flag = true
			fmt.Println("Todays exisit")
			continue
		}
		if flag {
			if strings.Contains(line, "---") {
				flag = false
				break
			}
			if regex.Match([]byte(line)) {
				task.Tasks = append(task.Tasks, line)
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error durning the file reading. Error ", err)
	}

}
