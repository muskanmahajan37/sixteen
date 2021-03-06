package utils

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type CommitMessage struct {
	Rev     string
	Author  string
	Date    string
	Message string
	Changes []FileChange
	Task    LogTaskModel
}

type LogTaskModel struct {
	Id string
}

type FileChange struct {
	Added   int
	Deleted int
	File    string
}

var currentCommitMessage CommitMessage
var currentFileChanges []FileChange
var commitMessages []CommitMessage

func BuildCommitMessages() []CommitMessage {
	historyArgs := []string{"log", "--pretty=format:[%h] %aN %ad %s", "--date=iso", "--numstat"}
	cmd := exec.Command("git", historyArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	splitStr := strings.Split(string(out), "\n")
	for _, str := range splitStr {
		parseLog(str)
	}

	return commitMessages
}

type TopAuthor struct {
	Name        string
	CommitCount int
	LineCount   int
}

func parseLog(text string) {
	rev := `\[([\d|a-f]{5,8})\]`
	author := `(.*?)\s\d{4}-\d{2}-\d{2}`
	date := `\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\s\+\d{4}`
	changes := `([\d-])*\t([\d-]*)\t(.*)`

	revReg := regexp.MustCompile(rev)
	authorReg := regexp.MustCompile(author)
	dateReg := regexp.MustCompile(date)
	changesReg := regexp.MustCompile(changes)

	taskModel := &LogTaskModel{
		Id: "",
	}

	allString := revReg.FindAllString(text, -1)
	if len(allString) == 1 {
		str := ""

		id := revReg.FindStringSubmatch(text)
		str = strings.Split(text, id[0])[1]
		auth := authorReg.FindStringSubmatch(str)
		str = strings.Split(str, auth[1])[1]
		dat := dateReg.FindStringSubmatch(str)
		msg := strings.Split(str, dat[0])[1]
		msg = msg[1:]

		if strings.HasPrefix(msg, "refactoring:") && string(msg[len(msg) - ID_LENGTH - 1]) == "-" {
			taskModel.Id = msg[len(msg) - ID_LENGTH:len(msg)]
		}

		date := dat[0][0:len(dat[0])-6]
		currentCommitMessage = *&CommitMessage{id[1], auth[1], date, msg, nil, *taskModel}
	} else if changesReg.MatchString(text) {
		changes := changesReg.FindStringSubmatch(text)
		deleted, _ := strconv.Atoi(changes[2])
		added, _ := strconv.Atoi(changes[1])
		change := &FileChange{added, deleted, changes[3]}

		currentFileChanges = append(currentFileChanges, *change)
	} else {
		if currentCommitMessage.Rev != "" {
			currentCommitMessage.Changes = currentFileChanges
			commitMessages = append(commitMessages, currentCommitMessage)

			currentCommitMessage = *&CommitMessage{"", "", "", "", nil, *taskModel}
			currentFileChanges = nil
		}

	}
}
