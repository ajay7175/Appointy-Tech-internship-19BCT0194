package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

var client *mongo.Client

type Instauser struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}
type Instapost struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserPostId string             `json:"userpostid,omitempty" bson:"userpostid,omitempty"`
	Caption    string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageUrl   string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	Timestamp  string             `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

func createuser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var instauser Instauser
	_ = json.NewDecoder(request.Body).Decode((&instauser))
	hashpassword := GetMD5Hash(instauser.Password) // Will not be able to retrieve back
	instauser.Password = hashpassword
	//fmt.Println(instauser.Password)
	collection := client.Database("appointy").Collection("instausers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, instauser)
	json.NewEncoder(response).Encode(result)
}

func getUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	//params=
	//id, _ := primitive.ObjectIDFromHex(params["id"])
	//urlParams := request.URL.Query()
	temp := strings.Split(request.URL.Path, "/")
	uid := temp[len(temp)-1]
	id, _ := primitive.ObjectIDFromHex(uid)
	//fmt.Println(id)
	var instauser Instauser
	collection := client.Database("appointy").Collection("instausers")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Instauser{ID: id}).Decode(&instauser)
	//fmt.Println(instauser.Email)
	//temb:=db.instausers.FindOne({"_id":"6161659378289d3f64f6e2ad"})
	//fmt.println(temb)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(instauser)

}
func createpost(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var instapost Instapost
	_ = json.NewDecoder(request.Body).Decode((&instapost))

	collection := client.Database("appointy").Collection("instaposts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, instapost)
	json.NewEncoder(response).Encode(result)
}
func getpost(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	temp := strings.Split(request.URL.Path, "/")
	uid := temp[len(temp)-1]
	id, _ := primitive.ObjectIDFromHex(uid)
	var instapost Instapost
	collection := client.Database("appointy").Collection("instaposts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Instapost{ID: id}).Decode(&instapost)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(instapost)

}

func getallpost(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	temp := strings.Split(request.URL.Path, "/")
	uid := temp[len(temp)-1]
	//id, _ := primitive.ObjectIDFromHex(uid)
	var instapost []Instapost
	collection := client.Database("appointy").Collection("instaposts")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.M{})
	for cur.Next(ctx) {
		var instaposts Instapost
		cur.Decode(&instaposts)
		if instaposts.UserPostId == uid {
			instapost = append(instapost, instaposts)
		}
	}
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	//err, := collection.Find(ctx,Instapost{UserPostId: uid}).Decode(&instapost)
	//err, _ := collection.Find(ctx, bson.M{"userpostid": uid})
	//cur,currErr :=collective.Find(ctx, bson.M{""})

	// err := collection.FindOne(ctx, bson.M{"userpostid": uid}).Decode(&instapost)

	// } else {
	// 	// Print out data from the document result
	// 	fmt.Println("result AFTER:", instapost, "\n")
	//err =db.collection.Find(ctx,{"userpostid":uid});
	// if err != nil {
	// 	response.WriteHeader(http.StatusInternalServerError)
	// 	//response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	// 	return
	// }
	json.NewEncoder(response).Encode(instapost)
}

func main() {
	fmt.Println("Hi vasuki")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	http.HandleFunc("/user", createuser)
	http.HandleFunc("/user/", getUser)
	http.HandleFunc("/post", createpost)
	http.HandleFunc("/post/", getpost)
	http.HandleFunc("/post/user/", getallpost)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
