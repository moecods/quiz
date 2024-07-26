package quiz

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @swagger:model Quiz
type Quiz struct {
	QuizBase
	Questions []Question `bson:"questions" json:"questions"`
}

// @swagger:model QuizBase
type QuizBase struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	StartAt     time.Time          `bson:"start_at" json:"start_at" example:"2024-07-30T12:00:00Z"`
	EndAt       time.Time          `bson:"end_at" json:"end_at" example:"2024-07-30T12:00:00Z"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at" swaggerignore:"true"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at" swaggerignore:"true"`
}

type Question struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type          string             `bson:"type" json:"type"` // Type of question: "descriptive" or "multiple-choice"
	Text          string             `bson:"text" json:"text"`
	Options       []string           `bson:"options" json:"options"`
	CorrectOption int                `bson:"correct_option" json:"correct_option"`
}
