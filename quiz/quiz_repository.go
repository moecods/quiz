package quiz

import (
	"log"
	"moecods/quiz/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuizRepository interface {
	ListQuizzes() ([]QuizBase, error)
	AddQuiz(quiz *Quiz) (*Quiz, error)
	UpdateQuiz(id primitive.ObjectID, quiz *Quiz) error
	DeleteQuiz(id primitive.ObjectID) error
	GetQuiz(id primitive.ObjectID) (*Quiz, error)
}

type MongoQuizRepository struct {
	collection *mongo.Collection
}

func NewQuizRepository(collection *mongo.Collection) *MongoQuizRepository {
	return &MongoQuizRepository{collection: collection}
}

func (r *MongoQuizRepository) ListQuizzes() ([]QuizBase, error) {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	projection := bson.M{"questions": 0}
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var quizzes []QuizBase
	for cursor.Next(ctx) {
		var quiz QuizBase
		err := cursor.Decode(&quiz)
		if err != nil {
			log.Fatal(err)
		}

		quizzes = append(quizzes, quiz)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return quizzes, nil
}

func (r *MongoQuizRepository) AddQuiz(quiz *Quiz) (*Quiz, error) {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, quiz)
	if err != nil {
		return nil, err
	}

	quiz.ID = result.InsertedID.(primitive.ObjectID)
	return quiz, err
}

func (r *MongoQuizRepository) UpdateQuiz(id primitive.ObjectID, quiz *Quiz) error {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": quiz,
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoQuizRepository) DeleteQuiz(id primitive.ObjectID) error {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoQuizRepository) GetQuiz(id primitive.ObjectID) (*Quiz, error) {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	var quiz Quiz
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&quiz)
	return &quiz, err
}
