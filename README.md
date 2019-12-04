# GO API

## Purpose
<p>This is a simple example on implementing a CRUD api using GO Lang.</p>

## Packages
<p>Three main packages are used<br>
	"github.com/gorilla/mux" - routing<br>
	"github.com/jinzhu/gorm" - orm library for GO<br>
	_ "github.com/jinzhu/gorm/dialects/mysql" - mysql<br>
</p>

## Installation
<p>Clone the repo<br>
Run go get in terminal for the three packages used above eg: go get github.com/gorilla/mux<br>
Add your database connections by replacing the texts in caps on the line below, also note the port and host of the mysql db my port is 8889 while my host is 127.0.0.1 <br>
    db, err = gorm.Open("mysql", "ROOT:ROOT@tcp(127.0.0.1:8889)/soccer?charset=utf8&parseTime=True")<br>
Edit PORT = ":8090" to suit your own port for running the app<br>
</p>

## Routes
<p>/ - homepage(GET)<br>
/new-booking(POST) - create new booking using USER and MEMBERS(int) fields<br>
/all-bookings(GET) - all bookings<br>
/booking/{id}(GET) - get single booking with ID<br>
/booking/{id}(DELETE) - deletes single booking with ID<br>
/booking/{id}(PUT) - updates set USER and MEMBERS(int) fields.
</p>

## TESTING
<p>
Test with postman referring to the routes above or import the collection below<br>
https://www.getpostman.com/collections/30f69be2e9db5e63d136
</p>