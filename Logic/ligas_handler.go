package Logic

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
)

// DTO para validar las peticiones JSON entrantes de Ligas
type LigaRequest struct {
	Nombre string `json:"nombre"`
	Pais   string `json:"pais"`
}

// ------------------------------------------------------------------
// RUTAS GENERALES: /ligas
// ------------------------------------------------------------------

func (app *APIHandler) LigasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getLigas(w, r)
	case http.MethodPost:
		app.createLiga(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *APIHandler) getLigas(w http.ResponseWriter, r *http.Request) {
	ligas, err := app.Repo.ListLigas(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener ligas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ligas)
}

func (app *APIHandler) createLiga(w http.ResponseWriter, r *http.Request) {
	var req LigaRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Nombre) == "" || strings.TrimSpace(req.Pais) == "" {
		http.Error(w, "El nombre y el país no pueden estar vacíos", http.StatusBadRequest)
		return
	}

	// Mapeamos nuestro DTO a los parámetros requeridos por sqlc
	params := DataBase.CreateLigaParams{
		Nombre: req.Nombre,
		Pais:   req.Pais,
	}

	nuevaLiga, err := app.Repo.CreateLiga(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(nuevaLiga)
}

// ------------------------------------------------------------------
// RUTAS ESPECÍFICAS POR ID: /ligas/{id}
// ------------------------------------------------------------------

func (app *APIHandler) LigaHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "El ID debe ser numérico", http.StatusBadRequest)
		return
	}

	// Castear a int32 para que coincida con el SERIAL de Postgres
	ligaID := int32(id)

	switch r.Method {
	case http.MethodGet:
		app.getLiga(w, r, ligaID)
	case http.MethodPut:
		app.updateLiga(w, r, ligaID)
	case http.MethodDelete:
		app.deleteLiga(w, r, ligaID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *APIHandler) getLiga(w http.ResponseWriter, r *http.Request, id int32) {
	liga, err := app.Repo.GetLiga(r.Context(), id)
	if err != nil {
		http.Error(w, "Liga no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(liga)
}

func (app *APIHandler) deleteLiga(w http.ResponseWriter, r *http.Request, id int32) {
	err := app.Repo.DeleteLiga(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al eliminar la liga", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (app *APIHandler) updateLiga(w http.ResponseWriter, r *http.Request, id int32) {
	var req LigaRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Construimos los parámetros requeridos por sqlc para el UPDATE
	params := DataBase.UpdateLigaParams{
		ID:     id,
		Nombre: req.Nombre,
		Pais:   req.Pais,
	}

	err := app.Repo.UpdateLiga(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(params)
}
