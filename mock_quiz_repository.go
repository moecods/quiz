package main

import (
	"fmt"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockQuizRepository struct {
	mock.Mock
}

func (m *MockQuizRepository) ListQuizzes() ([]Quiz, error) {
	fmt.Printf("Mock list quizzes")
	return []Quiz{}, nil
}

func (m *MockQuizRepository) AddQuiz(quiz *Quiz) error {
	fmt.Printf("Mock Add quiz")
	return nil
}

func (m *MockQuizRepository) UpdateQuiz(id primitive.ObjectID, quiz *Quiz) error {
	fmt.Printf("Mock Update quiz")
	return nil
}

func (m *MockQuizRepository) DeleteQuiz(id primitive.ObjectID) error {
	fmt.Printf("Mock Delete quiz")
	return nil
}

func (m *MockQuizRepository) GetQuiz(id primitive.ObjectID) (*Quiz, error) {
	fmt.Printf("Mock Get quiz")
	return &Quiz{}, nil
}
