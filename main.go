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
var sizeMap = make(map[int64]*data.SizeData)
var typeMap string = ""
var fileStartLen int64 = 100
var fileStartSeek int64 = 2000
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

func readFileStart(fileName string, size int64) ([]byte, int16) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	defer f.Close()
	if size > fileStartSeek {
		f.Seek(fileStartSeek-fileStartLen, 0)
	}
	buf := make([]byte, fileStartLen)
	bc, err1 := f.Read(buf)
	if err1 != nil {
		log.Fatalf("unable to read file: %v", err1)
	}

	//	fmt.Print(string(buf))
	return buf, int16(bc)
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
		if value.GetMatchCount() > 0 {
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
	destFileData, reason, exclude := deriveFileData(path, info)
	if exclude {
		if verbose {
			fmt.Printf("EXCLUDE: Key:%s Path:%s Reason:%s\n", destFileData.GetKey(), path, reason)
		}
	} else {
		filesInDestChecked++
		matchFileData, foundInMap := sourceMap[destFileData.GetKey()]
		if foundInMap {
			currentMC := matchFileData.GetMatchCount()
			matchFileData.SetMatchCount(currentMC + 1)
			matchFileData.SetMatchedOnName()
			matchFileData.AddDestName(path)
			if matchFileData.GetSize() == info.Size() {
				if verbose {
					fmt.Printf("NAME+SI: Key:%s Path:%s\n", matchFileData.GetKey(), path)
				}
				matchFileData.SetMatchedOnSize()
			} else {
				if verbose {
					fmt.Printf("NAME:    Key:%s Path:%s\n", matchFileData.GetKey(), path)
				}
			}
		} else {
			matchSizeData, matchedSize := sizeMap[info.Size()]
			if matchedSize {
				matchFileData, foundInMap := sourceMap[matchSizeData.FileKey()]
				if foundInMap {
					bytes, len := readFileStart(path, matchFileData.GetSize())
					if compareButes(bytes, len, matchFileData) {
						currentMC := matchFileData.GetMatchCount()
						if verbose {
							fmt.Printf("SIZE+BI: Key:%s Path:%s\n", matchSizeData.FileKey(), path)
						}
						matchFileData.SetMatchCount(currentMC + 1)
						matchFileData.SetMatchedOnSizeBytes()
						matchFileData.AddDestName(path)
					} else {
						fmt.Printf("SIZE-NB: Key:%s Path:%s\n", matchSizeData.FileKey(), path)
					}
				} else {
					fmt.Printf("ERROR:   Key:%s from SizeData not found in SourceMap for source:%s\n", matchSizeData.FileKey(), path)
				}
			} else {
				fmt.Printf("SIZE-NM: Size:%d Path:%s\n", info.Size(), path)
			}
		}
	}
	return nil
}

func compareButes(buff []uint8, len int16, dataFile *data.FileData) bool {
	if len != dataFile.GetFilePrefixLen() {
		return false
	}
	var i int16 = 0
	for i = 0; i < len; i++ {
		if buff[i] != dataFile.GetFilePrefix()[i] {
			return false
		}
	}
	return true
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
	fd, matchName := sourceMap[fileData.GetKey()]
	if matchName {
		fd.IncSourceCount()
		fmt.Printf("INC-DUP: Key:%s Size:%d Dupe:%d Source file '%s'\n", fd.GetKey(), info.Size(), fd.GetSourceCount(), path)
	} else {
		sourceMap[fileData.GetKey()] = fileData
		fileData.IncSourceCount()
		b, len := readFileStart(path, fileData.GetSize())
		fileData.SetFilePrefix(b, len)
		fmt.Printf("INCLUDE: Key:%s Size:%d Source file '%s'\n", fileData.GetKey(), info.Size(), path)
	}
	_, matchSd := sizeMap[info.Size()]
	if !matchSd {
		sizeMap[info.Size()] = data.NewSizeData(uint64(info.Size()), fileData.GetKey())
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
