// Package io provides functions for interfacing with input and output.
package io

import (
	"bytes"
	"errors"
	baseIO "io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	// Windows OS string.
	windows = "windows"
)

var (
	// OS-specific path separator.
	pathSeparator = string(os.PathSeparator)
)

type executeWithTimeoutResult struct {
	output []byte
	err    error
}

// init is a helper initializer to deal with invalid PATH variables deriving from mis-set GOROOT and GOPATH variables.
func init() {
	originalPath := os.Getenv("PATH")
	updatedPath := strings.TrimPrefix(originalPath, "undefined")
	if runtime.GOOS == windows {
		systemRoot := os.Getenv("SYSTEMROOT")
		if !strings.Contains(strings.ToLower(updatedPath), strings.ToLower(systemRoot)+"\\system32") {
			updatedPath += ";" + systemRoot + "\\system32"
		}
	}
	if updatedPath != originalPath {
		err := os.Setenv("PATH", updatedPath)
		if err != nil {
			log.Println("Error setting environment variable (PATH): " + err.Error())
		}
	}
}

// AddFileURI ensures the file:// prefix exists on a URI.
func AddFileURI(hostname string, uri string) string {
	if strings.HasPrefix(uri, "file://") {
		return uri
	}

	return "file://" + hostname + strings.Replace("/"+strings.Replace(uri, "\\", "/", -1), "//", "/", -1)
}

// CopyFile copies a file, overwriting if the destination exists.
func CopyFile(source string, destination string) error {
	in, err := os.Open(NormalizePathSeparators(source))
	if err != nil {
		return err
	}
	defer func() {
		err = in.Close()
		if err != nil {
			log.Println("Error closing input file (" + source + "): " + err.Error())
		}
	}()

	out, err := os.Create(NormalizePathSeparators(destination))
	if err != nil {
		return err
	}

	_, err = baseIO.Copy(out, in)
	if err != nil {
		closeErr := out.Close()
		if closeErr != nil {
			log.Println("Error closing output file (" + destination + "): " + closeErr.Error())
		}
		return err
	}
	return out.Close()
}

// DirectoryExists checks whether a directory exists.
func DirectoryExists(path string) bool {
	stats, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stats.IsDir()
}

// EnsureDirectory ensures a directory exists.
func EnsureDirectory(uri string) error {
	uri = NormalizePathSeparators(uri)
	uriParts := strings.Split(uri, pathSeparator)
	builtURI := uriParts[0]

	for i := 1; i < len(uriParts); i++ {
		builtURI += pathSeparator + uriParts[i]

		_, err := os.Stat(builtURI)
		if err != nil {
			err = os.Mkdir(builtURI, 0700)
			if err != nil {
				if !strings.HasPrefix(err.Error(), "Cannot create a file when that file already exists") {
					return err
				}
			}
		}
	}

	return nil
}

// ExecuteWithTimeout executes a command, killing it after too long.
func ExecuteWithTimeout(timeout time.Duration, name string, args ...string) ([]byte, error) {
	// Set default timeout of three minutes.
	if timeout == 0 {
		timeout = 3 * time.Minute
	}

	// Wait for completion or fail upon timing out.
	cmd := exec.Command(name, args...)
	closed := false
	closedMutex := new(sync.Mutex)
	closedChannel := make(chan executeWithTimeoutResult, 1)
	go func() {
		output := executeWithTimeoutResult{}
		closedMutex.Lock()
		if !closed {
			closedMutex.Unlock()
			// Map standard output and errors.
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			// Start command.
			err := cmd.Start()
			if err != nil {
				output.err = errors.New(err.Error() + " (" + stderr.String() + ")")
				closedChannel <- output
				return
			}

			// Wait for command to complete.
			err = cmd.Wait()
			closedMutex.Lock()
			if closed {
				closedMutex.Unlock()
				return
			}
			closedMutex.Unlock()
			if err != nil {
				output.err = errors.New(err.Error() + " (" + stderr.String() + ")")
			}

			// Read output and error.
			output.output = bytes.TrimRight(stdout.Bytes(), "\n")

			// Return result.
			closedMutex.Lock()
			closedChannel <- output
			closedMutex.Unlock()
		}
	}()
	select {
	case output := <-closedChannel:
		if output.err != nil {
			return nil, output.err
		}
		return output.output, nil
	case <-time.After(timeout):
		closedMutex.Lock()
		if !closed {
			closed = true
			closedMutex.Unlock()
			if cmd.Process != nil {
				err := cmd.Process.Kill()
				if err != nil {
					return nil, errors.New("error when killing process (" + name + " " + strings.Join(args, " ") + ") after timeout reached: " + err.Error())
				}
			}
			return nil, errors.New("process (" + name + " " + strings.Join(args, " ") + ") killed because timeout was reached")
		}
		closedMutex.Unlock()
		return nil, errors.New("no result returned from process (" + name + " " + strings.Join(args, " ") + ")")
	}
}

// FileExists checks whether a file exists.
func FileExists(path string) bool {
	stats, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !stats.IsDir()
}

// FileExistsInPath checks whether a file exists anywhere in the system path.
func FileExistsInPath(fileName string) (bool, error) {
	whichCommand := "which"
	if runtime.GOOS == windows {
		whichCommand = "where"
	}

	cmd := exec.Command(whichCommand, fileName)
	_, err := cmd.Output()
	if err != nil {
		if err.Error() == "exit status 1" {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// NormalizePathSeparators ensures that paths have correct slashes.
func NormalizePathSeparators(input string) string {
	switch runtime.GOOS {
	case windows:
		return strings.Replace(input, "/", "\\", -1)
	default:
		return input
	}
}

// ParseFileURI removes the file:// prefix and returns the hostname and file paths from a URI.
func ParseFileURI(uri string) (string, string) {
	hostname := ""
	if strings.HasPrefix(uri, "file://") {
		thirdSlash := strings.Index(uri[7:], "/")
		if thirdSlash > -1 {
			hostname = uri[7 : 7+thirdSlash]
			uri = uri[(thirdSlash + 8):]
			return hostname, NormalizePathSeparators(uri)
		}
	}

	return "", NormalizePathSeparators(uri)
}

// Remove removes a file or directory.
func Remove(uri string) error {
	uri = NormalizePathSeparators(uri)
	return os.Remove(uri)
}

// RemoveRecursive removes a directory and its children, recursively.
func RemoveRecursive(uri string) error {
	uri = NormalizePathSeparators(uri)
	if strings.Count(uri, pathSeparator) < 2 {
		return errors.New("too shallow of a path to delete")
	}

	fis, err := ioutil.ReadDir(uri)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if fi.IsDir() {
			err = RemoveRecursive(uri + pathSeparator + fi.Name())
		} else {
			err = os.Remove(uri + pathSeparator + fi.Name())
		}
		if err != nil {
			return err
		}
	}

	return os.Remove(uri)
}

// SanitizeDirectory ensures a directory path is safe.
func SanitizeDirectory(input string) string {
	return strings.Replace(strings.Replace(input, "/..", "", -1), "../", "", -1)
}
