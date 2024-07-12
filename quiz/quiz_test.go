package quiz

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockTimeProvider struct {
	mockNow func() time.Time
}

func (m MockTimeProvider) Now() time.Time {
	return m.mockNow()
}

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

func TestQuizHandler_AddQuiz(t *testing.T) {
	t.Run("successful add quiz", func(t *testing.T) {
		mockRepo := new(MockQuizRepository)

		fixedTime := time.Date(2024, time.July, 11, 12, 0, 0, 0, time.UTC)
		mockTimeProvider := MockTimeProvider{
			mockNow: func() time.Time { return fixedTime },
		}

		service := &QuizService{repo: mockRepo, timeProvider: mockTimeProvider}

		quiz := Quiz{
			Title:       "New Quiz",
			Description: "New Description",
			Questions: []Question{
				{Type: "descriptive", Text: "2+2=?"},
			},
			CreatedAt: mockTimeProvider.Now(),
			UpdatedAt: mockTimeProvider.Now(),
		}

		mockRepo.On("AddQuiz", &quiz).Return(&quiz, nil)

		storedQuiz, err := service.AddQuiz(quiz)

		assert.NoError(t, err)

		assert.Equal(t, quiz, storedQuiz)

		mockRepo.AssertCalled(t, "AddQuiz", &quiz)
	})

	t.Run("validation failed to add quiz", func(t *testing.T) {
		quiz := Quiz{
			Title: "",
		}

		validation := QuizValidation{quiz: &quiz}
		err := validation.QuizValidation()
		assert.Error(t, err)
		assert.EqualErrorf(t, err, "title is required", "Error should be: %v, got: %v", "title is required", err)
	})

	// t.Run("failed to add quiz", func(t *testing.T) {
	//     mockRepo := new(MockQuizRepository)
	//     service := &QuizService{repo: mockRepo}

	//     quiz := Quiz{
	//         Title:       "New Quiz",
	//         Description: "New Description",
	//         Questions: []Question{
	//             {Type: "descriptive", Text: "2+2=?"},
	//         },
	//     }

	//     expectedError := errors.New("failed to add quiz")
	//     mockRepo.On("AddQuiz", quiz).Return(Quiz{}, expectedError)

	//     storedQuiz, err := service.AddQuiz(quiz)

	//     assert.Error(t, err)
	//     assert.Equal(t, expectedError, err)

	//     assert.Equal(t, Quiz{}, storedQuiz)

	//     mockRepo.AssertCalled(t, "AddQuiz", quiz)
	// })
}
