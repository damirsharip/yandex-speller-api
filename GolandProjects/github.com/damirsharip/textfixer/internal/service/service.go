package service

import "textfixer/pkg/speller"

type Service struct {
	*speller.Client
}

func NewService() *Service {
	return &Service{
		Client: speller.NewClient(),
	}
}
