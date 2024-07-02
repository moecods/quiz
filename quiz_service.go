package main

import "go.mongodb.org/mongo-driver/bson/primitive"



type QuizService struct {
	repo QuizRepository
}

func NewQuizService(repo QuizRepository) *QuizService {
    return &QuizService{repo: repo}
}

func (s *QuizService) AddQuiz(quiz Quiz) (Quiz, error) {
	questions := quiz.Questions
	for i := range questions {
		questions[i].ID = primitive.NewObjectID()
	}

	quiz.Questions = questions

	err := s.repo.AddQuiz(&quiz)
	
	return quiz, err
}