package io

import (
	"io/ioutil"
	"os"
	"os/user"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	mutex sync.Mutex
)

// TestAddFileURI tests AddFileURI().
func TestAddFileURI(t *testing.T) {
	var hostname, uri string
	var expectedResult, result string
	hostname, _ = os.Hostname()

	// Test Windows path.
	uri = "c:\\grandparent\\parent\\child"
	expectedResult = "file://" + hostname + "/c:/grandparent/parent/child"
	result = AddFileURI(hostname, uri)
	assert.Equal(t, expectedResult, result, "Windows path.")

	// Test Linux path.
	uri = "/usr/home/tester"
	expectedResult = "file://" + hostname + "/usr/home/tester"
	result = AddFileURI(hostname, uri)
	assert.Equal(t, expectedResult, result, "Linux path.")
}

// TestCopyFile tests CopyFile().
func TestCopyFile(t *testing.T) {
	mutex.Lock()

	var err error
	var inputFileName, outputFileName string
	var inputBytes, outputBytes []byte

	// Test valid copy.
	currentDirectory, err := os.Getwd()
	assert.NoError(t, err, "Get current directory.")
	inputFileName = currentDirectory + "/io.go"
	err = EnsureDirectory(currentDirectory + "/test")
	assert.NoError(t, err, "Copy file.")
	outputFileName = currentDirectory + "/test/io.go.test"
	err = CopyFile(inputFileName, outputFileName)
	assert.NoError(t, err, "Copy file.")
	inputBytes, err = ioutil.ReadFile(inputFileName)
	assert.NoError(t, err, "Read input file.")
	outputBytes, err = ioutil.ReadFile(outputFileName)
	assert.NoError(t, err, "Read output file.")
	assert.Equal(t, inputBytes, outputBytes, "Copy file.")
	err = Remove(outputFileName)
	assert.NoError(t, err, "Remove file.")

	// Test invalid source.
	inputFileName = "io.go.missing"
	outputFileName = "io.go.test"
	err = CopyFile(inputFileName, outputFileName)
	assert.Error(t, err, "Invalid source.")

	// Test invalid destination.
	inputFileName = "io.go"
	outputFileName = "//zzz/io.go.test"
	err = CopyFile(inputFileName, outputFileName)
	assert.Error(t, err, "Invalid destination.")

	mutex.Unlock()
}

// TestDirectoryExists tests DirectoryExists().
func TestDirectoryExists(t *testing.T) {
	// Test valid directory.
	currentDirectory, err := os.Getwd()
	assert.NoError(t, err, "Get current directory.")
	assert.True(t, DirectoryExists(currentDirectory), "Valid directory.")

	// Test invalid directory.
	assert.False(t, DirectoryExists(currentDirectory+"/test/nonexistent/path"), "Invalid directory.")
}

// TestEnsureDirectory tests EnsureDirectory().
func TestEnsureDirectory(t *testing.T) {
	mutex.Lock()

	var err error
	var base string

	switch runtime.GOOS {
	case "windows":
		base, err = os.Getwd()
		assert.NoError(t, err)
	default:
		var usr *user.User
		usr, err = user.Current()
		if err != nil {
			t.Error("utils/user.Current: error: ", err.Error())
		}
		base = usr.HomeDir + "/"
	}
	err = EnsureDirectory(base + "/test/child/grandchild")
	assert.NoError(t, err, "Ensure directory.")
	assert.True(t, DirectoryExists(base+"/test/child/grandchild"), "Valid directory.")
	err = RemoveRecursive(base + "/test")
	assert.NoError(t, err, "Remove directory.")

	mutex.Unlock()
}

// TestExecuteWithTimeout tests ExecuteWithTimeout().
func TestExecuteWithTimeout(t *testing.T) {
	var timeout time.Duration
	var result, expectedResult []byte
	var err error

	// Run successful command..
	timeout = 15 * time.Second
	result, err = ExecuteWithTimeout(timeout, "echo", "success")
	expectedResult = []byte("success")
	assert.NoError(t, err, "Execute with timeout.")
	assert.Equal(t, expectedResult, result)

	// Run command with impossible timeout.
	timeout = 1 * time.Nanosecond
	result, err = ExecuteWithTimeout(timeout, "echo", "too short")
	expectedResult = nil
	assert.Equal(t, "process (echo too short) killed because timeout was reached", err.Error())
	assert.Equal(t, expectedResult, result)

	// Run command with timeout that will fail.
	timeout = 3 * time.Second
	result, err = ExecuteWithTimeout(timeout, "sleep", "5")
	assert.Equal(t, "process (sleep 5) killed because timeout was reached", err.Error())
	assert.Equal(t, expectedResult, result)
}

// TestFileExists tests FileExists().
func TestFileExists(t *testing.T) {
	// Test valid file.
	assert.True(t, FileExists("io.go"), "File exists.")

	// Test invalid file.
	assert.False(t, FileExists("nonexistent.file.dat"), "File does not exist.")
}

