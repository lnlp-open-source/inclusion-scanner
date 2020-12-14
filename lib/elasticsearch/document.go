package elasticsearch

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Document struct {
	Id                    string    `json:"-"`
	Filepath              string    `json:"filepath"`
	NonInclusiveTermsUsed []string  `json:"terms_abused"`
	Timestamp             time.Time `json:"timestamp"`
}

func NewDocument(path string, nonInclusiveTermsUsed []string) *Document {
	return &Document{
		Id:                    uniqueIdFromFilePath(path),
		Filepath:              path,
		NonInclusiveTermsUsed: nonInclusiveTermsUsed,
		Timestamp:             time.Now(),
	}
}

func (document *Document) GetPayload() ([]byte, error) {
	result, err := json.Marshal(document)
	if err != nil {
		return []byte{}, err
	}
	return result, nil
}

func uniqueIdFromFilePath(key string) string {
	hash := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", hash[:])
}
