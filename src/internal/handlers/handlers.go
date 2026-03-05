package handlers

import (
	"APIforElasticBD/internal/types"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	limit int = 10 // sets nums of string in output
)

// Used for dependency injection in main and transfer it to handlers
type HandlerParam struct {
	store types.Store
	index string
	key   []byte
}

func InitHandlerParams(s types.Store, i string, k []byte) *HandlerParam {
	return &HandlerParam{store: s, index: i, key: k}
}

// Output token for autentification to use HandlerApiRecommend
func (h *HandlerParam) GetToken(w http.ResponseWriter, r *http.Request) {
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(h.key)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{"token": tokenStr}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode the token", http.StatusBadRequest)
	}
}

// A wrapper for other handlers for protection
func (h *HandlerParam) MiddleWareJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unathtorized. Authorization header required", http.StatusUnauthorized)
			return
		}
		authParts := strings.Fields(auth)
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			http.Error(w, "Unathtorized. Usage: Authorization Bearer <token>", http.StatusUnauthorized)
			return
		}
		tokenStr := authParts[1]
		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) { return h.key, nil })
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Default handler to check the data through a browser in HTML-format
func (h *HandlerParam) Handler(w http.ResponseWriter, r *http.Request) {
	strpage := r.URL.Query().Get("page")
	if strpage == "" {
		strpage = "0"
	}
	page, err := strconv.Atoi(strpage)
	if err != nil || page < 0 {
		http.Error(w, "Invalid 'page' value: 'foo'", http.StatusBadRequest)
		return
	}

	offset := page * limit

	places, total, err := h.store.GetPlaces(limit, offset, h.index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastPage := total / limit
	if x := total % limit; x == 0 {
		lastPage -= 1
	}

	if page > lastPage {
		http.Error(w, "Invalid 'page' value: 'foo'", http.StatusBadRequest)
		return
	}

	tmp, err := template.ParseFiles("../../templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmp.Execute(w, map[string]interface{}{
		"Total":  total,
		"Places": places,
		"Page":   page,
		"First":  0,
		"Prev":   page - 1,
		"Next":   page + 1,
		"Last":   lastPage,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Default handler to check the data through a browser in JSON-format
func (h *HandlerParam) HandlerApi(w http.ResponseWriter, r *http.Request) {
	strpage := r.URL.Query().Get("page")
	if strpage == "" {
		strpage = "0"
	}
	page, err := strconv.Atoi(strpage)
	if err != nil || page < 0 {
		http.Error(w, "Invalid 'page' value: 'foo'", http.StatusBadRequest)
		return
	}

	offset := page * limit

	places, total, err := h.store.GetPlaces(limit, offset, h.index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastPage := total / limit
	if x := total % limit; x == 0 {
		lastPage -= 1
	}

	if page > lastPage {
		http.Error(w, "Invalid 'page' value: 'foo'", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"Total":  total,
		"Places": places,
		"Page":   page,
		"First":  0,
		"Prev":   page - 1,
		"Next":   page + 1,
		"Last":   lastPage,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler for using param lat & lon to see the three closest places in JSON-format
func (h *HandlerParam) HandlerApiRecommend(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		http.Error(w, "lat parsing error", http.StatusBadRequest)
	} else if lat > 90.0 || lat < -90.0 {
		http.Error(w, "wrong latitude: -90.0 < lat < 90.0", http.StatusBadRequest)
	}
	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil {
		http.Error(w, "lon parsing error", http.StatusBadRequest)
	} else if lon > 180.0 || lon < -180.0 {
		http.Error(w, "wrong longitude: -180.0 < lon < 180.0", http.StatusBadRequest)
	}

	places, err := h.store.GetClosestPlaces(lat, lon, h.index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"Name":   "Recommendation",
		"Places": places,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
