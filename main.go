package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sixteen/utils"
	"strconv"
	"strings"
)

type TaskModel struct {
	Id    string
	Title string
	Todos []string
}

func main() {
	prompt := promptui.Select{
		Label: "Refactoring",
		Items: []string{
			"list",
			"step",
			"switch",
			"delete",
			"create",
		},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "list":
		tasks := listTasks()
		fmt.Println(tasks)
	case "create":
		createNew()
	case "step":
		listSteps()
	default:
		validate()
	}
}

func listSteps() {

}

const task_path = "docs/refactoring/"

func listTasks() []TaskModel {
	files, err := ioutil.ReadDir(task_path)
	if err != nil {
		log.Fatal(err)
	}

	var tasks []TaskModel
	for _, f := range files {
		task, _ := ParseTask(task_path + f.Name())
		tasks = append(tasks, *task)
	}
	return tasks
}

func ParseTask(filePath string) (*TaskModel, error) {
	id := getIdFromFileName(filePath)
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	var todos []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
		var re = regexp.MustCompile(`\s-\s\[[ |x]\]\s(.*)`)

		for _, match := range re.FindAllString(scanner.Text(), -1) {
			all := re.ReplaceAllString(match, `$1`)
			todos = append(todos, all)
		}
	}

	file.Close()

	task := &TaskModel{
		Id:    id,
		Title: strings.ReplaceAll(txtlines[0], "# ", ""),
		Todos: todos,
	}

	return task, nil
}

func getIdFromFileName(filePath string) string {
	splitPath := strings.Split(filePath, "/")
	taskName := splitPath[len(splitPath)-1]
	id := taskName[0:utils.ID_LENGTH]
	return id
}

func createNew() {
	prompt := promptui.Prompt{
		Label: "title",
	}

	title, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	buildRefactoringFile(title)
}

func buildRefactoringFile(title string) {
	_ = os.MkdirAll("docs", os.ModePerm)
	_ = os.MkdirAll(task_path, os.ModePerm)

	fileName := utils.BuildFileName(utils.GenerateId(), title)
	_ = ioutil.WriteFile(task_path+"/"+fileName, []byte("# "+title+"\n\n"+" - [ ] todo"), 0644)
}

func validate() {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}
