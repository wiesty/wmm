package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/gookit/color"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Hints      []string `json:"hints"`
	InstallURL []string `json:"installurl"`
	ModsURL    string   `json:"modsurl"`
}

var config Config

func main() {
	// Set terminal title to "wmm"
	setTerminalTitle("wmm")

	clearTerminal()
	fmt.Println(`
           _             __                                             
 _    __  (_) ___   ___ / /_  __ __  ___                                
| |/|/ / / / / -_) (_-</ __/ / // / (_-<                                
|__,__/ /_/  \__/ /___/\__/  \_, / /___/                                
                            /___/                                       
                 __        _              __          __   __           
  __ _  ___  ___/ /       (_)  ___   ___ / /_ ___ _  / /  / / ___   ____ 
 /  ' \/ _ \/ _  /       / /  / _ \ (_-</ __// _  / / /  / / / -_) / __/ 
/_/_/_/\___/\_,_/       /_/  /_//_//___/\__/ \_,_/ /_/  /_/  \__/ /_/    
	`)

	loadConfig()

	options := []string{
		"Install Mods",
		"Restore latest MC-Backup",
		"List Backups",
		"Show Config",
		"Create Backup",
		"Exit",
	}

	var result string
	prompt := &survey.Select{
		Message: "Choose an option:",
		Options: options,
	}

	err := survey.AskOne(prompt, &result)
	if err != nil {
		color.Red.Println("Selection error: %v", err)
		return
	}

	switch result {
	case "Install Mods":
		installMods()
	case "Restore latest MC-Backup":
		restoreBackup()
	case "List Backups":
		listBackups()
	case "Show Config":
		showConfig()
	case "Create Backup":
		createBackup()
	case "Exit":
		color.Cyan.Println("Exiting... Goodbye!")
		os.Exit(0)
	}
}

func setTerminalTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}

func loadConfig() {
	file, err := os.Open("wmm.json")
	if err != nil {
		color.Red.Println("Failed to load config: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		color.Red.Println("Error parsing config: %v", err)
		os.Exit(1)
	}
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func installMods() {
	clearTerminal()

	color.Cyan.Println("Installing mods...")
	mcPath := choosePath(".minecraft")

	backupMinecraft(mcPath)

	if _, err := os.Stat(mcPath); os.IsNotExist(err) {
		color.Red.Println("The .minecraft folder was not found!")
		return
	}

	for _, hint := range config.Hints {
		color.Red.Println(hint)
	}

	if len(config.InstallURL) == 0 {
		color.Red.Println("No install URL found. Skipping installer step...")
	} else {
		for _, url := range config.InstallURL {
			color.Cyan.Println("Downloading and running installer from ", url)
			downloadAndRunInstaller(url)
		}
	}

	downloadAndExtractMods(mcPath, config.ModsURL)

	color.Green.Println("Installation complete!")
	showPostTaskMenu()
}

func downloadAndRunInstaller(url string) {
	tempPath := filepath.Join(os.TempDir(), "installer.exe")

	if _, err := os.Stat(tempPath); err == nil {
		os.Remove(tempPath)
	}

	err := downloadFile(url, tempPath)
	if err != nil {
		color.Red.Println("Error downloading installer: %v", err)
		return
	}

	color.Cyan.Println("Running installer...")
	cmd := exec.Command(tempPath)
	cmd.Start()

	fmt.Println("Press Enter after installation is complete...")
	fmt.Scanln()
}


func downloadAndExtractMods(mcPath, modsUrl string) {
	color.Cyan.Println("Downloading mods from ", modsUrl)
	modsZip := filepath.Join(os.TempDir(), "mods.zip")
	
	if _, err := os.Stat(modsZip); err == nil {
		os.Remove(modsZip)
	}

	err := downloadFile(modsUrl, modsZip)
	if err != nil {
		color.Red.Println("Error downloading mods: %v", err)
		return
	}

	color.Cyan.Println("Extracting mods...")
	err = unzip(modsZip, mcPath)
	if err != nil {
		color.Red.Println("Error extracting mods: %v", err)
		return
	}
	color.Green.Println("Mods successfully installed!")
}


func showConfig() {
	clearTerminal()
	color.Cyan.Println("Current Configuration:")
	fmt.Printf("Hints: %v\n", config.Hints)
	fmt.Printf("Install URL(s): %v\n", config.InstallURL)
	fmt.Printf("Mods URL: %s\n", config.ModsURL)
	showPostTaskMenu()
}

func createBackup() {
	color.Cyan.Println("Creating Minecraft backup...")
	mcPath := choosePath(".minecraft")
	backupMinecraft(mcPath)
}

func backupMinecraft(mcPath string) {
	backupPath := mcPath + "BACKUP"

	if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
		currentTime := time.Now().Format("2006-01-02_15-04-05")
		backupPath = mcPath + "BACKUP_" + currentTime
	}

	color.Cyan.Println("Creating backup...")

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()
	time.Sleep(2 * time.Second)
	err := copyDir(mcPath, backupPath)
	s.Stop()

	if err != nil {
		color.Red.Println("Error creating backup: %v", err)
		os.Exit(1)
	}
	color.Green.Println("Backup successfully created at ", backupPath)
}

func restoreBackup() {
	color.Cyan.Println("Restoring Minecraft backup...")

	mcPath := choosePath(".minecraft")
	backupPath := mcPath + "BACKUP"

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		color.Red.Println("Backup folder not found!")
		return
	}

	color.Cyan.Println("Deleting current .minecraft folder...")
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()

	err := os.RemoveAll(mcPath)
	if err != nil {
		s.Stop()
		color.Red.Println("Error deleting .minecraft folder: %v", err)
		return
	}
	s.Stop()
	color.Green.Println("The .minecraft folder was successfully deleted.")

	color.Cyan.Println("Copying the backup...")
	s.Start()

	err = copyDir(backupPath, mcPath+"Temp")
	if err != nil {
		s.Stop()
		color.Red.Println("Error copying the backup: %v", err)
		os.Exit(1)
	}
	s.Stop()
	color.Green.Println("Backup copied successfully.")

	color.Cyan.Println("Activating the backup...")
	s.Start()
	err = os.Rename(mcPath+"Temp", mcPath)
	if err != nil {
		s.Stop()
		color.Red.Println("Error renaming the folder: %v", err)
		return
	}
	s.Stop()

	color.Green.Println("Backup successfully restored.")
	showPostTaskMenu()
}

