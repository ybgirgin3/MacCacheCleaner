package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/getlantern/systray"
)

var __VERSION__ string = "0.0.1"

// var prodCommands = map[string]string{
// 	"root":            "",
// 	"unixRemove":      "rm -rf",
// 	"dotcachedirpath": filepath.Join(pwd, "DummyCacheFile"),
// }

/*
define false in prod !!
define DEBUG=1 as env var
ex: DEBUG=1 go run main.go
*/
var DEFAULT_DEBUG string = "0"

var (
	appName     string = "MacCacheCleaner"
	runningIcon string = " ♻️"
)

var Commands = map[string]string{
	"root":            "sudo",
	"unixRemove":      "rm -rf",
	"dotcachedirpath": ".cache",
}

func main() {
	DEBUG := os.Getenv("DEBUG")
	if len(DEBUG) == 0 {
		os.Setenv("DEBUG", DEFAULT_DEBUG)
	}

	systray.Run(onReady, onExit)
}

// helpers
func cacheFileDir(targetFilePath string) (map[string]string, error) {
	// var _files = []string{}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return map[string]string{}, err
	}
	fmt.Println(targetFilePath)

	dir := filepath.Join(homedir, Commands[targetFilePath])

	// get names off files
	// entries, err := os.ReadDir(dir)
	// if err != nil {
	// 	return nil, fmt.Errorf("errrrr %v", err)
	//
	// }
	// for _, e := range entries {
	// 	fmt.Println(e.Name(), e)
	// 	_files = append(_files, e.Name())
	// }

	// get len of all files
	files, _ := ioutil.ReadDir(dir)
	// returnString := fmt.Sprintf("%d found in %s", len(files), dir)
	var filePaths = []string{}

	for _, e := range files {
		fmt.Println(e.Name(), e)
		filePaths = append(filePaths, e.Name())
	}

	fmt.Println("files in cacheFileDir", filePaths)

	return map[string]string{
		"len":      strconv.Itoa(len(files)),
		"allFiles": fmt.Sprintf(strings.Join(filePaths[:], "\n")),
	}, nil
}

func onReady() {

	// ** DEFINITIONS **
	// ** define barebone of macos app
	systray.SetTitle(appName + runningIcon)
	systray.SetTooltip("Automatic Cache Cleaner For Mac")

	// ** buttons
	// .cache
	dirInfo, err := cacheFileDir("dotcachedirpath")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("dirInfooo", dirInfo)

	// dirinf, _ := dirInfo.(map[string]interface{})
	// jsonData, err := json.MarshalIndent(dirInfo, "", "  ")
	// if err != nil {
	// 	fmt.Println("erro amk", err)
	// }
	// info, _ := json.Marshal(dirInfo)

	dotcacheButtonName := fmt.Sprintf(".cache (%s file found)", dirInfo["len"])
	dotcacheButton := systray.AddMenuItem(dotcacheButtonName, dirInfo["allFiles"])

	// ** quit
	mQuit := systray.AddMenuItem("Quit", "Quit App")

	// *** HANDLERS ***
	dotcacheHandler := func(isSudoNeeded bool) error {
		var sudo string = ""

		if isSudoNeeded {
			sudo = Commands["root"]
		}

		if os.Getenv("DEBUG") == "1" {
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println("err", err)
			}
			Commands["dotcachedirpath"] = filepath.Join(pwd, "DummyCacheFile")
		} else if os.Getenv("DEBUG") == "0" {
			homedir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			Commands["dotcachedirpath"] = filepath.Join(homedir, ".cache")
			// Commands["dotcachedirpath"] = "~/.cache"
		}
		_commandString := fmt.Sprintf("%s %s %s", sudo, Commands["unixRemove"], Commands["dotcachedirpath"])
		// fmt.Printf("command %v", _commandString)
		if err := runTerminalCommand(_commandString); err != nil {
			return err
		}
		cacheFileDir("dotcachedirpath")
		return nil
	}

	quitHandler := func() error {
		systray.Quit()
		return nil
	}

	for {
		select {
		case <-dotcacheButton.ClickedCh:
			if err := dotcacheHandler(false); err != nil {
				fmt.Println("error in start handler", err)
			}

		case <-mQuit.ClickedCh:
			if err := quitHandler(); err != nil {
				fmt.Println("error in quit", err)
			}

		}
	}

}

func onExit() {
	// cleanup

}

func runTerminalCommand(command string) error {
	// fmt.Println("current running command", command)

	// print len and name of files
	// dir := strings.Fields(command)
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return nil
	}
	dir := parts[len(parts)-1]

	fmt.Printf("This files in %s will be gone!! in \n", dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("errrrr %v", err)

	}
	for _, e := range entries {
		fmt.Println(e.Name(), e)
	}

	// define command
	cmd := exec.Command("sh", "-c", command)

	// run command
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error while running", command, err)
		return fmt.Errorf("error occured while running command %s", err)
	}

	fmt.Println("command run successfully")

	return nil

}
