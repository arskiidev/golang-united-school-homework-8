package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var (
	actions = map[string]func(Arguments, io.Writer) error{
		"list":     getItemsList,
		"add":      addItem,
		"remove":   removeUser,
		"findById": findById,
	}
)

func checkFileExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func getJsonFileContent(path string, jsonContent *[]Item) error {
	err := checkFileExist(path)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, jsonContent)
}

func getItemsList(args Arguments, writer io.Writer) error {
	err := checkFileExist(args["fileName"])
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(args["fileName"])
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func addItem(args Arguments, writer io.Writer) error {
	jsonContent := []Item{}
	var itemContent Item
	if args["item"] == "" {
		return errors.New("-item flag has to be specified")
	}
	err := getJsonFileContent(args["fileName"], &jsonContent)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(args["item"]), &itemContent)
	if err != nil {
		return err
	}
	for _, i := range jsonContent {
		if i.Id == itemContent.Id {
			msg := fmt.Sprintf("Item with id %v already exists", itemContent.Id)
			writer.Write([]byte(msg))
			return nil
		}
	}
	jsonContent = append(jsonContent, itemContent)
	data, err := json.Marshal(jsonContent)
	if err != nil {
		return nil
	}
	err = ioutil.WriteFile(args["fileName"], data, 0644)
	return err
}

func removeUser(args Arguments, writer io.Writer) error {
	jsonContent := []Item{}
	isNotExist := true
	var data []byte
	if args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}
	err := getJsonFileContent(args["fileName"], &jsonContent)
	if err != nil {
		return err
	}
	for i := len(jsonContent) - 1; i >= 0; i-- {
		if jsonContent[i].Id == args["id"] {
			jsonContent = append(jsonContent[:i], jsonContent[i+1:]...)
			isNotExist = false
		}
	}
	if isNotExist {
		msg := []byte(fmt.Sprintf("Item with id %v not found", args["id"]))
		writer.Write(msg)
	} else {
		data, err = json.Marshal(jsonContent)
	}
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(args["fileName"], data, 0644)
	return err
}

func findById(args Arguments, writer io.Writer) error {
	jsonContent := []Item{}
	isNotExists := true
	var data []byte
	var targetItem Item
	if args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}
	err := getJsonFileContent(args["fileName"], &jsonContent)
	if err != nil {
		return err
	}
	for _, item := range jsonContent {
		if item.Id == args["id"] {
			targetItem = item
			isNotExists = false
		}
	}
	if isNotExists {
		data = []byte("")
	} else {
		data, err = json.Marshal(targetItem)
	}
	writer.Write(data)
	return err
}

func parseArgs() Arguments {
	operation := flag.String("operation", "", "operation type")
	item := flag.String("item", "", "user data")
	fileName := flag.String("fileName", "", "name of the file")
	id := flag.String("id", "", "id of user")
	flag.Parse()
	return Arguments{"operation": *operation, "item": *item, "fileName": *fileName, "id": *id}
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}
	if _, ok := actions[args["operation"]]; !ok {
		return fmt.Errorf("Operation %v not allowed!", args["operation"])
	}
	err := actions[args["operation"]](args, writer)
	return err
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
