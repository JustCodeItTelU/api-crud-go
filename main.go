package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

var db *gorm.DB
var err error

type Topic struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type Details struct {
	ID       int         `json:"id"`
	Title    string      `json:"title"`
	Content  string      `json:"content"`
	Comments interface{} `json:"comments"`
}

type Comments struct {
	ID      int    `json:"id"`
	Comment string `json:"comment"`
	IDTopic int    `json:"id_topic"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	db, err = gorm.Open("mysql", "root:1235813@/db_jci?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection Failed", err)
	} else {
		log.Println("Connected Success")
	}
	db.AutoMigrate(&Topic{})
	db.AutoMigrate(&Comments{})

	handleRequest()
}

func handleRequest() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/topics", createTopic).Methods("POST")
	myRouter.HandleFunc("/topics", getTopics).Methods("GET")
	myRouter.HandleFunc("/topics/{id}", getTopic).Methods("GET")
	myRouter.HandleFunc("/topics/{id}", updateTopic).Methods("PUT")
	myRouter.HandleFunc("/topics/{id}", deleteTopic).Methods("DELETE")
	myRouter.HandleFunc("/topics/{id}", createComment).Methods("POST")
	myRouter.HandleFunc("/topics/{id}/{id_comment}", updateComment).Methods("PUT")
	myRouter.HandleFunc("/topics/{id}/{id_comment}", deleteComment).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func createTopic(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var topic Topic
	json.Unmarshal(payloads, &topic)

	db.Create(&topic)

	res := Response{
		Code:    200,
		Message: "Success Create Topic",
		Data:    topic,
	}
	response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func getTopics(w http.ResponseWriter, r *http.Request) {
	topics := []Topic{}

	db.Find(&topics)

	res := Response{
		Code:    200,
		Message: "Succes get All Topic",
		Data:    topics,
	}
	responses, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responses)
}
func getTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["id"]

	var topic Topic
	comment := []Comments{}
	//var topicDetails TopicDetails

	db.First(&topic, topicID)
	db.Find(&comment, "id_topic = ? ", topicID)

	res := Response{
		Code:    200,
		Message: "Succes get Topic",
		Data: Details{
			ID:       topic.ID,
			Title:    topic.Title,
			Content:  topic.Content,
			Comments: comment,
		},
	}
	response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func updateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicId := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var topicUpdates Topic

	json.Unmarshal(payloads, &topicUpdates)
	var topic Topic
	db.First(&topic, topicId)
	db.Model(&topic).Updates(topicUpdates)

	res := Response{
		Code:    200,
		Message: "Success Update Topic",
		Data:    topic,
	}
	response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func deleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicId := vars["id"]

	var topic Topic

	db.First(&topic, topicId)
	db.Delete(&topic)

	res := Response{
		Code:    200,
		Message: "Success Delete Topic",
	}
	response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
func createComment(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var comment Comments

	json.Unmarshal(payloads, &comment)
	vars := mux.Vars(r)
	topicID := vars["id"]

	//var topic Topic

	db.Exec("INSERT INTO comments (comment, id_topic) VALUES ('" + comment.Comment + "'," + topicID + ")")
	res := Response{
		Code:    200,
		Message: "Success Create Comment",
		Data:    comment.Comment,
	}
	response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func updateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicId := vars["id"]
	commentId := vars["id_comment"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var commentUpdates Comments

	json.Unmarshal(payloads, &commentUpdates)

	db.Exec("UPDATE comments SET comment = '" + commentUpdates.Comment + "' WHERE id = " + commentId + " AND " + "id_topic = " + topicId)

	res := Response{
		Code:    200,
		Message: "Success Create Comment",
		Data:    commentUpdates.Comment,
	}
	response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
func deleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicId := vars["id"]
	commentId := vars["id_comment"]

	db.Exec("DELETE FROM comments WHERE id = " + commentId + " AND " + "id_topic = " + topicId)
	res := Response{
		Code:    200,
		Message: "Success Delete Comment",
	}
	response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
