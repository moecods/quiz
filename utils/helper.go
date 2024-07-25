package utils

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ExtractIDFromRequest(r *http.Request) (primitive.ObjectID, error) {
    vars := mux.Vars(r)
    idStr := vars["id"]

    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        return primitive.ObjectID{}, err
    }
    return id, nil
}