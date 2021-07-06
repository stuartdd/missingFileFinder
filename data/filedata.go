package data

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
	NONE uint16 = 0b0000000000000000
	NAME uint16 = 0b0000000000000001
	SIZE uint16 = 0b0000000000000010
	BYTE uint16 = 0b0000000000000100
)

type FileData struct {
	key           string
	sourceName    string
	destNames     string
	size          int64
	match         uint16
	matchCount    uint32
	sourceCount   uint32
	filePrefix    []uint8
	filePrefixLen int16
}

func NewFileData(pKey string, pSourceName string, pSize int64) *FileData {
	return &FileData{key: pKey, sourceName: pSourceName, destNames: "", size: pSize, match: NONE, matchCount: 0, sourceCount: 0, filePrefixLen: 0}
}

func (r *FileData) GetKey() string {
	return r.key
}

func (r *FileData) GetSourceName() string {
	return r.sourceName
}

func (r *FileData) GetMatchCount() uint32 {
	return r.matchCount
}

func (r *FileData) GetSize() int64 {
	return r.size
}

func (r *FileData) GetSourceCount() uint32 {
	return r.sourceCount
}

func (r *FileData) SetMatchCount(count uint32) {
	r.matchCount = count
}

func (r *FileData) SetCountSource(count uint32) {
	r.sourceCount = count
}

func (r *FileData) SetFilePrefix(b []uint8, bc int16) {
	r.filePrefix = b
	r.filePrefixLen = bc
}

func (r *FileData) GetFilePrefix() []uint8 {
	return r.filePrefix
}

func (r *FileData) GetFilePrefixLen() int16 {
	return r.filePrefixLen
}

func (r *FileData) IncSourceCount() {
	r.sourceCount++
}

func (r *FileData) AddDestName(name string) {
	if r.destNames == "" {
		r.destNames = name
	} else {
		r.destNames = r.destNames + " | " + name
	}
}

func (r *FileData) SetMatchedOnName() {
	r.match = r.match | NAME
}

func (r *FileData) SetMatchedOnSize() {
	r.match = r.match | SIZE
}

func (r *FileData) SetMatchedOnSizeBytes() {
	r.match = r.match | BYTE
}

func (r *FileData) Match() string {
	desc := ""
	if (r.match & NAME) == NAME {
		desc = desc + "name"
	}
	if (r.match & SIZE) == SIZE {
		desc = desc + "+size+"
	}
	if (r.match & BYTE) == BYTE {
		desc = desc + "-byte-"
	}
	if desc == "" {
		return "no-match"
	}
	return desc
}

func (r *FileData) BashString(def string, line int) string {
	if strings.HasPrefix(def, r.Match()) || strings.HasPrefix(def, "all") {
		resp := ""
		if strings.HasPrefix(def, "all") {
			resp = def[3:]
		} else {
			resp = def[len(r.Match()):]
		}
		resp = strings.Replace(resp, "$line", strconv.Itoa(line), -1)
		resp = strings.Replace(resp, "$key", r.key, -1)
		resp = strings.Replace(resp, "$source", r.sourceName, -1)
		resp = strings.Replace(resp, "$match", r.Match(), -1)
		resp = strings.Replace(resp, "$nl", "\n", -1)
		resp = strings.Replace(resp, "$size", strconv.Itoa(int(r.size)), -1)
		resp = strings.Trim(resp, " ")
		return resp
	}
	return ""
}

func (r *FileData) String() string {
	hex.EncodeToString(r.filePrefix)
	return fmt.Sprintf("{\"key\":\"%s\", \"mc\":%d, \"sc\":%d, \"source\":\"%s\", \"dest\":\"%s\", \"si\":%d, \"match\":\"%s\"}",
		r.key, r.matchCount, r.sourceCount, r.sourceName, r.destNames, r.size, r.Match())
}
