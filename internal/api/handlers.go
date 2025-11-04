package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KrlLevchenko/micomido/internal/models"
	"github.com/KrlLevchenko/micomido/internal/repository"
	"github.com/KrlLevchenko/micomido/internal/s3"
)

type API struct {
	MealRepo repository.MealRepository
	S3Client *s3.Client
}

func (a *API) GetMeals(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to *time.Time
	var err error

	if fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid 'from' date format. Use ISO 8601.")
			return
		}
		from = &t
	}

	if toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid 'to' date format. Use ISO 8601.")
			return
		}
		to = &t
	}

	if from != nil && to != nil && from.After(*to) {
		respondWithError(w, http.StatusBadRequest, "'from' date cannot be after 'to' date.")
		return
	}

	meals, err := a.MealRepo.GetMeals(r.Context(), from, to)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get meals")
		return
	}

	respondWithJSON(w, http.StatusOK, meals)
}

func (a *API) CreateMeal(w http.ResponseWriter, r *http.Request) {
	var newMeal models.Meal
	if err := json.NewDecoder(r.Body).Decode(&newMeal); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if newMeal.ID == "" {
		respondWithError(w, http.StatusBadRequest, "Meal ID is required")
		return
	}

	newMeal.At = time.Now().UTC()

	if err := a.MealRepo.CreateMeal(r.Context(), &newMeal); err != nil {
		// Here you might want to check for specific DB errors, like duplicate ID
		respondWithError(w, http.StatusInternalServerError, "Failed to create meal")
		return
	}

	respondWithJSON(w, http.StatusCreated, newMeal)
}

func (a *API) DeleteMeal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Meal ID is required")
		return
	}

	rowsAffected, err := a.MealRepo.DeleteMeal(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete meal")
		return
	}

	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Meal with the specified ID not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *API) AddMealPhoto(w http.ResponseWriter, r *http.Request) {
	mealID := r.PathValue("meal_id")
	photoID := r.PathValue("photo_id")

	if err := a.S3Client.Upload(r.Context(), photoID, r.Body); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to upload photo to S3")
		return
	}

	if err := a.MealRepo.AddMealPhoto(r.Context(), mealID, photoID); err != nil {
		// Attempt to clean up the uploaded S3 object if DB insertion fails
		_ = a.S3Client.Delete(r.Context(), photoID)
		respondWithError(w, http.StatusConflict, "Failed to create meal photo record") // 409 for duplicate
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) DeleteMealPhoto(w http.ResponseWriter, r *http.Request) {
	photoID := r.PathValue("photo_id")

	rowsAffected, err := a.MealRepo.DeleteMealPhoto(r.Context(), photoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete meal photo record")
		return
	}

	if rowsAffected == 0 {
		// No record in DB, so no file in S3 to delete.
		// Depending on desired behavior, you could still attempt a delete from S3.
		respondWithError(w, http.StatusNotFound, "Meal photo not found")
		return
	}

	if err := a.S3Client.Delete(r.Context(), photoID); err != nil {
		// The DB record was deleted, but S3 deletion failed.
		// This state is inconsistent and should be logged for manual cleanup.
		// For the client, the primary action (DB deletion) succeeded.
		// Log the error but still return success to the client.
		// log.Printf("error: failed to delete photo %s from S3: %v", photoID, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
