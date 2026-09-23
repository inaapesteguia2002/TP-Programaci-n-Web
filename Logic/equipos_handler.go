package Logic

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
)

// 1. Estructura principal que recibe el repositorio de sqlc
type APIHandler struct {
	Repo *DataBase.Queries
}

// DTO para validar las peticiones JSON entrantes
type EquipoRequest struct {
	Nombre string `json:"nombre"`
}

// ------------------------------------------------------------------
// RUTAS GENERALES: /equipos
// ------------------------------------------------------------------

func (app *APIHandler) EquiposHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getEquipos(w, r)
	case http.MethodPost:
		app.createEquipo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *APIHandler) getEquipos(w http.ResponseWriter, r *http.Request) {
	equipos, err := app.Repo.ListEquipos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipos)
}

func (app *APIHandler) createEquipo(w http.ResponseWriter, r *http.Request) {
	var req EquipoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Nombre) == "" {
		http.Error(w, "El nombre del equipo no puede estar vacío", http.StatusBadRequest)
		return
	}

	// sqlc pide un string directamente al haber un solo parámetro
	nuevoEquipo, err := app.Repo.CreateEquipo(r.Context(), req.Nombre)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(nuevoEquipo)
}

// ------------------------------------------------------------------
// RUTAS ESPECÍFICAS POR ID: /equipos/{id}
// ------------------------------------------------------------------

func (app *APIHandler) EquipoHandler(w http.ResponseWriter, r *http.Request) {
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
	equipoID := int32(id)

	switch r.Method {
	case http.MethodGet:
		app.getEquipo(w, r, equipoID)
	case http.MethodPut:
		app.updateEquipo(w, r, equipoID)
	case http.MethodDelete:
		app.deleteEquipo(w, r, equipoID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *APIHandler) getEquipo(w http.ResponseWriter, r *http.Request, id int32) {
	equipo, err := app.Repo.GetEquipo(r.Context(), id)
	if err != nil {
		http.Error(w, "Equipo no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipo)
}

func (app *APIHandler) deleteEquipo(w http.ResponseWriter, r *http.Request, id int32) {
	err := app.Repo.DeleteEquipo(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al eliminar el equipo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (app *APIHandler) updateEquipo(w http.ResponseWriter, r *http.Request, id int32) {
	var req EquipoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Construimos los parámetros requeridos por sqlc para el UPDATE
	params := DataBase.UpdateEquipoParams{
		ID:     id,
		Nombre: req.Nombre,
	}

	err = app.Repo.UpdateEquipo(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(params)
}
