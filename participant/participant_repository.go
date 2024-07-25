package participant

import (
	"log"
	"moecods/quiz/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ParticipantRepository struct {
	collection *mongo.Collection
}

func NewParticipantRepository(collection *mongo.Collection) *ParticipantRepository {
	return &ParticipantRepository{collection: collection}
}

func (r *ParticipantRepository) ListParticipantzes() ([]Participant, error) {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var participantzes []Participant
	for cursor.Next(ctx) {
		var participant Participant
		err := cursor.Decode(&participant)
		if err != nil {
			log.Fatal(err)
		}
		participantzes = append(participantzes, participant)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return participantzes, nil
}

func (r *ParticipantRepository) GetParticipantsByQuiz(quizId primitive.ObjectID) ([]ParticipantBase, error) {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	projection := bson.M{"answers": 0}
	cursor, err := r.collection.Find(ctx, bson.M{"quiz_id": quizId}, options.Find().SetProjection(projection))
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(ctx)

	var participants []ParticipantBase
	for cursor.Next(ctx) {
		var participant ParticipantBase
		err := cursor.Decode(&participant)
		if err != nil {
			log.Fatal(err)
		}
		participants = append(participants, participant)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return participants, nil
}

func (r *ParticipantRepository) AddParticipant(participant *Participant) error {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, participant)
	participant.ParticipantBase.ID = result.InsertedID.(primitive.ObjectID)
	return err
}

func (r *ParticipantRepository) AddManyParticipants(participants []Participant) error {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	docs := make([]interface{}, len(participants))
	for i, participant := range participants {
		docs[i] = participant
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

func (r *ParticipantRepository) UpdateParticipant(id primitive.ObjectID, participant *Participant) error {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": participant,
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *ParticipantRepository) DeleteParticipant(id primitive.ObjectID) error {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ParticipantRepository) GetParticipant(id primitive.ObjectID) (*Participant, error) {
	ctx, cancel := utils.WithTimeoutContext(5 * time.Second)
	defer cancel()

	var participant Participant
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&participant)
	return &participant, err
}
