package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Game represents a chess game
type Game struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	Player1     string    `json:"player1,omitempty" bson:"player1,omitempty"`
	Player2     string    `json:"player2,omitempty" bson:"player2,omitempty"`
	Moves       []string  `json:"moves,omitempty" bson:"moves,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	LastUpdated time.Time `json:"lastUpdated,omitempty" bson:"lastUpdated,omitempty"`
}

var client *mongo.Client

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get MongoDB connection URI from environment variables
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI")
	}

	// Create MongoDB client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Initialize router
	router := mux.NewRouter()

	// Define API endpoints
	router.HandleFunc("/games", createGame).Methods("POST")
	router.HandleFunc("/games/{id}", getGame).Methods("GET")
	router.HandleFunc("/games/{id}", updateGame).Methods("PUT")
	router.HandleFunc("/games/{id}", deleteGame).Methods("DELETE")

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// Helper function to get the MongoDB collection
func getCollection() *mongo.Collection {
	return client.Database("chess").Collection("games")
}

// Handler function to create a new game
func createGame(w http.ResponseWriter, r *http.Request) {
	var game Game
	json.NewDecoder(r.Body).Decode(&game)
	game.CreatedAt = time.Now()
	game.LastUpdated = game.CreatedAt
	_, err := getCollection().InsertOne(context.Background(), game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Handler function to get a game by ID
func getGame(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var game Game
	err := getCollection().FindOne(context.Background(), bson.M{"_id": id}).Decode(&game)
	if err != nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(game)
}

// Handler function to update a game by ID
func updateGame(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var game Game
	json.NewDecoder(r.Body).Decode(&game)
	game.LastUpdated = time.Now()
	_, err := getCollection().ReplaceOne(context.Background(), bson.M{"_id": id}, game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Handler function to delete a game by ID
func deleteGame(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	_, err := getCollection().DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
