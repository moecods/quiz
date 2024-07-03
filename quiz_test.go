package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestQuizeHandler_GetQuiz(t *testing.T) {
	mockRepo := new(MockQuizRepository)
	service := &QuizService{repo: mockRepo}

	QuizId := primitive.NewObjectID()
	mockRepo.On("GetQuiz", QuizId).Return(&Quiz{ID: QuizId, Title: "Me Before You"}, nil)
	quiz, err := service.GetQuiz(QuizId)
	assert.NoError(t, err)
	assert.Equal(t, QuizId, quiz.ID)

	QuizId = primitive.NewObjectID()
	mockRepo.On("GetQuiz", QuizId).Return(nil, errors.New("post not found"))
	_, err = service.GetQuiz(QuizId)
	assert.Error(t, err)
}
