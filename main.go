package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	// "sync"
)

type GmmInterface interface {
	// setGoRoot()
	// setGoPath()
	// setGoPathLocale()

	getGoRoot() string
	getGoPath() string
	// getGoPathLocale() string

	// init()
	// install()
	// installDependencies()
	// addDependency()

	execCmd(cmd string)

	checkErr(err error)

	warningMessage(message string)
	successMessage(message string)
	errorMessage(err error)
}

type GMM struct {
	GoPath string
	GoRoot string
}

const (
	HEADER    = "\033[95m"
	OKBLUE    = "\033[94m"
	OKGREEN   = "\033[92m"
	WARNING   = "\033[93m"
	FAIL      = "\033[91m"
	ENDC      = "\033[0m"
	BOLD      = "\033[1m"
	UNDERLINE = "\033[4m"
)

func (g GMM) successMessage(message string) {
	fmt.Println(OKGREEN, message, ENDC)
}

func (g GMM) warningMessage(message string) {
	fmt.Println(WARNING, message, ENDC)
}

func (g GMM) errorMessage(err error) {
	fmt.Println(FAIL, err, ENDC)
}

func (g GMM) checkErr(err error) {
	if err != nil {
		g.errorMessage(err)
		os.Exit(1)
	}
}

func (g GMM) getGoPath() string {
	return os.Getenv("GOPATH")
}

func (g GMM) getGoRoot() string {
	return os.Getenv("GOROOT")
}

// done

func (g GMM) execCmd(cmd string) {
	parts := strings.Fields(cmd)
	_, err := exec.Command(parts[0], parts[1], parts[2]).Output()
	g.checkErr(err)
}

func main() {
	// var gmm GmmInterface = new(GMM)

	// fmt.Println(gmm.getGoRoot())
	// fmt.Println(gmm.getGoPath())
}

// func (g GMM) execCmdWait(cmd string) {
// 	wg := new(sync.WaitGroup)
// 	wg.add(1)
// 	g.execCmd(cmd)
// 	wg.Wait()
// 	wg.Done()
// }
