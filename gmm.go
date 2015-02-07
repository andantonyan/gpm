package gmm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
)

type gmm interface {
	setGoRoot
	setGoPath
	setGoPathLocale

	getGoRoot
	getGoPath
	getGoPathLocale

	init
	install
	installDependencies
	addDependency
}
