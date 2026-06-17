package model

import (
	"errors"
	"time"
)

type URL struct {
	Hash      string    `json:"hash,omitempty"`
	URL       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Expiry    time.Time `json:"expiry"`
}

var (
	memoryMap []URL
	hashMap   map[string]bool
)

func initStorage() {
	if hashMap == nil {
		hashMap = make(map[string]bool)
	}
}

func (this URL) GetModel(hash string) (URL, error) {
	initStorage()

	if _, exists := hashMap[hash]; !exists {
		return URL{}, errors.New("Value Doesnt Exists")
	}

	for _, it := range memoryMap {
		if it.Hash == hash {
			return it, nil
		}
	}
	return URL{}, errors.New("Value Doesnt Exists")
}

func (this URL) AddURL() (string, error) {
	initStorage()

	if this.Hash == "" {
		return "", errors.New("Hash is required")
	}

	if _, exists := hashMap[this.Hash]; exists {
		return "", errors.New("Value Exists")
	}

	if this.CreatedAt.IsZero() {
		this.CreatedAt = time.Now()

		this.Expiry = this.CreatedAt.Add(10 * time.Second)
	}

	memoryMap = append(memoryMap, this)
	hashMap[this.Hash] = true

	return this.Hash, nil
}
