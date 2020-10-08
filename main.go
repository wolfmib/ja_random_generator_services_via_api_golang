package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Jira-Ticket:  DAT-148
/*   	Intn 5 is from 0,1,2,3,4
		With low = 1, high = 5
        return value is 1,2,3,4,5
*/
func get_random(low int, high int) int {
	return rand.Intn(high-low+1) + low
}

func jason_random_v1(low int, high int) int {
	get_times := get_random(1, 5)
	var value int
	for i := 0; i < get_times; i++ {
		value = get_random(low, high)
	}
	return value
}

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type RandomGen struct {
	Value int `json:"value,omitempty" bson:"value,omitempty`
}

// Global client
var client *mongo.Client

func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var my_user User

	/*
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

	log.Println(request.Body)

	//custom_user_map := make(map[string]interface{})
	//custom_user_map["firstname"] = "This is interface"
	//custom_user_map["lastname"] = "hahahaha"

	err := json.NewDecoder(request.Body).Decode(&my_user)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(my_user)
	// Johnny: this interface, will help you to dynamic access the api data format
	// no need to limit the structure in struct.

	log.Println("[INFO]: Get the request")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("test_user").Collection("test_user_collection")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("my user")
	log.Println(my_user)

	// Check the data is not the same

	var _key string = "firstname"
	var _value string = my_user.FirstName

	// cursor, err := collection.Find(ctx, bson.M{"firstname": "hi"})
	cursor, err := collection.Find(ctx, bson.M{_key: _value})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	defer cursor.Close(ctx)

	var _tem_user User
	for cursor.Next(ctx) {

		log.Println("####### message  Duplicate Data ...##############")
		log.Println("Input firstname", _value)
		cursor.Decode(&_tem_user)
		fmt.Println(_tem_user)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "Duplicate Data ..."} `))

		return
	}

	log.Println("Saveing data ...")
	result, _ := collection.InsertOne(ctx, my_user)
	json.NewEncoder(response).Encode(result)
}

// Query (min, max )
func GetRandomEndpoint(response http.ResponseWriter, request *http.Request) {

	// Header
	response.Header().Add("content-type", "application/json")

	var get_random RandomGen

	// load dynamic endpoint request to parameters
	random_dynamic_parameter := mux.Vars(request)
	var low int
	var high int
	low, _ = strconv.Atoi(random_dynamic_parameter["low"])
	high, _ = strconv.Atoi(random_dynamic_parameter["high"])

	get_random.Value = jason_random_v1(low, high)
	fmt.Println(get_random.Value)

	json.NewEncoder(response).Encode(get_random)
	//response.Write([]byte(get_random.value))

}

// Query one specific user by id
func CetUser_by_id_Endpoint(response http.ResponseWriter, request *http.Request) {

	// Header
	response.Header().Add("content-type", "application/json")

	// Make user structure instance
	var user User

	// load dynamic endpoint request to parameters
	user_dynamic_parameter := mux.Vars(request)

	// Convert to mongo object id (I won't use that )
	id, _ := primitive.ObjectIDFromHex(user_dynamic_parameter["id"])
	log.Println("Get id:    , ", id)

	// Setting ctx
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("test_user").Collection("test_user_collection")

	// Query ici
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	// Drror checking
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": err for query by find()}`))
		return
	}

	// Make json response
	log.Println("==================")
	log.Println("Return the query message user: ", user)
	log.Println("==================")
	json.NewEncoder(response).Encode(user)

}

// Query one specific user by firstname
// /user/name/{firstname}
func CetUser_by_name_Endpoint(response http.ResponseWriter, request *http.Request) {

	// Header
	response.Header().Add("content-type", "application/json")

	// Make user structure instance
	var user User

	// load dynamic endpoint request to parameters
	user_dynamic_parameter := mux.Vars(request)

	// Convert to mongo object id (I won't use that )
	firstname := user_dynamic_parameter["firstname"]
	log.Println("Get firstname:    , ", firstname)

	// Setting ctx
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("test_user").Collection("test_user_collection")

	// Query ici
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"firstname": firstname}).Decode(&user)

	// Drror checking
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": err for query by find()}`))
		return
	}

	// Make json response
	log.Println("==================")
	log.Println("Return the query message user: ", user)
	log.Println("==================")
	json.NewEncoder(response).Encode(user)

}

func main() {

	fmt.Println("Start the application ....")
	fmt.Println("Please source ja_test_get_random_api.sh to test the api ")

	fmt.Println("Test jason_random_v1(1,7)")
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println(jason_random_v1(1, 7))
	fmt.Println("--------------------------------")

	//Router
	router := mux.NewRouter()
	router.HandleFunc("/user", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/get_random/{low}/{high}", GetRandomEndpoint).Methods("GET")      //Test 1/7, return 1~7 random
	router.HandleFunc("/user/{id}", CetUser_by_id_Endpoint).Methods("GET")               //User, one user with id
	router.HandleFunc("/user/name/{firstname}", CetUser_by_name_Endpoint).Methods("GET") //User, one user with firstanme
	http.ListenAndServe(":12345", router)

}
