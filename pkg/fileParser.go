package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const DefaultFolder = "~/notes"

type TodayTasks struct {
	FolderPath string
	Tasks      []string
	filename   string
}

func checkFolerExisting(folderPath string) (bool, error) {
	_, err := os.Stat(folderPath)
	if err != nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func Init(folderPath string) *TodayTasks {
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Month().String()
	dayNumber := strconv.Itoa(time.Now().Day())

	//  note_2025-Dec-31.md
	filename := "note_" + year + "_" + month + "_" + dayNumber + ".md"

	_, err := checkFolerExisting(folderPath)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	return &TodayTasks{
		Tasks:      make([]string, 0),
		FolderPath: folderPath,
		filename:   filename,
	}

}

func (task *TodayTasks) ParseFile() {
	fullPath := task.FolderPath + task.filename
	log.Printf("%s\n", fullPath)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		log.Fatal("Error: {}", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		task.Tasks = append(task.Tasks, line)
	}

	for _, line := range task.Tasks {
		fmt.Printf("%s\n", line)
	}

}

func (task *TodayTasks) PrintFile() {
	for _, line := range task.Tasks {
		println(line)
	}
}

func (task *TodayTasks) PrintTodaysTasks() {

}

func (task *TodayTasks) ReplaceTodos() {
	tmp, err := os.OpenFile("tmpTodos.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error creating tmp file ", err)
	}

	defer os.Remove(tmp.Name())
	for _, val := range task.Tasks {
		tmp.WriteString(val + "\n")
	}
	if err := tmp.Sync(); err != nil {
		log.Fatal("Error sync ", err)
	}
	if err := tmp.Close(); err != nil {
		log.Fatal("Error Close ", err)
	}

	os.Rename(tmp.Name(), task.FolderPath+task.filename)

}
