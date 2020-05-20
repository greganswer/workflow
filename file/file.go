package file

import (
	"os"
	"time"
)

// Exists determines if the file exists.
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Touch a file. Return true if the file exists.
// Ref: https://golangbyexample.com/touch-file-golang/
func Touch(name string) (bool, error) {
	exists, err := Exists(name)
	if err != nil {
		return exists, err
	}
	if exists {
		currentTime := time.Now().Local()
		err = os.Chtimes(name, currentTime, currentTime)
		return exists, err
	}
	file, err := os.Create(name)
	if err != nil {
		return false, err
	}
	defer file.Close()
	return exists, nil
}
