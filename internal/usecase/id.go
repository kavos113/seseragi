package usecase

import "github.com/google/uuid"

type IDGenerator interface {
	GenerateID() string
}

type UUIDGenerator struct{}

func (g *UUIDGenerator) GenerateID() string {
	return uuid.New().String()
}

func NewUUIDGenerator() IDGenerator {
	return &UUIDGenerator{}
}