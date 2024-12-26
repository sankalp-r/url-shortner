package storage

import (
	"errors"
	"math/rand"
	"net/url"
	"sync"

	"github.com/sankalp-r/url-shortner/internal/models"
)

// base-32 characters set
const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

// Store interface defines methods for creating and retrieving URL mappings
type Store interface {
	Create(string) (string, error)
	Get(string) (string, error)
}

// NewStore creates a new instance of InMemoryStore
func NewStore() Store {
	return &InMemoryStore{
		urlMap: make(map[string]models.URLMapping),
		lock:   sync.RWMutex{},
	}
}

// InMemoryStore is an in-memory implementation of the Store interface
type InMemoryStore struct {
	urlMap map[string]models.URLMapping
	lock   sync.RWMutex
}

// Create generates a short code for the given URL and stores the mapping
func (ims *InMemoryStore) Create(URL string) (string, error) {
	// Validate the URL
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return "", err
	}

	ims.lock.Lock()
	defer ims.lock.Unlock()
	shortCode := generateCode(7)
	ims.urlMap[shortCode] = models.URLMapping{
		OriginalURL:  URL,
		ShortURLCode: shortCode,
	}
	return shortCode, nil
}

// Get retrieves the original URL for the given short code
func (ims *InMemoryStore) Get(shortURLCode string) (string, error) {
	ims.lock.RLock()
	defer ims.lock.RUnlock()
	if val, exist := ims.urlMap[shortURLCode]; exist {
		return val.OriginalURL, nil
	}
	return "", errors.New("URL not found")
}

// generateCode generates a random code of the specified length using the charset
func generateCode(length int) string {
	output := make([]byte, 0, length)

	for i := 0; i < length; i++ {
		index := rand.Intn(32)
		output = append(output, charset[index])
	}

	return string(output)
}
