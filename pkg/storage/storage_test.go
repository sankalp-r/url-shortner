package storage

import (
	"testing"
)

func TestStore(t *testing.T) {
	testStore := NewStore()
	testURL := "https://test.com"
	shortCode, err := testStore.Create(testURL)
	if err != nil {
		t.Error("short URL creation failed")
	}

	url, err := testStore.Get(shortCode)
	if err != nil {
		t.Error("getting URL failed")
	}

	if url != testURL {
		t.Error("URL not equal")
	}

}
