package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var IsInt = regexp.MustCompile(`^[0-9]+$`).MatchString

type HashcatConfig struct {
	PotPath         string `json:"potfile-path"`
	RestorePath     string `json:"restorefile-path"`
	Workload        int    `json:"workload-profile"`
	BitmapMax       int    `json:"bitmap-max"`
	OptimizedKernel int    `json:"optimized-kernel"`
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("error: not enough arguments provided to run")
		fmt.Println("usage: <hashcat hash-mode> <hashcat command>")
		fmt.Println("Files available on the server:", "\n")
		printFiles()

		os.Exit(0)

	}

	hashType := os.Args[1]
	hashCommand := prepareHashcatCommand(hashType, os.Args[2:])

	if IsInt(hashType) == false {
		fmt.Println("Input contained non-numeric characters: ", fmt.Sprint("%i", hashType))
		os.Exit(0)
	}

	commandHashcat(hashCommand)
}

// Parses configuration json file
func parseConf() map[string]interface{} {
	fileContent, err := os.Open(os.Getenv("CONF_DIR"))

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	defer fileContent.Close()
	byteResult, _ := ioutil.ReadAll(fileContent)
	var res map[string]interface{}
	json.Unmarshal([]byte(byteResult), &res)

	return res
}

// Prints the files available on the server
func printFiles() {
	err := filepath.Walk(os.Getenv("FILE_DIR"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("%s\n", path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

// Prepares a Hashcat command from CLI arguments
func prepareHashcatCommand(hashType string, args []string) string {
	conf := parseConf()
	kernel := "--optimized-kernel-enable" + " "
	base := " -m" + hashType + " "
	potpath := "--potfile-path=" + fmt.Sprintf("%v", conf["potfile-path"]) + hashType + ".potfile" + " "
	respath := "--restore-file-path=" + fmt.Sprintf("%v", conf["restorefile-path"]) + hashType + ".restore --session=" + hashType + " "
	workload := "--workload-profile=" + fmt.Sprintf("%v", conf["workload-profile"]) + " "

	if fmt.Sprintf("%v", conf["optimized-kernel"]) == "0" {
		kernel = " "
	}

	c := strings.Join(args, " ")
	fin := base + potpath + respath + workload + kernel + c
	return fin
}

// Starts Hashcat and connects std to cmd
func commandHashcat(hashCommand string) {
	cmd := exec.Command("hashcat", strings.Fields(hashCommand)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
