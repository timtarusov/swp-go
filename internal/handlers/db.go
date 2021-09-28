package handlers

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	Db     *sqlx.DB
	Logger *log.Logger
}

func NewService(l *log.Logger) *Service {
	return &Service{
		Logger: l,
	}
}

func (s *Service) Info(format string, v ...interface{}) {
	s.Logger.SetPrefix("\u001b[34mINFO: ")
	s.Logger.Printf(format, v...)
}
func (s *Service) Debug(format string, v ...interface{}) {
	s.Logger.SetPrefix("\u001b[37mDEBUG: ")
	s.Logger.Printf(format, v...)
}
func (s *Service) Warning(format string, v ...interface{}) {
	s.Logger.SetPrefix("\u001b[33mWARNING: ")
	s.Logger.Printf(format, v...)
}
func (s *Service) Error(format string, v ...interface{}) {
	s.Logger.SetPrefix("\u001b[31mERROR: ")
	s.Logger.Printf(format, v...)
}
