package crawler

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"path/filepath"
	"strings"
)

// getCategory take a filepath and extract category based on path
func getCategory(r string, p string) (string, error) {
	var (
		err      error
		basePath string //Base path to strip from root
		relPath  string //Relative path from doc root
		category string //Calculated category
	)

	// Remove last folder from path to use it as category root
	basePath = filepath.Dir(filepath.ToSlash(r))

	// extract relative path from doc root
	relPath, err = filepath.Rel(basePath, filepath.ToSlash(p))
	if err != nil {
		return "", errors.New(fmt.Sprintf("crawler/helper/getCategory returned error while processing relative filepath: %v\n", err))
	}

	// extract category from relative path, stripping filename
	category = filepath.Dir(relPath)

	return category, nil
}

// getDisplayName take a path and calculate human-readable file name, stripping useless char
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
		// file          *os.File  //In memory file
		hashInterface hash.Hash //Hash interface to write to
		hashInByte    []byte    //Calculated hash as []byte
		checksum      string    //SHA-1 checksum
		err           error
	)

	// -------------------------
	// Old method using whole file content , deprecated vs filename only calculation
	// -------------------------

	//file, err = os.Open(p)
	//if err != nil {
	//	return "", err
	//}
	//
	//defer file.Close()
	//
	//hashInterface = sha1.New()
	//
	//if _, err = io.Copy(hashInterface, p); err != nil {
	//	return "", err
	//}

	// Calculate SAH-1
	//hashInByte = hashInterface.Sum(nil)

	// -------------------------
	// New method using only filename for calculation
	// -------------------------

	hashInterface = sha1.New()

	_, err = io.WriteString(hashInterface, p)
	if err != nil {
		return "", err
	}

	hashInByte = hashInterface.Sum(nil)

	// Convert [] byte in string
	checksum = hex.EncodeToString(hashInByte)

	return checksum, nil
}

// getFirstNChar return first n char from given string (unicode safe)
func getFirstNChar(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

// TODO: Implement
// parseVehicle take an encoded vehicle string in the form [ente]-[stazionamento].[lotto] and return 3 separated value.
//
// Actually [lotto] always return ""
func parseVehicle(v string) (string, string, string) {
	var (
		working       []string
		ente          []string
		stazionamento []string
	)

	// Check if multiple ente ad extract them
	working = strings.Split(v, "/")
	// For each ente-stazionamento pair extract single value and append to relative slice
	for _, s := range working {
		splitted := strings.Split(s, "-")
		// FIXME: Check if malformed entry without "-" separator, resulting in array of length 1
		ente = append(ente, splitted[0])
		stazionamento = append(stazionamento, splitted[1])
	}
	return strings.Join(ente, " - "), strings.Join(stazionamento, " - "), ""
}
