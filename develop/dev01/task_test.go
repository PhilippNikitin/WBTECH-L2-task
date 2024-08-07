package main

import (
	"testing"
	"time"
)

// MockNTPClient представляет мок для NTP клиента
type MockNTPClient struct{}

// Time возвращает фиксированное время для тестирования
func (m MockNTPClient) Time(host string) (time.Time, error) {
	return time.Date(2023, time.March, 12, 15, 30, 0, 0, time.UTC), nil
}

func TestGetCurrentTime(t *testing.T) {
	// Используем мок для тестирования
	client := MockNTPClient{}

	expectedTime := time.Date(2023, time.March, 12, 15, 30, 0, 0, time.UTC)

	currentTime, err := getCurrentTime(client)
	if err != nil {
		t.Fatalf("getCurrentTime вернула ошибку: %v", err)
	}

	if !currentTime.Equal(expectedTime) {
		t.Errorf("Ожидалось время %v, но получено %v", expectedTime, currentTime)
	}
}
