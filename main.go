package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
)

type GMM struct{}

type moduleJson struct {
	Name         string
	Description  string
	Version      string
	Dependencies []string
	Author       string
}

type GmmInterface interface {
	init()

	getGoRoot() string
	getGoPath() string

	setGoPathTmp(path string) error

	execCmd(cmd string, wg *sync.WaitGroup) []byte
	runBinary(name string)
	checkErr(err error)

	install(name string)
	installDependencies()
	saveDependency(name string)

	headerMessage(message string)
	successMessage(message string)
	infoMessage(message string)
	boldMessage(message string)
	warningMessage(message string)
	errorMessage(err error)
}

const (
	// Colors
	HEADER    = "\033[95m"
	INFO      = "\033[94m"
	SUCCESS   = "\033[92m"
	WARNING   = "\033[93m"
	ERROR     = "\033[91m"
	ENDC      = "\033[0m"
	BOLD      = "\033[1m"
	UNDERLINE = "\033[4m"
	// Paths
	MODULE_DIR = "go_modules"
)

func (g GMM) headerMessage(message string) {
	fmt.Println(HEADER, message, ENDC)
}

func (g GMM) successMessage(message string) {
	fmt.Println(SUCCESS, message, ENDC)
}

func (g GMM) infoMessage(message string) {
	fmt.Print(INFO, message, ENDC)
}

func (g GMM) boldMessage(message string) {
	fmt.Print(BOLD, message, ENDC)
}

func (g GMM) warningMessage(message string) {
	fmt.Println(WARNING, "warning", message, ENDC)
}

func (g GMM) errorMessage(err error) {
	fmt.Println(ERROR, "error:", err, ENDC)
}

func (g GMM) checkErr(err error) {
	if err != nil {
		g.errorMessage(err)
		os.Exit(0)
	}
}

func (g GMM) getGoPath() string {
	return os.Getenv("GOPATH")
}

func (g GMM) getGoRoot() string {
	return os.Getenv("GOROOT")
}

func (g GMM) setGoPathTmp(path string) error {
	err := os.Setenv("GOPATH", path)
	g.checkErr(err)
	return err
}

func (g GMM) execCmd(cmd string, wg *sync.WaitGroup) []byte {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	g.checkErr(err)

	wg.Done()
	return out
}

func (g GMM) runBinary(name string) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	out := g.execCmd(g.getGoPath()+"/bin/"+name, wg)
	g.headerMessage(string(out))
	wg.Wait()
}

func (g GMM) install(name string) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	g.execCmd("go get "+name, wg)
	wg.Wait()
}

func (g GMM) installDependencies() {
	var module moduleJson

	file, err := ioutil.ReadFile("module.json")
	g.checkErr(err)

	err = json.Unmarshal(file, &module)
	g.checkErr(err)

	for _, el := range module.Dependencies {
		g.install(el)
	}
}

func (g GMM) saveDependency(name string) {
	var module moduleJson

	file, err := ioutil.ReadFile("module.json")
	g.checkErr(err)

	err = json.Unmarshal(file, &module)
	g.checkErr(err)

	module.Dependencies = append(module.Dependencies, name)

	moduleByte, err := json.MarshalIndent(module, "", "\t")
	g.checkErr(err)

	err = ioutil.WriteFile("module.json", moduleByte, 0644)
	g.checkErr(err)
}

func (g GMM) init() {
	g.headerMessage("Press ^C at any time to quit.")

	var module moduleJson

	pwd, err := os.Getwd()
	g.checkErr(err)
	pwdSlice := strings.Split(pwd, "/")
	module.Name = pwdSlice[len(pwdSlice)-1]
	g.infoMessage(" Name: (" + module.Name + ") ")
	fmt.Scanln(&module.Name)

	file, err := ioutil.ReadFile("README.md")
	// g.checkErr(err)
	module.Description = string(file)

	module.Version = "0.0.1"
	g.infoMessage(" Version: (" + module.Version + ") ")
	fmt.Scanln(&module.Version)

	usr, err := user.Current()
	g.checkErr(err)
	module.Author = usr.Username
	g.infoMessage(" Author: (" + module.Author + ") ")
	fmt.Scanln(&module.Author)

L1:
	isOk := "yes"
	g.boldMessage(" Is this ok?: (" + isOk + ") ")
	fmt.Scanln(&isOk)

	if isOk == "no" {
		g.headerMessage("Aborted")
		os.Exit(0)
	} else if isOk != "no" && isOk != "yes" {
		goto L1
	}

	moduleByte, err := json.MarshalIndent(module, "", "\t")
	g.checkErr(err)

	err = ioutil.WriteFile("module.json", moduleByte, 0644)
	g.checkErr(err)
	g.successMessage("Done.")
}

func main() {

}
