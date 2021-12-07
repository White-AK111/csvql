package scanner

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type IScanner interface {
}

type Scanner struct {
	File    *os.File
	Headers []string
	Results [][]string
}

func NewScanner(filePath string) (*Scanner, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	scanner := &Scanner{
		File:    file,
		Headers: make([]string, 10),
		Results: make([][]string, 10),
	}

	return scanner, nil
}

func (a *Scanner) GetHeaders(comma string, comment string) error {
	reader := csv.NewReader(a.File)
	reader.Comma = []rune(comma)[0]
	reader.Comment = []rune(comment)[0]

	var headers []string
	lineCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		lineCount++

		if lineCount == 2 {
			break
		}

		headers = record
	}

	if lineCount < 2 {
		return errors.New("the numbers of lines in the file cannot be less then 2 (including headers)")
	}

	var cleanHeaders []string
	for _, header := range headers {
		if containsCount(headers, header) > 1 {
			return errors.New("duplicate header in source file")
		}
		header = strings.ToLower(header)
		cleanHeaders = append(cleanHeaders, header)
	}

	a.Headers = cleanHeaders
	return nil
}

func containsCount(s []string, e string) int {
	count := 0
	for _, a := range s {
		if a == e {
			count++
		}
	}
	return count
}
