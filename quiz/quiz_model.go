package quiz

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
    QuizBase
    Questions []Question `bson:"questions" json:"questions"`
}

type QuizBase struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Title       string             `bson:"title" json:"title"`
    Description string             `bson:"description" json:"description"`
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Question struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type          string             `bson:"type" json:"type"` // Type of question: "descriptive" or "multiple-choice"
	Text          string             `bson:"text" json:"text"`
	Options       []string           `bson:"options" json:"options"`
	CorrectOption int                `bson:"correct_option" json:"correct_option"`
}
