package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sntruong12/josie-rides/internal/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Josie Rides" as the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// handles edge case when users nav to non existing routes
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	rides, err := app.rides.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Rides = rides

	// Use the new render helper.
	app.render(w, http.StatusOK, "home.html", data)
}

// Add a rideView handler function.
func (app *application) rideView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	ride, err := app.rides.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Ride = ride

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) rideCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.html", data)
}

// Add a rideCreate handler function.
func (app *application) rideCreatePost(w http.ResponseWriter, r *http.Request) {
	// Limit the request body size to 4096 bytes
	// r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	trailName := r.PostForm.Get("trail_name")
	distance := r.PostForm.Get("distance_miles")
	durationHours := r.PostForm.Get("duration_hours")
	durationMinutes := r.PostForm.Get("duration_minutes")
	durationSeconds := r.PostForm.Get("duration_seconds")
	rodeAt := r.PostForm.Get("rode_at")
	timezone := r.PostForm.Get("timezone")
	// need to implement this later
	media := json.RawMessage("null")

	// validation
	if title == "" || len(title) > 255 || description == "" || trailName == "" || len(trailName) > 255 || distance == "" || len(distance) > 8 || durationHours == "" || durationMinutes == "" || durationSeconds == "" || rodeAt == "" || timezone == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// convert distance to float64
	distanceFloat, err := strconv.ParseFloat(distance, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// convert rodeAt to time.Time utc
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	rodeAtTime, err := time.ParseInLocation("2006-01-02T15:04", rodeAt, loc)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	rodeAtUTC := rodeAtTime.UTC()

	// convert duration into seconds
	h, _ := strconv.Atoi(durationHours)
	m, _ := strconv.Atoi(durationMinutes)
	s, _ := strconv.Atoi(durationSeconds)
	convertedDuration := (h * 3600) + (m * 60) + s

	id, err := app.rides.Create(title, description, trailName, distanceFloat, convertedDuration, rodeAtUTC, media)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ride/view/%d", id), http.StatusSeeOther)
}
