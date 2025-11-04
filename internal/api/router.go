package api

import "net/http"

func NewRouter(api *API) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello this is MiComido proj"))
	})

	mux.HandleFunc("GET /api/meals", api.GetMeals)
	mux.HandleFunc("POST /api/meal", api.CreateMeal)
	mux.HandleFunc("DELETE /api/meal/{id}", api.DeleteMeal)
	mux.HandleFunc("POST /api/meal/{meal_id}/photo/{photo_id}", api.AddMealPhoto)
	mux.HandleFunc("DELETE /api/meal/{meal_id}/photo/{photo_id}", api.DeleteMealPhoto)

	return mux
}
