package repository

import (
	"encoding/json"
	"github.com/Guram-Gurych/shortenerURL.git/internal/logger"
	"github.com/Guram-Gurych/shortenerURL.git/internal/model"
	"go.uber.org/zap"
	"io"
	"os"
	"strconv"
	"sync"
)

type FileRepository struct {
	urls       map[string]string
	mu         sync.RWMutex
	filePath   string
	descriptor *os.File
	encoder    *json.Encoder
	uuidCount  int
}

func NewFileRepository(filePath string) (*FileRepository, error) {
	fileRepository := &FileRepository{
		urls:      make(map[string]string),
		filePath:  filePath,
		uuidCount: 0,
	}

	if filePath == "" {
		return fileRepository, nil
	}

	err := fileRepository.loadFromFile()
	if err != nil {
		logger.Log.Error("Не удалось загрузить URL из файла", zap.Error(err))
		return nil, err
	}

	file, err := os.OpenFile(fileRepository.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	fileRepository.descriptor = file
	fileRepository.encoder = json.NewEncoder(file)

	return fileRepository, nil
}

func (rep *FileRepository) loadFromFile() error {
	file, err := os.OpenFile(rep.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Error("Не удалось открыть файл для чтения", zap.String("path", rep.filePath), zap.Error(err))
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var record model.URLModel

		if err = decoder.Decode(&record); err != nil {
			if err == io.EOF {
				break
			}

			logger.Log.Error("Не удалось декодировать запись из файла", zap.String("path", rep.filePath), zap.Error(err))
			return err
		}

		rep.urls[record.ShortURL] = record.OriginalURL
		if uuid, err := strconv.Atoi(record.UUID); err == nil {
			if uuid > rep.uuidCount {
				rep.uuidCount = uuid
			}
		}
	}

	return nil
}

func (rep *FileRepository) Save(id, url string) error {
	rep.mu.Lock()
	defer rep.mu.Unlock()

	if _, ok := rep.urls[id]; ok {
		return ErrorAlreadyExists
	}

	rep.urls[id] = url

	if rep.encoder == nil {
		return nil
	}

	rep.uuidCount++
	record := model.URLModel{
		UUID:        strconv.Itoa(rep.uuidCount),
		ShortURL:    id,
		OriginalURL: url,
	}

	if err := rep.encoder.Encode(&record); err != nil {
		logger.Log.Error("Не удалось записать URL в файл", zap.String("path", rep.filePath), zap.Error(err))
		return err
	}

	return nil
}

func (rep *FileRepository) Get(id string) (string, error) {
	rep.mu.RLock()
	defer rep.mu.RUnlock()

	val, ok := rep.urls[id]
	if !ok {
		return "", ErrorNotFound
	}

	return val, nil
}

func (rep *FileRepository) Close() error {
	if rep.descriptor != nil {
		return rep.descriptor.Close()
	}
	return nil
}
