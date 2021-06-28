package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"stuartdd.com/data"
)

var sourcePath string
var destPath string
var resultFile string = "results.txt"
var bashFile string = "bashrun.sh"
var bashDef string = ""
var verbose bool = false
var compressFileKey bool = false

var sourceMap = make(map[string]*data.FileData)
var typeMap string = ""

var filesInDestTotal = 0
var filesInDestChecked = 0

func main() {
	if len(os.Args) == 1 {
		exitWithHelp("No arguments provided", true)
	}
	readArgs(os.Args[1:])
	err := filepath.Walk(sourcePath, foundSourceFile)
	if err != nil {
		exitWithErr(err, 2)
	}
	err = filepath.Walk(destPath, foundDestinationFile)
	if err != nil {
		exitWithErr(err, 2)
	}
	writeResults()
	if bashDef != "" {
		writeBashFile()
	}
	if verbose {
		fmt.Printf(
			"Files from destination %d.\nFiles included         %d\n", filesInDestTotal, filesInDestChecked)
	}
}

func sortedKeys() []string {
	keys := make([]string, 0, len(sourceMap))
	for k := range sourceMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func writeBashFile() {
	count := 0
	f, err := os.Create(bashFile)
	if err != nil {
		exitWithErr(err, 2)
	}
	defer f.Close()
	for _, key := range sortedKeys() {
		count++
		value := sourceMap[key]
		line := value.BashString(bashDef, count)
		if line != "" {
			f.WriteString(value.BashString(bashDef, count))
			f.WriteString("\n")
		}
	}
	if verbose {
		fmt.Printf("BASHOUT: File %s\n", bashFile)
	}

}

func writeResults() {
	count := 0
	hitCount := 0
	size := len(sourceMap)
	f, err := os.Create(resultFile)
	if err != nil {
		exitWithErr(err, 2)
	}
	defer f.Close()
	f.WriteString("[\n")
	for _, key := range sortedKeys() {
		count++
		value := sourceMap[key]
		if value.MatchCount() > 0 {
			hitCount++
		}
		f.WriteString(value.String())
		if count >= size {
			f.WriteString("\n")
		} else {
			f.WriteString(",\n")
		}
	}
	f.WriteString("]\n")
	if verbose {
		fmt.Printf("HITS:    %d\n", hitCount)
		fmt.Printf("MISSES:  %d\n", count-hitCount)
		fmt.Printf("RESULTS: File %s\n", resultFile)
	}
}

func foundDestinationFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	filesInDestTotal++
	fileData, reason, exclude := deriveFileData(path, info)
	if exclude {
		if verbose {
			fmt.Printf("EXCLUDE: Key:%s Path:%s Reason:%s\n", fileData.Key(), path, reason)
		}
	} else {
		filesInDestChecked++
		fd, matchedName := sourceMap[fileData.Key()]
		if verbose {
			if matchedName {
				fmt.Printf("MATCH:   Key:%s Path:%s\n", fileData.Key(), path)
			} else {
				fmt.Printf("FILE:    Key:%s Path:%s\n", fileData.Key(), path)
			}
		}
		if matchedName {
			fd.CountMatch()
			fd.SetMatchedName(path)
			fd.SetMatchedSize(info.Size())
		}
	}
	return nil
}

func foundSourceFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	fileData, reason, exclude := deriveFileData(path, info)
	if exclude {
		if verbose {
			fmt.Printf("EXCLUDE: Source file '%s'. %s\n", path, reason)
		}
		return nil
	}
	fd, matchName := sourceMap[fileData.Key()]
	if matchName {
		fd.CountSource()
	} else {
		sourceMap[fileData.Key()] = fileData
		fileData.CountSource()
	}
	return nil
}

func deriveFileData(path string, info os.FileInfo) (*data.FileData, string, bool) {
	fileKey := info.Name()
	dotPos := strings.LastIndexByte(fileKey, '.')
	ext := ""
	if dotPos > 0 {
		ext = fileKey[dotPos:]
	}

	if compressFileKey {
		if ext == "" {
			dotPos = len(fileKey)
		}
		var buf bytes.Buffer
		for i := 0; i < dotPos; i++ {
			c := fileKey[i]
			if (c != '-') && (c != '.') && (c != ' ') && (c != '_') {
				buf.WriteByte(c)
			}
		}
		fileKey = buf.String() + ext
	}

	if typeMap == "" {
		return data.NewFileData(fileKey, path, info.Size()), "", false
	} else {
		if ext == "" {
			return data.NewFileData(fileKey, path, info.Size()), "No file extension", true
		} else {
			ind := strings.Index(typeMap, ext)
			if ind >= 0 {
				return data.NewFileData(fileKey, path, info.Size()), "", false
			} else {
				return data.NewFileData(fileKey, path, info.Size()), fmt.Sprintf("File extension '%s' not matched", ext), true
			}
		}
	}
}

func readArgs(args []string) {
	for _, a := range args {
		if strings.HasPrefix(strings.ToLower(a), "source:") {
			sourcePath = a[7:]
		} else {
			if strings.HasPrefix(strings.ToLower(a), "dest:") {
				destPath = a[5:]
			} else {
				if strings.HasPrefix(strings.ToLower(a), "ext:") {
					typeMap = typeMap + " ." + a[4:]
				} else {
					if strings.HasPrefix(strings.ToLower(a), "-v") {
						verbose = true
					} else {
						if strings.HasPrefix(strings.ToLower(a), "-c") {
							compressFileKey = true
						} else {
							if strings.HasPrefix(strings.ToLower(a), "out:") {
								resultFile = a[4:]
							} else {
								if strings.HasPrefix(strings.ToLower(a), "bash:") {
									bashDef = a[5:]
								} else {
									if strings.HasPrefix(strings.ToLower(a), "bashfile:") {
										bashFile = a[9:]
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if sourcePath == "" {
		exitWithHelp("source:<path> argument missing", true)
	}
	if destPath == "" {
		exitWithHelp("dest:<path> argument missing", true)
	}
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		exitWithHelp("source:"+sourcePath+" does not exist", false)
	}
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		exitWithHelp("dest:"+destPath+" does not exist", false)
	}
	if verbose {
		fmt.Println("SOURCE: ", sourcePath)
		fmt.Println("DEST:   ", destPath)
		fmt.Println("OUT:    ", resultFile)
		if bashDef != "" {
			fmt.Println("BASH:   ", bashDef)
			fmt.Println("BASHFIL:", bashFile)
		}
		if typeMap == "" {
			fmt.Println("EXT:     <any>")
		} else {
			fmt.Println("EXT:     ", typeMap)
		}
		if compressFileKey {
			fmt.Println("OPT:     Compressing File Key.")
		}
	}
}

func exitWithErr(err error, code int) {
	fmt.Println("ERROR :", err.Error())
	log.Println(err)
	os.Exit(code)
}

func exitWithHelp(message string, help bool) {
	fmt.Println("ERROR :", message)
	if help {
		fmt.Println("Args:  source:<path> Path to the file we are looking for")
		fmt.Println("Args:  dest:<path> Path to the where they should be")
	}
	os.Exit(1)
}
