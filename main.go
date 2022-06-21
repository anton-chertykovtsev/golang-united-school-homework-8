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
	var user User

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

		item := json.RawMessage(args["item"])
		err = json.Unmarshal(item, &user)
		if err != nil {
			return err
		}

		u, out, err := getById(users, user.Id)
		if err != nil {
			return err
		}

		if out != "" && len(u) != 0 {
			writer.Write([]byte(out))
			return nil
		}

		users = append(users, user)
		fileContent, err = json.Marshal(users)
		if err != nil {
			return err
		}
		if _, err := file.WriteAt(fileContent, 0); err != nil {
			return err
		}
		writer.Write(fileContent)

	case "findById":
		if args["id"] == "" {
			return errors.New("-id flag has to be specified")
		}

		user, _, err := getById(users, args["id"])
		if err != nil {
			return err
		}

		writer.Write(user)

	case "remove":
		if args["id"] == "" {
			return errors.New("-id flag has to be specified")
		}

		u, out, err := getById(users, args["id"])
		if err != nil {
			return nil
		}

		if out != "" && len(u) == 0 {
			writer.Write([]byte(out))
			return nil
		}

		var newUsers []User
		for i, v := range users {
			if v.Id == args["id"] {
				newUsers = append(users[:i], users[i+1:]...)
			}
		}
		bytes, _ := json.Marshal(newUsers)
		file.Truncate(0)
		file.WriteAt(bytes, 0)
		writer.Write(bytes)

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
