package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   uint   `json:"age"`
}

const (
	chmod = 0644
)

func Perform(args Arguments, writer io.Writer) error {

	var users []User

	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}

	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, chmod)
	if err != nil {
		return err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(fileContent) != 0 {
		json.Unmarshal(fileContent, &users)
	}

	switch args["operation"] {
	case "list":
		writer.Write(fileContent)

	case "add":
		if args["item"] == "" {
			return errors.New("-item flag has to be specified")
		}
		addUser(users, args["item"], file, writer)

	case "findById":
		if args["id"] == "" {
			return errors.New("-id flag has to be specified")
		}

		findById(users, args["id"], writer)

	case "remove":
		if args["id"] == "" {
			return errors.New("-id flag has to be specified")
		}

		removeUser(users, args["id"], file, writer)

	default:
		err := fmt.Sprintf("Operation %s not allowed!", args["operation"])
		return errors.New(err)
	}

	return nil
}

func getById(users []User, id string) ([]byte, string, error) {
	for _, v := range users {
		if v.Id == id {
			bytes, err := json.Marshal(v)
			if err != nil {
				return nil, "", err
			}
			return bytes, fmt.Sprintf("Item with id %s already exists", id), nil
		}
	}
	return nil, fmt.Sprintf("Item with id %s not found", id), nil
}

func addUser(users []User, item string, file *os.File, writer io.Writer) {
	var user User
	err := json.Unmarshal(json.RawMessage(item), &user)
	if err != nil {
		panic(err)
	}

	u, out, err := getById(users, user.Id)
	if err != nil {
		panic(err)
	}

	if out != "" && len(u) != 0 {
		writer.Write([]byte(out))
		return
	}

	users = append(users, user)
	fileContent, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}
	if _, err := file.WriteAt(fileContent, 0); err != nil {
		panic(err)
	}
	writer.Write(fileContent)
}

func findById(users []User, id string, writer io.Writer) {
	user, _, err := getById(users, id)
	if err != nil {
		panic(err)
	}
	writer.Write(user)
}

func removeUser(users []User, id string, file *os.File, writer io.Writer) {
	u, out, err := getById(users, id)
	if err != nil {
		panic(err)
	}

	if out != "" && len(u) == 0 {
		writer.Write([]byte(out))
		return
	}

	var newUsers []User
	for i, v := range users {
		if v.Id == id {
			newUsers = append(users[:i], users[i+1:]...)
		}
	}
	bytes, err := json.Marshal(newUsers)
	if err != nil {
		panic(err)
	}
	err = file.Truncate(0)
	if err != nil {
		panic(err)
	}
	file.WriteAt(bytes, 0)
	writer.Write(bytes)
}

func parseArgs() Arguments {
	operation := flag.String("operation", "", "Type of operations like add, list, findById or remove")
	item := flag.String("item", "", "Item example: {\"id\": \"1\", \"email\": \"email@test.com\", \"age\": 23}")
	fileName := flag.String("fileName", "", "Path to json database file")
	id := flag.String("id", "", "Record Id")

	flag.Parse()

	return Arguments{
		"operation": *operation,
		"item":      *item,
		"fileName":  *fileName,
		"id":        *id,
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
