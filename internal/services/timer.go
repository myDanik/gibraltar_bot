package services

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type TimerService struct {
	Cache    *Cache
	filename string
}

func NewTimerService(cache *Cache, filename string) *TimerService {
	service := &TimerService{
		Cache:    cache,
		filename: filename,
	}
	service.FillCacheFromFile()
	return service
}

func (s *TimerService) FillCacheFromFile() error {
	file, err := os.Open(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" {
			id, err := strconv.ParseInt(line, 10, 64)
			if err != nil {
				return err
			}
			s.Cache.Set(id, 1)
		}
	}
	return nil
}

func (s *TimerService) AddNewChatToTimer(chatID int64) (success bool) {
	v, ok := s.Cache.Get(chatID)
	if ok && v == 1 {
		return false
	}
	data := strconv.FormatInt(chatID, 10)

	err := appendToFile(s.filename, data)
	if err != nil {
		log.Println("append to file: %w", err)
		return false
	}
	s.Cache.Set(chatID, 1)
	return true
}

func appendToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data + "\n")
	return err
}
