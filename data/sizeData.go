package data

import (
	"fmt"
)

type SizeData struct {
	key     uint64
	fileKey string
}

func NewSizeData(pKey uint64, pFileKey string) *SizeData {
	return &SizeData{key: pKey, fileKey: pFileKey}
}

func (r *SizeData) Key() uint64 {
	return r.key
}

func (r *SizeData) FileKey() string {
	return r.fileKey
}

func (r *SizeData) String() string {
	return fmt.Sprintf("{\"key\":\"%d\", \"fk\":%s\"}", r.key, r.fileKey)
}
