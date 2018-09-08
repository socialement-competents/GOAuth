package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/socialement-competents/goauth/database"
)

func runFile(filename string, client *database.Client) {
	query, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("running migration ", filename)

	if _, err := client.Connection.Exec(string(query)); err != nil {
		panic(fmt.Sprintf("migration failed: %v", err))
	}
}

func runDir(dirname string, client *database.Client) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(fmt.Sprintf("reading the %s directory failed: %v", dirname, err))
	}

	for _, f := range files {
		filepath := filepath.Join(dirname, f.Name())
		runFile(filepath, client)
	}
}

func main() {
	client, err := database.NewClient()
	if err != nil {
		fmt.Println("connecting to the database failed: ", err)
		return
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("please provide folders or sql files")
		return
	}

	for _, arg := range args {
		path := filepath.Join(arg)
		fi, err := os.Stat(path)
		if err != nil {
			fmt.Println("couldn't open file or directory: ", err)
			return
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			runDir(path, client)
		case mode.IsRegular():
			runFile(path, client)
		}
	}
}
