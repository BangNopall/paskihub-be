package helpers

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/pkg/log"
)

func Contains(search string, words []string) bool {
	for _, word := range words {
		if search == word {
			return true
		}
	}

	return false
}

func ReadFile(filepath string, separator string) ([]string, error) {
	file, err := os.Open(filepath)

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[HELPERS][ReadFiles] failed to open file")
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	results := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()

		results = append(results, line)
	}

	if err := scanner.Err(); err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[HELPERS][ReadFiles] failed to scan file")
		return nil, err
	}

	return results, nil
}

func CheckRowsAffected(rows int64) error {
	if rows == 0 {
		return domain.ErrNotFound
	}

	if rows > 1 {
		return fmt.Errorf("weird behaviour. rows affected : %v", rows)
	}

	return nil
}

func GenerateRandomString(lenght int) string {
	alphaNumRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	randomRune := make([]rune, lenght)

	for i := 0; i < lenght; i++ {
		randomRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}

	return string(randomRune)
}
