package data

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	NONE uint16 = 0b0000000000000000
	NAME uint16 = 0b0000000000000001
	SIZE uint16 = 0b0000000000000010
)

type FileData struct {
	key         string
	sourceName  string
	destNames   string
	size        int64
	match       uint16
	matchCount  uint32
	sourceCount uint32
}

func NewFileData(pKey string, pSourceName string, pSize int64) *FileData {
	return &FileData{key: pKey, sourceName: pSourceName, destNames: "", size: pSize, match: NONE, matchCount: 0, sourceCount: 0}
}

func (r *FileData) Key() string {
	return r.key
}

func (r *FileData) SourceName() string {
	return r.sourceName
}

func (r *FileData) MatchCount() uint32 {
	return r.matchCount
}

func (r *FileData) SourceCount() uint32 {
	return r.sourceCount
}

func (r *FileData) CountMatch() {
	r.matchCount++
}

func (r *FileData) CountSource() {
	r.sourceCount++
}

func (r *FileData) SetMatchedName(name string) {
	if r.destNames == "" {
		r.destNames = name
	} else {
		r.destNames = r.destNames + " | " + name
	}
	r.match = r.match | NAME
}

func (r *FileData) SetMatchedSize(pSize int64) {
	if r.size == pSize {
		r.match = r.match | SIZE
	}
}

func (r *FileData) Match() string {
	desc := ""
	if (r.match & NAME) == NAME {
		desc = "name"
	}
	if (r.match & SIZE) == SIZE {
		desc = desc + "+size"
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
	return fmt.Sprintf("{\"key\":\"%s\", \"mc\":%d, \"sc\":%d, \"source\":\"%s\", \"dest\":\"%s\", \"si\":%d, \"match\":\"%s\"}",
		r.key, r.matchCount, r.sourceCount, r.sourceName, r.destNames, r.size, r.Match())
}