func listBackups() {
	color.Cyan.Println("Listing available backups...")
	mcPath := choosePath(".minecraft")
	backups := []string{}

	backupDir := filepath.Dir(mcPath)
	files, err := os.ReadDir(backupDir)
	if err != nil {
		color.Red.Println("Error reading the directory: %v", err)
		return
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "BACKUP") {
			backups = append(backups, file.Name())
		}
	}

	if len(backups) == 0 {
		color.Red.Println("No backups found.")
		return
	}

	for {
		prompt := &survey.Select{
			Message: "Choose a backup:",
			Options: backups,
		}

		var selectedBackup string
		err := survey.AskOne(prompt, &selectedBackup)
		if err != nil {
			color.Red.Println("Error selecting backup: %v", err)
			return
		}

		promptAction := &survey.Select{
			Message: fmt.Sprintf("Selected backup: %s. What do you want to do?", selectedBackup),
			Options: []string{"Restore", "Delete", "Return to Main Menu"},
		}

		var action string
		err = survey.AskOne(promptAction, &action)
		if err != nil {
			color.Red.Println("Error selecting action: %v", err)
			return
		}

		switch action {
		case "Restore":
			restoreFromBackup(filepath.Join(backupDir, selectedBackup))
			return
		case "Delete":
			deleteBackup(filepath.Join(backupDir, selectedBackup))
			return
		case "Return to Main Menu":
			return
		}
	}
}

func deleteBackup(backupPath string) {
	color.Cyan.Println("Deleting backup...")

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()

	err := os.RemoveAll(backupPath)
	s.Stop()

	if err != nil {
		color.Red.Println("Error deleting backup: %v", err)
	} else {
		color.Green.Println("Backup deleted successfully.")
	}

	showPostTaskMenu()
}

func restoreFromBackup(backupPath string) {
	color.Cyan.Println("Restoring backup from %s...", backupPath)

	mcPath := choosePath(".minecraft")

	color.Cyan.Println("Deleting current .minecraft folder...")
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()

	err := os.RemoveAll(mcPath)
	if err != nil {
		s.Stop()
		color.Red.Println("Error deleting .minecraft folder: %v", err)
		return
	}
	s.Stop()

	color.Cyan.Println("Copying the backup...")
	s.Start()

	err = copyDir(backupPath, mcPath+"Temp")
	if err != nil {
		s.Stop()
		color.Red.Println("Error copying the backup: %v", err)
		os.Exit(1)
	}
	s.Stop()
	color.Green.Println("Backup copied successfully.")

	color.Cyan.Println("Activating the backup...")
	s.Start()
	err = os.Rename(mcPath+"Temp", mcPath)
	if err != nil {
		s.Stop()
		color.Red.Println("Error renaming the folder: %v", err)
		return
	}
	s.Stop()

	color.Green.Println("Backup successfully restored.")
	showPostTaskMenu()
}

func showPostTaskMenu() {
	options := []string{"Back to Main Menu", "Exit"}

	prompt := &survey.Select{
		Message: "What would you like to do next?",
		Options: options,
	}

	var result string
	err := survey.AskOne(prompt, &result)
	if err != nil {
		color.Red.Println("Selection error: %v", err)
		os.Exit(1)
	}

	switch result {
	case "Back to Main Menu":
		main()
	case "Exit":
		color.Cyan.Println("Exiting program...")
		os.Exit(0)
	}
}

func choosePath(folder string) string {
	usr, _ := user.Current()
	defaultPath := filepath.Join(usr.HomeDir, "AppData", "Roaming", folder)
	if folder == "Downloads" {
		defaultPath = filepath.Join(usr.HomeDir, "Downloads")
	}

	prompt := &survey.Select{
		Message: fmt.Sprintf("Choose the %s path:", folder),
		Options: []string{"Default Path", "Custom Path"},
	}

	var result string
	err := survey.AskOne(prompt, &result)
	if err != nil {
		color.Red.Println("Selection error: %v", err)
		os.Exit(1)
	}

	if result == "Custom Path" {
		return promptForInput(fmt.Sprintf("Enter the custom %s path: ", folder))
	}

	return defaultPath
}

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func copyDir(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, src)
		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			os.MkdirAll(destPath, info.Mode())
		} else {
			err := copyFile(path, destPath)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func promptForInput(label string) string {
	prompt := &survey.Input{
		Message: label,
	}

	var result string
	err := survey.AskOne(prompt, &result)
	if err != nil {
		color.Red.Println("Input error: %v", err)
		os.Exit(1)
	}
	return result
}
