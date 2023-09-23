package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

type CacheBufSizeError struct{}

func (err *CacheBufSizeError) Error() string {
	return "bufSize value must be less than 10"
}

func main() {
	cache, err := NewFileCache("./cache.txt", 9)
	if err != nil {
		log.Fatal(err)
	}
	cache.Save("Item1")
}

type Cacher interface {
	Get(int) (string, error)
	Save([]byte) error
}

func NewFileCache(filePath string, bufSize int) (fileCache *FileCache, err error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err = os.Create(filePath)
		if err != nil {
			return nil, err
		}
	}
	if bufSize > 10 {
		return nil, &CacheBufSizeError{}
	}
	return &FileCache{filePath: filePath, bufSize: bufSize}, nil
}

type FileCache struct {
	filePath string
	bufSize  int
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
	scanner := bufio.NewScanner(bufio.NewReader(cacheFile))
	writer := bufio.NewWriter(tempFile)
	_, err = io.WriteString(writer, item+"\n")

	for i := 0; i < fc.bufSize; i++ {
		scanner.Scan()
		line := scanner.Text()
		_, err = writer.WriteString(line + "\n")
		log.Print(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	writer.Flush()
	_ = os.Rename(tempFile.Name(), cacheFile.Name())
	return err
}
