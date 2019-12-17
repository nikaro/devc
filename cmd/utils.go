package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// GetConfig get setting from devcontainer.json file
func GetConfig(name string) string {
	var config string

	// open file
	jsonFile, err := os.Open(rootPath + ".devcontainer/devcontainer.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	// load content
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	// check if set
	if name != "name" && result[name] == nil {
		fmt.Println(name + " is not set")
		os.Exit(1)
	}

	// process values
	switch name {
	case "dockerComposeFile":
		dockerComposeFile := rootPath + ".devcontainer/" + result[name].(string)
		// check if exists
		if _, err := os.Stat(dockerComposeFile); os.IsNotExist(err) {
			fmt.Println("docker-compose file not found: " + dockerComposeFile)
			os.Exit(1)
		}
		config = dockerComposeFile
	case "name":
		var projectName string
		// set name to current path if not set
		if result[name] == nil {
			path, _ := os.Getwd()
			dirs := strings.Split(path, "/")
			projectName = dirs[len(dirs)-1]
		} else {
			projectName = result[name].(string)
		}
		config = projectName + "_devcontainer"
	default:
		config = result[name].(string)
	}

	return config
}

// Run the given command
func Run(command []string, verbose bool) error {
	if verbose {
		fmt.Println(strings.Join(command, " "))
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
