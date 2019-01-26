package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var dao *DBServer

type People struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	FullName string        `bson:"fullname" json:"fullname"`
	Phone    string        `bson:"phone" json:"phone"`
	Age      string        `bson:"age" json:"age"`
}

func GetAllPeoples(w http.ResponseWriter, r *http.Request) {
	peoples, err := dao.FindAll()
	if err != nil {
		http.Error(w, "Error to fetch data from db", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(peoples)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	people, err := dao.FindById(params["id"])
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	resp, _ := json.Marshal(people)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func CreatePeople(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var p People

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.ID = bson.NewObjectId()
	if err := dao.Insert(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func UpdatePeople(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var p People

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err := dao.Update(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeletePeople(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var p People

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err := dao.Delete(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func init() {
	dao = &DBServer{Server: "localhost:27017", Database: "address_book"}
	dao.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/peoples", GetAllPeoples).Methods("GET")
	r.HandleFunc("/peoples/{id}", GetPeople).Methods("GET")
	r.HandleFunc("/peoples", CreatePeople).Methods("POST")
	r.HandleFunc("/peoples", UpdatePeople).Methods("PUT")
	r.HandleFunc("/peoples", DeletePeople).Methods("DELETE")

	log.Println("Starting API")
	log.Fatal(http.ListenAndServe(":8080", r))
}
