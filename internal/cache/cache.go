package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	cache, err := NewFileCache("./cache.txt")
	if err != nil {
		log.Fatal(err)
	}
	cache.Save("abcdefg")
}

type Cacher interface {
	Get(int) (string, error)
	Save([]byte) error
}

func NewFileCache(path string) (fileCache *FileCache, err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}
	return &FileCache{filePath: path}, nil
}

type FileCache struct {
	filePath string
}

func (fc FileCache) Get() (entry string, err error) {
	cacheFile, err := os.Open(fc.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer cacheFile.Close()
	scanner := bufio.NewScanner(cacheFile)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return
}

func (fc FileCache) Save(item string) error {
	cacheFile, err := os.OpenFile(fc.filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer cacheFile.Close()
	tempFile, err := os.CreateTemp(".", "tempfile")
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()
	reader := bufio.NewReader(cacheFile)
	scanner := bufio.NewScanner(reader)
	writer := bufio.NewWriter(tempFile)
	_, err = io.WriteString(writer, item+"\n")

	for scanner.Scan() {
		line := scanner.Text()
		_, err = writer.WriteString(line + "\n")
		log.Print(line)
	}
	writer.Flush()
	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }
	// _ = os.Rename(tempFile.Name(), cacheFile.Name())
	return err
}
