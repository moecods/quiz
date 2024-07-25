package participant

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParticipantBase struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	QuizID     primitive.ObjectID `bson:"quiz_id" json:"quiz_id"`
	Status     string             `bson:"status" json:"status"` // not_started, started, finished
	StartAt    time.Time
	FinishedAt time.Time
}

type Answer struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	QuestionID     primitive.ObjectID `bson:"question_id" json:"question_id"`
	AnswerText     string             `bson:"answer_text" json:"answer_text"`
	SelectedOption int                `bson:"selection_option" json:"selection_option"`
	IsCorrect      bool               `bson:"is_correct" json:"is_correct"`
	AnsweredAt     time.Time          `bson:"answered_at" json:"answered_at"`
}

type Participant struct {
	ParticipantBase
	Answers         []Answer `bson:"answers" json:"answers"`
}
