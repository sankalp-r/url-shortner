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

type Store interface {
	Create(string) (string, error)
	Get(string) (string, error)
}

func NewStore() Store {
	return &InMemoryStore{
		urlMap: make(map[string]models.URLMapping),
		lock:   sync.RWMutex{},
	}
}

type InMemoryStore struct {
	urlMap map[string]models.URLMapping
	lock   sync.RWMutex
}

func (ims *InMemoryStore) Create(URL string) (string, error) {
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

func (ims *InMemoryStore) Get(shortURLCode string) (string, error) {
	ims.lock.RLock()
	defer ims.lock.RUnlock()
	if val, exist := ims.urlMap[shortURLCode]; exist {
		return val.OriginalURL, nil
	}
	return "", errors.New("URL not found")
}

func generateCode(length int) string {
	output := make([]byte, 0, length)

	for i := 0; i < length; i++ {
		index := rand.Intn(32)
		output = append(output, charset[index])
	}

	return string(output)
}
