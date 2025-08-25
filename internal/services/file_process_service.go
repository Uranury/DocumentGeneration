package services

import "log/slog"

type FileProcessService struct {
	logger *slog.Logger
}

func NewFileProcessService(logger *slog.Logger) *FileProcessService {
	return &FileProcessService{logger: logger}
}
