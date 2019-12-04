package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

const (
	PORT = ":8090"
)

type Booking struct {
	ID      int    `json:"id"`
	User    string `json:"user"`
	Members int    `json:"members"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: HomePage")
}

func handleRequests() {
	log.Println("Starting development server at http://127.0.0.1:8090/")
	log.Println("Quit the server with CONTROL-C.")
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/new-booking", createNewBooking).Methods("POST")
	myRouter.HandleFunc("/all-bookings", returnAllBookings).Methods("GET")
	myRouter.HandleFunc("/booking/{id}", returnSingleBooking).Methods("GET")
	myRouter.HandleFunc("/booking/{id}", deleteBooking).Methods("DELETE")
	myRouter.HandleFunc("/booking/{id}", updateBooking).Methods("PUT")
	log.Fatal(http.ListenAndServe(PORT, myRouter))
}

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func createNewBooking(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var booking Booking
	json.Unmarshal(reqBody, &booking)
	db.Create(&booking)
	fmt.Println("Endpoint Hit: Creating New Booking")
	respondJSON(w, http.StatusCreated, booking)
}

func returnAllBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Location")
	bookings := []*Booking{}
	db.Find(&bookings)
	fmt.Println("Endpoint Hit: returnAllBookings")
	respondJSON(w, http.StatusOK, bookings)

}

func returnSingleBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	bookings := []Booking{}
	db.Find(&bookings)
	for _, booking := range bookings {
		// string to int
		s, err := strconv.Atoi(key)
		if err == nil {
			if booking.ID == s {
				fmt.Println(booking)
				fmt.Println("Endpoint Hit: Booking No:", key)
				respondJSON(w, http.StatusOK, booking)
			}
		}
	}
}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var booking Booking
	db.First(&booking, params["id"])
	db.Delete(&booking)

	var bookings []Booking
	db.Find(&bookings)
	// json.NewEncoder(w).Encode(&bookings)
	respondJSON(w, http.StatusNoContent, booking)
}

func updateBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["user"]
	booking := getBookingOr404(db, name, w, r)
	if booking == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&booking); err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&booking).Error; err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	respondJSON(w, http.StatusOK, booking)
}

func getBookingOr404(db *gorm.DB, user string, w http.ResponseWriter, r *http.Request) *Booking {
	booking := Booking{}
	if err := db.First(&booking, Booking{User: user}).Error; err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return nil
	}
	return &booking
}

// func withCORS(fn http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Expose-Headers", "Location")

// 	}
// }

func main() {
	// Please define your username and password for MySQL.
	db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:8889)/soccer?charset=utf8&parseTime=True")
	// NOTE: See we’re using = to assign the global var
	// instead of := which would assign it only in this function

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Connection Established")
	}

	db.AutoMigrate(&Booking{})
	handleRequests()

}
