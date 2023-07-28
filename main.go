package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	Users = []*User{}
)

type User struct {
	Name     string `json:"name"`
	Balance  int    `json:"balance"`
	IsLoggIn bool   `json:"isLoggIn"`
}

func main() {
	fmt.Println("welcome to simple cli")
	fmt.Println("---------------------")

	if err := readIntoFile(); err != nil {
		panic(err)
	}

	if Users == nil {
		Users = make([]*User, 0)
	}
	args := os.Args

	for i := 1; i < len(args); i++ {
		arg := args[i]

		switch {
		case strings.HasPrefix(arg, "login="):
			name := strings.TrimPrefix(arg, "login=")
			login(name)
		case strings.HasPrefix(arg, "deposit="):
			amountString := strings.TrimPrefix(arg, "deposit=")
			amount, err := strconv.Atoi(amountString)
			if err != nil {
				panic(err)
			}
			deposit(amount)
		}
	}
}

func findUser(name string) *User {
	for _, user := range Users {
		if user.Name == name {
			return user
		}
	}
	return nil
}

func findUserIsLoggedIn() *User {
	for _, user := range Users {
		if user.IsLoggIn {
			return user
		}
	}
	return nil
}

func login(name string) {
	user := findUser(name)
	if user == nil {
		user = &User{
			Name:     name,
			IsLoggIn: true,
		}
		Users = append(Users, user)
		for _, u := range Users {
			if u.Name == name {
				u.IsLoggIn = true
			} else {
				u.IsLoggIn = false
			}
		}
	}

	fmt.Println("your balance is :", user.Balance)
	if err := writeIntoFile(); err != nil {
		panic(err)
	}

}

func deposit(amount int) {
	loggedInUser := findUserIsLoggedIn()
	if loggedInUser == nil {
		fmt.Println("You need to login before making a deposit.")
		return
	}
	loggedInUser.Balance += amount
	fmt.Printf("Successfully deposited %d into %s. New balance: %d\n", amount, loggedInUser.Name, loggedInUser.Balance)

	if err := writeIntoFile(); err != nil {
		panic(err)
	}
}

// func withdraw(amount int) {
// 	loggedInUser := findUserIsLoggedIn()
// 	if loggedInUser == nil {
// 		fmt.Println("You need to login before making a withdraw.")
// 		return
// 	}
// 	if loggedInUser.Balance < amount {
// 		fmt.Println("your balance less than amount")
// 		return
// 	}
// 	loggedInUser.Balance -= amount
// 	fmt.Printf("Successfully withdraw %d into %s. New balance: %d\n", amount, loggedInUser.Name, loggedInUser.Balance)

// 	if err := writeIntoFile(); err != nil {
// 		panic(err)
// 	}
// }

func readIntoFile() error {
	file, err := os.Open("users.json")
	if err != nil {
		if os.IsNotExist(err) {
			Users = make([]*User, 0)
			return nil
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&Users)
}

func writeIntoFile() error {
	writer, err := os.Create("users.json")
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	return json.NewEncoder(writer).Encode(&Users)
}
