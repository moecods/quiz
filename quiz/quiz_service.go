package quiz

import (
	"moecods/quiz/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizService struct {
	repo         QuizRepository
	timeProvider utils.TimeProvider
}

func NewQuizService(repo QuizRepository) *QuizService {
	return &QuizService{
		repo:         repo,
		timeProvider: utils.RealTimeProvider{},
	}
}

func (s *QuizService) AddQuiz(quiz Quiz) (Quiz, error) {
	questions := quiz.Questions
	for i := range questions {
		questions[i].ID = primitive.NewObjectID()
	}

	quiz.Questions = questions
	quiz.CreatedAt = s.timeProvider.Now()
	quiz.UpdatedAt = s.timeProvider.Now()

	_, err := s.repo.AddQuiz(&quiz)
	return quiz, err
}

func (s *QuizService) GetQuiz(id primitive.ObjectID) (*Quiz, error) {
	quiz, error := s.repo.GetQuiz(id)
	return quiz, error
}