// TestFileExistsInPath tests FileExistsInPath().
func TestFileExistsInPath(t *testing.T) {
	var exists bool
	var err error

	// Run different command based on OS.
	switch runtime.GOOS {
	case "windows":
		// Test valid file.
		exists, err = FileExistsInPath("cmd.exe")
		assert.NoError(t, err, "cmd.exe")
		assert.True(t, exists, "cmd.exe")
	default:
		// Test valid file.
		exists, err = FileExistsInPath("echo")
		assert.NoError(t, err, "echo")
		assert.True(t, exists, "echo")
	}

	// Test invalid file.
	exists, err = FileExistsInPath("nonexistent.file.dat")
	assert.NoError(t, err, "Invalid file.")
	assert.False(t, exists, "Invalid file.")
}

// TestNormalizePathSeparators tests NormalizePathSeparators().
func TestNormalizePathSeparators(t *testing.T) {
	var input, result, expectedResult string

	// Run different tests based on OS.
	switch runtime.GOOS {
	case "windows":
		// Test for valid replacement.
		input = "c:/domain/kingdom/phylum/class"
		expectedResult = "c:\\domain\\kingdom\\phylum\\class"
		result = NormalizePathSeparators(input)
		assert.Equal(t, expectedResult, result, "Normalized path separators.")

		// Test for passthrough.
		input = "c:\\order\\family\\genus\\species"
		expectedResult = "c:\\order\\family\\genus\\species"
		result = NormalizePathSeparators(input)
		assert.Equal(t, expectedResult, result, "Normalized path separators.")
	default:
		// Test for passthrough.
		input = "order/family/genus/species"
		expectedResult = "order/family/genus/species"
		result = NormalizePathSeparators(input)
		assert.Equal(t, expectedResult, result, "Normalized path separators.")
	}
}

// TestParseFileURI tests ParseFileURI().
func TestParseFileURI(t *testing.T) {
	var input string
	var hostnameResult, expectedHostnameResult string
	var pathResult, expectedPathResult string

	// Test valid local URI.
	input = "file:///C:/Windows/System32/calc.exe"
	hostnameResult, pathResult = ParseFileURI(input)
	expectedHostnameResult = ""
	switch runtime.GOOS {
	case "windows":
		expectedPathResult = "C:\\Windows\\System32\\calc.exe"
	default:
		expectedPathResult = "C:/Windows/System32/calc.exe"
	}
	assert.Equal(t, expectedHostnameResult, hostnameResult, "Parsed file URI.")
	assert.Equal(t, expectedPathResult, pathResult, "Parsed file URI.")

	// Test valid remote URI.
	input = "file://machinename/C:/Windows/System32/calc.exe"
	hostnameResult, pathResult = ParseFileURI(input)
	expectedHostnameResult = "machinename"
	switch runtime.GOOS {
	case "windows":
		expectedPathResult = "C:\\Windows\\System32\\calc.exe"
	default:
		expectedPathResult = "C:/Windows/System32/calc.exe"
	}
	assert.Equal(t, expectedHostnameResult, hostnameResult, "Parsed file URI.")
	assert.Equal(t, expectedPathResult, pathResult, "Parsed file URI.")

	// Test invalid.
	input = "file:Washington/Adams/Jefferson"
	hostnameResult, pathResult = ParseFileURI(input)
	expectedHostnameResult = ""
	switch runtime.GOOS {
	case "windows":
		expectedPathResult = "file:Washington\\Adams\\Jefferson"
	default:
		expectedPathResult = "file:Washington/Adams/Jefferson"
	}
	assert.Equal(t, expectedHostnameResult, hostnameResult, "Parsed file URI.")
	assert.Equal(t, expectedPathResult, pathResult, "Parsed file URI.")
}

// TestRemove tests Remove()
func TestRemove(t *testing.T) {
	mutex.Lock()

	base, err := os.Getwd()
	assert.NoError(t, err)

	err = EnsureDirectory(base + "/test/to/be/removed")
	assert.NoError(t, err)

	err = Remove(base + "/test/to/be/removed")
	assert.NoError(t, err)

	err = Remove(base + "/test/to")
	assert.Error(t, err)

	err = Remove(base + "/test/to/be")
	assert.NoError(t, err)

	err = Remove(base + "/test/to")
	assert.NoError(t, err)

	err = Remove(base + "/test")
	assert.NoError(t, err)

	mutex.Unlock()
}

// TestRemoveRecursive tests RemoveRecursive()
func TestRemoveRecursive(t *testing.T) {
	mutex.Lock()

	base, err := os.Getwd()
	assert.NoError(t, err)

	err = EnsureDirectory(base + "/test/recursive/to/be/removed")
	assert.NoError(t, err)

	err = RemoveRecursive(base + "/test")
	assert.NoError(t, err)

	mutex.Unlock()
}

// TestSanitizeURI tests SanitizeURI().
func TestSanitizeURI(t *testing.T) {
	var input, result, expectedResult string

	// Test valid input.
	input = "test.txt"
	expectedResult = "test.txt"
	result = SanitizeDirectory(input)
	assert.Equal(t, expectedResult, result, "Valid input.")

	// Test invalid input.
	input = "../test.txt"
	expectedResult = "test.txt"
	result = SanitizeDirectory(input)
	assert.Equal(t, expectedResult, result, "Invalid input.")
}
