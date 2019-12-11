package main

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"os"
	"sixteen/domain"
	. "sixteen/domain"
	"sixteen/utils"
	"strconv"
)

func main() {
	prompt := promptui.Select{
		Label: "Refactoring",
		Items: []string{
			"list",
			"step",
			"switch",
			"delete",
			"commit",
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
		tasks := domain.GetTasks()
		fmt.Println(tasks)
	case "create":
		createNew()
	case "step":
		tasks := domain.GetTasks()
		index := selectTask(tasks)
		name := getStepName(tasks[index])
		fmt.Println(name)
	case "commit":
		doCommit()
	default:
		validate()
	}
}

func doCommit() {
	prompt := promptui.Prompt{
		Label: "Commit Message",
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	utils.CommitByMessage("refactoring: " + result + "-" + utils.GenerateId())
}

func selectTask(tasks []TaskModel) int {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Id | cyan }}-{{ .Title | red }} {{ if eq .Done false }} ⌛ {{end}}",
		Inactive: "  {{ .Id | cyan }}-{{ .Title | red }} {{ if eq .Done false }} ⌛ {{end}}",
		Selected: "\U0001F336 {{ .Id | red | cyan }}-{{ if eq .Done false }} ⌛ {{end}}",
	}

	prompt := promptui.Select{
		Label:     "Refactoring",
		Templates: templates,
		Items:     tasks,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 0
	}

	return i
}

func getStepName(model TaskModel) string {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Content | red }} {{ if eq .Done true }} 👍 {{end}}",
		Inactive: "  {{ .Content | red }} {{ if eq .Done true }} 👍 {{end}}",
		Selected: "\U0001F336 {{ if eq .Done true }} 👍 {{end}}",
	}

	prompt := promptui.Select{
		Label:     model.Title,
		Templates: templates,
		Items:     model.Todos,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
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
	_ = os.MkdirAll(TASK_PATH, os.ModePerm)

	fileName := utils.BuildFileName(utils.GenerateId(), title)
	_ = ioutil.WriteFile(TASK_PATH+"/"+fileName, []byte("# "+title+"\n\n"+" - [ ] todo"), 0644)
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
