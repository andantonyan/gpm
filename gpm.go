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

const (
	HEADER    = "\033[95m"
	INFO      = "\033[94m"
	SUCCESS   = "\033[92m"
	WARNING   = "\033[93m"
	ERROR     = "\033[91m"
	ENDC      = "\033[0m"
	BOLD      = "\033[1m"
	UNDERLINE = "\033[4m"

	MODULE_DIR = ""

	HELP_MESSAGE = "gpm init, initialization \n" +
		"gpm -g, use system GOPATH \n" +
		"gpm -i, install dependencies \n" +
		"gpm -i <name>, install package \n" +
		"gpm -i -s <name>, install package and save \n" +
		"gpm -e <command>, execute installed go binary \n" +
		"gpm -c <command>, run system command 'ex. gpm -c go build'"
)

type GPM struct{}

type packageJson struct {
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
	isInStrings(str string, list []string) bool
	removeFromSilce(str string, list []string) []string
	checkErr(err error)

	install(name string)
	installDependencies()
	saveDependency(name string)

	headerMessage(message string)
	successMessage(message string)
	infoMessage(message string)
	boldMessage(message string)
	underlineMessage(message string)
	warningMessage(message string)
	errorMessage(err error)
}

func (g GPM) headerMessage(message string) {
	fmt.Print(HEADER, message, ENDC)
}

func (g GPM) successMessage(message string) {
	fmt.Print(SUCCESS, message, ENDC)
}

func (g GPM) infoMessage(message string) {
	fmt.Print(INFO, message, ENDC)
}

func (g GPM) boldMessage(message string) {
	fmt.Print(BOLD, message, ENDC)
}

func (g GPM) underlineMessage(message string) {
	fmt.Print(UNDERLINE, message, ENDC)
}

func (g GPM) warningMessage(message string) {
	fmt.Print(WARNING, "warning:", message, ENDC)
}

func (g GPM) errorMessage(err error) {
	fmt.Print("\n", ERROR, "error: ", err, ENDC, "\n")
}

func (g GPM) checkErr(err error) {
	if err != nil {
		g.errorMessage(err)
		os.Exit(0)
	}
}

func (g GPM) getGoPath() string {
	return os.Getenv("GOPATH")
}

func (g GPM) getGoRoot() string {
	return os.Getenv("GOROOT")
}

func (g GPM) setGoPathTmp(path string) error {
	err := os.Setenv("GOPATH", path)
	g.checkErr(err)
	return err
}

func (g GPM) execCmd(cmd string, wg *sync.WaitGroup) []byte {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	g.checkErr(err)

	wg.Done()
	return out
}

func (g GPM) runBinary(name string) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	out := g.execCmd(g.getGoPath()+"/bin/"+name, wg)
	g.successMessage(string(out))
	wg.Wait()
}

func (g GPM) isInStrings(str string, list []string) bool {
	for _, part := range list {
		if part == str {
			return true
		}
	}
	return false
}

func (g GPM) removeFromSilce(str string, list []string) []string {
	s := list
	for i, part := range list {
		if part == str {
			s = append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (g GPM) install(name string) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	g.execCmd("go get "+name, wg)
	wg.Wait()
}

func (g GPM) installDependencies() {
	var packageFile packageJson

	file, err := ioutil.ReadFile("package.json")
	g.checkErr(err)

	err = json.Unmarshal(file, &packageFile)
	g.checkErr(err)

	for _, el := range packageFile.Dependencies {
		g.install(el)
	}
}

func (g GPM) saveDependency(name string) {
	var packageFile packageJson

	file, err := ioutil.ReadFile("package.json")
	g.checkErr(err)

	err = json.Unmarshal(file, &packageFile)
	g.checkErr(err)

	if !g.isInStrings(name, packageFile.Dependencies) {
		packageFile.Dependencies = append(packageFile.Dependencies, name)
	}

	packageByte, err := json.MarshalIndent(packageFile, "", "\t")
	g.checkErr(err)

	err = ioutil.WriteFile("package.json", packageByte, 0644)
	g.checkErr(err)
}

func (g GPM) init() {
	g.headerMessage("Press ^C at any time to quit.\n")

	var packageFile packageJson

	pwd, err := os.Getwd()
	g.checkErr(err)
	pwdSlice := strings.Split(pwd, "/")
	packageFile.Name = pwdSlice[len(pwdSlice)-1]
	g.infoMessage("Name: (" + packageFile.Name + ") ")
	fmt.Scanln(&packageFile.Name)

	file, err := ioutil.ReadFile("README.md")
	// g.checkErr(err)
	packageFile.Description = string(file)

	packageFile.Version = "0.0.1"
	g.infoMessage("Version: (" + packageFile.Version + ") ")
	fmt.Scanln(&packageFile.Version)

	usr, err := user.Current()
	g.checkErr(err)
	packageFile.Author = usr.Username
	g.infoMessage("Author: (" + packageFile.Author + ") ")
	fmt.Scanln(&packageFile.Author)

L1:
	isOk := "yes"
	g.boldMessage("Is this ok?: (" + isOk + ") ")
	fmt.Scanln(&isOk)

	if isOk == "no" {
		g.headerMessage("\nAborted")
		os.Exit(0)
	} else if isOk != "no" && isOk != "yes" {
		goto L1
	}

	packageByte, err := json.MarshalIndent(packageFile, "", "\t")
	g.checkErr(err)

	err = ioutil.WriteFile("package.json", packageByte, 0644)
	g.checkErr(err)
	g.successMessage("Done.\n")
}

func main() {
	var gpm GmmInterface = new(GPM)

	args := os.Args[1:]

	if !gpm.isInStrings("-g", args) {
		pwd, err := os.Getwd()
		gpm.checkErr(err)
		gpm.setGoPathTmp(pwd + MODULE_DIR)
	}

	args = gpm.removeFromSilce("-g", args)

	if gpm.isInStrings("-c", args) && len(args) != 1 {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		out := gpm.execCmd(strings.Join(args[1:], " "), wg)
		gpm.headerMessage(string(out))
		wg.Wait()
		os.Exit(0)
	}

	if gpm.isInStrings("init", args) {
		gpm.init()
		os.Exit(0)
	}

	if gpm.isInStrings("-i", args) && !gpm.isInStrings("-s", args) && len(args) == 1 {
		gpm.infoMessage("Installing dependencies...\n")
		gpm.installDependencies()
		gpm.successMessage("Done.\n")
		os.Exit(0)
	}

	if gpm.isInStrings("-i", args) && !gpm.isInStrings("-s", args) && len(args) == 2 {
		args = gpm.removeFromSilce("-i", args)
		gpm.infoMessage("Installing package " + args[0] + "...\n")
		gpm.install(args[0])
		gpm.successMessage("Done.\n")
		os.Exit(0)
	}

	if gpm.isInStrings("-i", args) && gpm.isInStrings("-s", args) && len(args) == 3 {
		args = gpm.removeFromSilce("-i", args)
		args = gpm.removeFromSilce("-s", args)

		gpm.infoMessage("Installing " + args[0] + "...\n")
		gpm.install(args[0])

		gpm.saveDependency(args[0])
		gpm.infoMessage("Package " + args[0] + " added in package.json.\n")

		gpm.successMessage("Done.\n")
		os.Exit(0)
	}

	args = gpm.removeFromSilce("-i", args)

	if gpm.isInStrings("-e", args) {
		args = gpm.removeFromSilce("-e", args)
		gpm.runBinary(strings.Join(args, " "))
		os.Exit(0)
	}

	gpm.headerMessage(HELP_MESSAGE + "\n")

}
