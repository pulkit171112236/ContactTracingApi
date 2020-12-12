package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    // "go.mongodb.org/mongo-driver/mongo/readpref"
)


var client *mongo.Client



type Marshaler interface {
    MarshalJSON() ([]byte, error)
}

type JSONTime time.Time

func (t JSONTime)MarshalJSON() ([]byte, error) {
    //do your serializing here
    stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
    return []byte(stamp), nil
}


type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Dob string `json:"dob,omitempty" bson: "dob,omitempty"`
	PhoneNumber string `json:"phone_num,omitempty" bson: "phone_num,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Time_stamp time.Time `json:"time_stamp" bson: "time_stamp"`
}
type Contact struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId1 string             `json:"user_id_1,omitempty" bson:"user_id_1,omitempty"`
	UserId2 string `json:"user_id_2,omitempty" bson: "user_id_2,omitempty"`
	Time_stamp time.Time `json:"time_stamp" bson: "time_stamp"`
}


// type User struct {

func GetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client.Database("contactTracingApi").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func GetAllUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var users []User
	collection := client.Database("contactTracingApi").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(users)
}

func CreateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)

	user.Time_stamp = time.Now()

	collection := client.Database("contactTracingApi").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}



func CreateContact(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var contact Contact
	_ = json.NewDecoder(request.Body).Decode(&contact)

	contact.Time_stamp = time.Now()

	collection := client.Database("contactTracingApi").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, _ := collection.InsertOne(ctx, contact)
	json.NewEncoder(response).Encode(result)
}


func GetContact(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	var contact Contact
	collection := client.Database("contactTracingApi").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	
	err := collection.FindOne(ctx, Contact{ID: id}).Decode(&contact)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(contact)
}

func GetAllContacts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var contacts []Contact
	collection := client.Database("contactTracingApi").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var contact Contact
		cursor.Decode(&contact)
		contacts = append(contacts, contact)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(contacts)
}

func main() {
	fmt.Println("Starting the application...")
	// for connecting to mongodb_server
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()

	router.HandleFunc("/users", GetAllUsers).Methods("GET")
	router.HandleFunc("/users", CreateUser).Methods("POST")

	router.HandleFunc("/users/{id}", GetUser).Methods("GET")

	router.HandleFunc("/contacts", GetAllContacts).Methods("GET")
	router.HandleFunc("/contacts", CreateContact).Methods("POST")

	router.HandleFunc("/contacts?user={id}&infection_timestamp={ts}", GetContact).Methods("GET")

	http.ListenAndServe(":12345", router)
}

