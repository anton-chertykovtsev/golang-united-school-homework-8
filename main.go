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
	Age   int64  `json:"age"`
}

const (
	chmod = 0644
)

func Perform(args Arguments, writer io.Writer) error {

	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}

	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch args["operation"] {
	case "list":
		file, err := os.OpenFile(args["fileName"], os.O_RDONLY|os.O_CREATE, chmod)
		defer file.Close()

		if err != nil {
			return err
		}

		bytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		writer.Write(bytes)

	case "add":
		if args["item"] == "" {
			return errors.New("-item flag has to be specified")
		}

		file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, chmod)
		defer file.Close()

		if err != nil {
			return err
		}

		var users []User
		var user User

		bytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		if len(bytes) != 0 {
			err := json.Unmarshal(bytes, &users)
			if err != nil {
				return err
			}
		}

		fmt.Println(users)

		bytes, err = json.Marshal(args["item"])
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &user)
		if err != nil {
			return err
		}

		fmt.Println(user.Id)

		// if err != nil {
		// 	return err
		// }
		// fmt.Println(string(bytes))
		// file.Write(bytes)

	default:
		err := fmt.Sprintf("Operation %s not allowed!", args["operation"])
		return errors.New(err)
	}

	return nil
}

func parseArgs() Arguments {
	operation := flag.String("operation", "", "Type of operations like add, list, findById or remove")
	item := flag.String("item", "", "Item example: {\"id\": \"1\", \"email\": \"email@test.com\", \"age\": 23}")
	fileName := flag.String("fileName", "", "Path to json database file")
	flag.Parse()

	return Arguments{
		"operation": *operation,
		"item":      *item,
		"fileName":  *fileName,
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
