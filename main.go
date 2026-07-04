package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func OpenFile(rootDir, fileName string) ([]*os.File, error) {
	//Locating All Files Related to the parameters set
	//FileName or Parts of the FileName
	//Root Directory, if not admin will just search home directory down
	//If admin, can choose which directory you'd like to start from
	var filePTR []*os.File
	fileName = strings.ToUpper(fileName)
	err := filepath.WalkDir(rootDir, func(search string, file fs.DirEntry, e error) error {
		if e != nil {
			if file != nil && file.IsDir() && e == fs.ErrPermission {
				return filepath.SkipDir
			}
		}
		if !strings.Contains(strings.ToUpper(file.Name()), fileName) {
			return nil
		}
		files, e := os.Open(search)
		if e != nil {
			return e
		}
		filePTR = append(filePTR, files)
		return nil

	})
	if len(filePTR) == 0 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Not Found")
	}
	return filePTR, nil
}

func UserInput() string {
	//User Terminal Input
	stdInput := bufio.NewReader(os.Stdin)
	input, err := stdInput.ReadString('\n')
	if err != nil {
		log.Fatalf("Error Reading Input %v\n", err)
	}
	input = strings.TrimSpace(input)
	return input
}

func IsAdmin() bool {
	//Checks if the Application is being ran as administrator
	if runtime.GOOS == "windows" {
		return windows.GetCurrentProcessToken().IsElevated()
	}
	return os.Getuid() == 0
}

func main() {
	//Variable Initialization and Setting Variable State
	var hash256STR []byte
	var hashMD5STR []byte
	fmt.Println("Which File Are You Searching For?")
	fileInQuestion := UserInput()
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error Processing Home Directory %v\n", err)
	}
	if IsAdmin() {
		//Runs the check if the user is admin
		fmt.Println("Please Set Root Directory.")
		homeDirectory = UserInput()
	}
	files, e := OpenFile(homeDirectory, fileInQuestion)
	if e != nil {
		log.Print("Error Reading Files or Directories\nPlease Ensure Item was Not Mispelled.\n")
	}
	for _, file := range files {
		//Loop that'll print out each file found and their respective hashes
		hasher256 := sha256.New()
		hasherMD5 := md5.New()
		defer file.Close()

		if _, err := io.Copy(hasher256, file); err != nil {
		}
		hash256STR = hasher256.Sum(nil)

		if _, err := io.Copy(hasherMD5, file); err != nil {
		}
		hashMD5STR = hasherMD5.Sum(nil)
		fmt.Printf("File Path: %s\n   SHA256 HASH: %x\n   MD5 Hash: %x\n", file.Name(), hash256STR, hashMD5STR)
	}
}
