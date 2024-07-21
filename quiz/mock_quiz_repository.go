package quiz

import (
	"fmt"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockQuizRepository struct {
	mock.Mock
}

func (m *MockQuizRepository) ListQuizzes() ([]QuizBase, error) {
	fmt.Printf("Mock list quizzes")
	return []QuizBase{}, nil
}

func (m *MockQuizRepository) AddQuiz(quiz *Quiz) (*Quiz, error) {
	args := m.Called(quiz)
	if quiz := args.Get(0); quiz != nil {
		return quiz.(*Quiz), args.Error(1)
	}
	return quiz, args.Error(1)
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
	args := m.Called(id)
	if quiz := args.Get(0); quiz != nil {
		return quiz.(*Quiz), args.Error(1)
	}
	return nil, args.Error(1)
}
