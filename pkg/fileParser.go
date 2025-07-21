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
	FilePath      string
	Tasks         []string
	weekDay       string
	month         string
	weekDayNumber int
	FileContent   []string
}

func Init(filePath string) *TodayTasks {
	weekday := time.Now().Weekday().String()
	month := time.Now().Month().String()
	dataNumber := time.Now().Day()

	return &TodayTasks{
		Tasks:         make([]string, 0),
		weekDay:       weekday,
		month:         month,
		weekDayNumber: dataNumber,
		FilePath:      filePath,
	}

}

func (task *TodayTasks) ParseFile() {
	flag := false
	file, err := os.Open(task.FilePath)
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
			}
			if regex.Match([]byte(line)) {
				task.Tasks = append(task.Tasks, line)
			}
		}

		task.FileContent = append(task.FileContent, line)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error durning the file reading. Error ", err)
	}

}

func (task *TodayTasks) PrintFile() {
	for _, line := range task.FileContent {
		println(line)
	}
}

func (task *TodayTasks) PrintTodaysTasks() {
	file, err := os.OpenFile("myfile.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("Could not open the file. Error ", err)
	}
	for _, todo := range task.Tasks {
		file.WriteString(todo + "\n")
	}
	defer file.Close()

}

func (task *TodayTasks) ReplaceTodos() {
	flag := false
	today := fmt.Sprintf("%s – %s %d", task.weekDay, task.month, task.weekDayNumber)
	regex := regexp.MustCompile(`- \[.*\]`)

	taskIndex := 0 // index into updated task.Tasks

	for i, line := range task.FileContent {
		if strings.Contains(line, today) {
			flag = true
			continue
		}
		if flag {
			if strings.Contains(line, "---") {
				break // End of today’s section
			}
			if regex.Match([]byte(line)) && taskIndex < len(task.Tasks) {
				task.FileContent[i] = task.Tasks[taskIndex]
				taskIndex++
			}
		}
	}
}
