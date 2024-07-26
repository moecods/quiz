package quiz

import (
	"moecods/quiz/utils"
	"time"

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

func (s *QuizService) SaveQuestionsToQuiz(quiz *Quiz, newQuestions []Question) {
	questionMap := make(map[primitive.ObjectID]int)

	for i, q := range quiz.Questions {
		questionMap[q.ID] = i
	}

	for _, newQuestion := range newQuestions {
		index, exists := questionMap[newQuestion.ID]
		if exists {
			quiz.Questions[index] = newQuestion
		} else {
			if newQuestion.ID == primitive.NilObjectID {
				newQuestion.ID = primitive.NewObjectID()
			}

			quiz.Questions = append(quiz.Questions, newQuestion)
		}
	}
}

func (s *QuizService) getQuizEndTime(id primitive.ObjectID) (time.Time, error) {
	quiz, err := s.GetQuiz(id)
	if err != nil {
		return time.Time{}, err
	}

	return quiz.EndAt, err
}

// func (s *QuizService) getQuizStartTime(id primitive.ObjectID) (time.Time, error) {
// 	quiz, err := s.GetQuiz(id)
// 	if err != nil {
// 		return time.Time{}, err
// 	}

// 	return quiz.StartAt, err
// }