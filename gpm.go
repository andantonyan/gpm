package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
    "os/exec"
    "sync"
    "strings"
    "os"
    "os/user"
)

type pachageList struct {
    Name string
    Description string
    Version string
    Dependencies []string
    Author string
}

func exe_cmd(cmd string, wg *sync.WaitGroup) {
    fmt.Println(cmd)
    parts := strings.Fields(cmd)
    out, err := exec.Command(parts[0],parts[1], parts[2]).Output()
    if err != nil {
        fmt.Println("error occured")
        log.Fatal(err)
        os.Exit(1)
    }
    fmt.Printf("%s", out)
    wg.Done()
}

func install_dependencies() {
  	var packages pachageList
	file, err := ioutil.ReadFile("package.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &packages)
	if err != nil {
		log.Fatal(err)
        os.Exit(1)
	}
    
    wg := new(sync.WaitGroup)
    for _, el := range packages.Dependencies {
        wg.Add(1)
        go exe_cmd("go get " + el, wg)
    }
    wg.Wait()
}

func install_package(packageSource string) {
    wg := new(sync.WaitGroup)
    wg.Add(1)
    go exe_cmd("go get " + packageSource, wg)
    wg.Wait()
}

func add_dependency(packageSource string) {
  	var packages pachageList
	file, err := ioutil.ReadFile("package.json")
	if err != nil {
		log.Fatal(err)
        os.Exit(1)
	}
	err = json.Unmarshal(file, &packages)
	if err != nil {
		log.Fatal(err)
        os.Exit(1)
	}
    
    packages.Dependencies = append(packages.Dependencies, packageSource)
    packagesByte, err := json.MarshalIndent(packages, "", "\t")
    
    if err != nil {
		log.Fatal(err)
        os.Exit(1)
	}
    
    err = ioutil.WriteFile("package.json", packagesByte, 0644)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
}

func string_in_slice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func initialization() {
    fmt.Println("Press ^C at any time to quit.")
    
    var packages pachageList
    
    packages.Name = ""
    fmt.Printf("Name: (%v) ", packages.Name)
    fmt.Scanln(&packages.Name)
    
    file, err := ioutil.ReadFile("README.md")
    packages.Description = string(file)
    
    packages.Version = "1.0.0"
    fmt.Printf("version: (%v) ", packages.Version)
    fmt.Scanln(&packages.Version)
    
    usr, err := user.Current()
    packages.Author = usr.Username
    fmt.Printf("Author: (%v) ", packages.Author)
    fmt.Scanln(&packages.Author)
    
    L1:
        is_it_ok:= "yes"
        fmt.Printf("Is this ok?: (%v) ", is_it_ok)
        fmt.Scanln(&is_it_ok)

        if is_it_ok == "no" {
            fmt.Printf("Aborted\n")
            os.Exit(0)
        } else if is_it_ok != "no" && is_it_ok != "yes" {
            goto L1
        }

    packagesByte, err := json.MarshalIndent(packages, "", "\t")
    
    if err != nil {
		log.Fatal(err)
        os.Exit(1)
	}
    
    err = ioutil.WriteFile("package.json", packagesByte, 0644)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
}

func main() {
    argsWithoutProg := os.Args[1:]
    
    if len(argsWithoutProg) == 1 && string_in_slice("init", argsWithoutProg) {
		fmt.Println("initialization...")
        initialization()
	} else if len(argsWithoutProg) == 1 && string_in_slice("install", argsWithoutProg) {
		fmt.Println("installing dependencies...")
        install_dependencies()
	} else if len(argsWithoutProg) == 2 && argsWithoutProg[0] == "install" {
        fmt.Println("installing package", argsWithoutProg[1])
        install_package(argsWithoutProg[1])
    } else if len(argsWithoutProg) == 3 && argsWithoutProg[0] == "install" && argsWithoutProg[2] == "save" {
        fmt.Println("installing package", argsWithoutProg[1])
        install_package(argsWithoutProg[1])
        add_dependency(argsWithoutProg[1])
        fmt.Println("package added in package.json", argsWithoutProg[1])
    } else {
        fmt.Println("init (initialization)\ninstall (install dependencies)\ninstall <packageSource> (install package)\ninstall <packageSource> save (install package and add source in package.json)")
    }
}
