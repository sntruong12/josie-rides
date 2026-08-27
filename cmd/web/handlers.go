package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sntruong12/josie-rides/internal/models"
	"github.com/sntruong12/josie-rides/internal/validator"
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
	data.Form = rideCreateForm{}

	app.render(w, http.StatusOK, "create.html", data)
}

type rideCreateForm struct {
	Title           string
	Description     string
	TrailName       string
	Distance        string
	DurationHours   string
	DurationMinutes string
	DurationSeconds string
	RodeAt          string
	Timezone        string
	Media           string

	validator.Validator
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

	form := rideCreateForm{
		Title:           strings.TrimSpace(r.PostForm.Get("title")),
		Description:     strings.TrimSpace(r.PostForm.Get("description")),
		Distance:        r.PostForm.Get("distance_miles"),
		DurationHours:   r.PostForm.Get("duration_hours"),
		DurationMinutes: r.PostForm.Get("duration_minutes"),
		DurationSeconds: r.PostForm.Get("duration_seconds"),
		RodeAt:          r.PostForm.Get("rode_at"),
		Timezone:        r.PostForm.Get("timezone"),
		TrailName:       strings.TrimSpace(r.PostForm.Get("trail_name")),
	}

	// need to implement this later
	media := json.RawMessage("null")

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 255), "title", "This field must be less than 255 characters")

	form.CheckField(validator.NotBlank(form.Description), "description", "This field cannot be blank")

	form.CheckField(validator.NotBlank(form.TrailName), "trail_name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.TrailName, 255), "trail_name", "This field must be less than 255 characters")

	form.CheckField(validator.NotBlank(form.Distance), "distance", "Distance is required")
	form.CheckField(validator.MaxChars(form.Distance, 8), "distance", "Distance is too long")
	form.CheckField(validator.ValidFloat(form.Distance, 0.01, 99999.99), "distance", "Distance is invalid")

	form.CheckField(validator.NotBlank(form.DurationHours), "duration_hours", "Duration hours is required")
	form.CheckField(validator.MaxChars(form.DurationHours, 2), "duration_hours", "Duration hours is too long")
	form.CheckField(validator.ValidInt(form.DurationHours, 0, 99), "duration_hours", "Duration hours is invalid")

	form.CheckField(validator.NotBlank(form.DurationMinutes), "duration_minutes", "Duration minutes is required")
	form.CheckField(validator.MaxChars(form.DurationMinutes, 2), "duration_minutes", "Duration minutes is too long")
	form.CheckField(validator.ValidInt(form.DurationMinutes, 0, 59), "duration_minutes", "Duration minutes is invalid")

	form.CheckField(validator.NotBlank(form.DurationSeconds), "duration_seconds", "Duration seconds is required")
	form.CheckField(validator.MaxChars(form.DurationSeconds, 2), "duration_seconds", "Duration seconds is too long")
	form.CheckField(validator.ValidInt(form.DurationSeconds, 0, 59), "duration_seconds", "Duration seconds is invalid")

	form.CheckField(validator.NotBlank(form.RodeAt), "rode_at", "Rode at is required")
	form.CheckField(validator.ValidTimeAndTimezone(form.RodeAt, form.Timezone), "rode_at", "Rode at is invalid")

	form.CheckField(validator.NotBlank(form.Timezone), "timezone", "Timezone is required")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	// convert rodeAt to time.Time utc
	loc, _ := time.LoadLocation(form.Timezone)
	rodeAtTime, _ := time.ParseInLocation("2006-01-02T15:04", form.RodeAt, loc)
	rodeAtUTC := rodeAtTime.UTC()
	distanceFloat, _ := strconv.ParseFloat(form.Distance, 64)
	h, _ := strconv.Atoi(form.DurationHours)
	m, _ := strconv.Atoi(form.DurationMinutes)
	s, _ := strconv.Atoi(form.DurationSeconds)
	convertedDuration := (h * 3600) + (m * 60) + s

	id, err := app.rides.Create(form.Title, form.Description, form.TrailName, distanceFloat, convertedDuration, rodeAtUTC, media)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ride/view/%d", id), http.StatusSeeOther)
}
