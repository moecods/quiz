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

func TestQuizService_AddQuiz(t *testing.T) {
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

	t.Run("sucessful save new question in quiz", func(t *testing.T) {

		existQuestionId1 := primitive.NewObjectID()
		existQuestionId2 := primitive.NewObjectID()

		quiz := Quiz{
			Title: "New Quiz",
			Questions: []Question{
				{ID: existQuestionId1, Type: "descriptive", Text: "2+2=?"},
				{ID: existQuestionId2, Type: "descriptive", Text: "5+5=?"},
			},
		}

		newQuestions := []Question{
			{Type: "descriptive", Text: "3+3=?"},
			{Type: "descriptive", Text: "4+4=?"},
			{ID: existQuestionId2, Type: "descriptive", Text: "10+10=?"},
		}

		service := &QuizService{}
		service.SaveQuestionsToQuiz(&quiz, newQuestions)

		// Check that the existing question remains unchanged
		assert.Equal(t, existQuestionId1, quiz.Questions[0].ID, "Existing question ID should not change")
		assert.Equal(t, "descriptive", quiz.Questions[0].Type, "Existing question type should not change")
		assert.Equal(t, "2+2=?", quiz.Questions[0].Text, "Existing question text should not change")

		// Check that the existing question update with new data
		assert.Equal(t, existQuestionId2, quiz.Questions[1].ID, "Existing question ID should not change")
		assert.Equal(t, "10+10=?", quiz.Questions[1].Text, "Existing question text should updated with new data")

		// Check that new questions are added with new IDs
		for i := 2; i < 4; i++ {
			assert.NotEqual(t, primitive.NilObjectID, quiz.Questions[i].ID, "New question ID should be set")
			assert.Contains(t, []string{"3+3=?", "4+4=?"}, quiz.Questions[i].Text, "New question text should match input")
		}
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
