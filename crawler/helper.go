package crawler

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// getCategory take a filepath and extract category based on path
func getCategory(r string, p string) (string, error) {
	var (
		err error
		//docRoot  string //Documents disk path
		relPath  string //Relative path from doc root
		category string //Calculated category
	)

	//Read doc root from env
	//docRoot, err = utils.ReadEnv("DOC_ROOT")
	//if err != nil {
	//	log.Fatalf("Error retrieving documents root from env")
	//}

	// extract relative path from doc root
	relPath, err = filepath.Rel(r, p)
	if err != nil {
		return "", errors.New(fmt.Sprintf("crawler/helper/getCategory returned error while processing relative filepath: %v\n", err))
	}

	// extract category from relative path, stripping filename
	category = filepath.Dir(relPath)

	return category, nil
}

// getDisplayName take a path and calculate human readable file name, stripping useless char
func getDisplayNameFromPath(p string) string {
	var (
		fileName                 string //Filename from path
		extension                string //Filename extension
		fileNameWithoutExtension string //Filename without ext
		displayName              string //Calculated display name
	)

	// extract filename from path, trimming lead and trail spaced
	fileName = strings.TrimSpace(filepath.Base(p))

	// extract extension from filename
	extension = filepath.Ext(p)

	// trim extension from filename
	fileNameWithoutExtension = strings.TrimSuffix(fileName, extension)

	displayName = strings.ReplaceAll(fileNameWithoutExtension, "_", " ")
	return displayName
}

// getSha1 take a valid file path and calculate SHA-1 checksum
func getSha1(p string) (string, error) {
	var (
		err           error
		file          *os.File  //In memory file
		hashInterface hash.Hash //Hash interface to write to
		hashInByte    []byte    //Calculated hash as []byte
		checksum      string    //SHA-1 checksum
	)

	file, err = os.Open(p)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hashInterface = sha1.New()

	if _, err = io.Copy(hashInterface, file); err != nil {
		return "", err
	}

	// Calculate SAH-1
	hashInByte = hashInterface.Sum(nil)

	// Convert [] byte in string
	checksum = hex.EncodeToString(hashInByte)

	return checksum, nil
}
