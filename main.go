package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	runningIcon string = " W"
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
func cacheFileLen(targetFilePath string) (int, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return 0, err
	}
	fmt.Println(targetFilePath)
	dir := filepath.Join(homedir, Commands[targetFilePath])
	fmt.Println("cacheFileLen dir", dir)
	files, _ := ioutil.ReadDir(dir)
	// returnString := fmt.Sprintf("%d found in %s", len(files), dir)
	return len(files), nil
}

func onReady() {

	// ** DEFINITIONS **
	// ** define barebone of macos app
	systray.SetTitle(appName + runningIcon)
	systray.SetTooltip("Automatic Cache Cleaner For Mac")

	// ** buttons
	// .cache
	filelen, err := cacheFileLen("dotcachedirpath")
	if err != nil {
		fmt.Println(err)
	}
	dotcacheButtonName := fmt.Sprintf(".cache (%d file found)", filelen)
	dotcacheButton := systray.AddMenuItem(dotcacheButtonName, "Clean .cache file")

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
