package setting

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Record struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type Zone struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Records []Record
}

var (
	appPath    string
	CustomPath string

	Cfg struct {
		CheckTime time.Duration

		ClientAPI   string
		ClientEmail string
		IPServer    string
		Zones       []Zone
	}
)

// execPath returns the executable path.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func init() {
	var err error
	if appPath, err = execPath(); err != nil {
		// log.Fatal(4, "fail to get app path: %v\n", err)
	}

	// Note: we don't use path.Dir here because it does not handle case
	//	which path starts with two "/" in Windows: "//psf/Home/..."
	appPath = strings.Replace(appPath, "\\", "/", -1)
}

// WorkDir returns absolute path of work directory.
func WorkDir() (string, error) {
	i := strings.LastIndex(appPath, "/")
	if i == -1 {
		return appPath, nil
	}
	return appPath[:i], nil
}

// NewContext initializes configuration context.
func NewContext() {
	workDir, err := WorkDir()
	if err != nil {
		log.Fatal("Fail to get work directory: %v", err)
	}

	CustomPath = workDir + "/custom"

	externalConf, err := ioutil.ReadFile(CustomPath + "/conf/conf.json")
	if err == nil {
		json.Unmarshal(externalConf, &Cfg)
	}
}
