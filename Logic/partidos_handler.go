package Logic

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
)

// DTO para la creación de un partido
type CreatePartidoRequest struct {
	LigaID            int32  `json:"liga_id"`
	EquipoLocalID     int32  `json:"equipo_local_id"`
	EquipoVisitanteID int32  `json:"equipo_visitante_id"`
	GolesLocal        int32  `json:"goles_local"`
	GolesVisitante    int32  `json:"goles_visitante"`
	Estado            string `json:"estado"`
}

// DTO para la actualización de un partido (sólo resultado y estado)
type UpdatePartidoRequest struct {
	GolesLocal     int32  `json:"goles_local"`
	GolesVisitante int32  `json:"goles_visitante"`
	Estado         string `json:"estado"`
}

// ------------------------------------------------------------------
// RUTAS GENERALES: /partidos
// ------------------------------------------------------------------

func (app *APIHandler) PartidosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getPartidos(w, r)
	case http.MethodPost:
		app.createPartido(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *APIHandler) getPartidos(w http.ResponseWriter, r *http.Request) {
	partidos, err := app.Repo.ListPartidos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener partidos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partidos)
}

func (app *APIHandler) createPartido(w http.ResponseWriter, r *http.Request) {
	var req CreatePartidoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validaciones básicas de negocio
	if req.LigaID <= 0 || req.EquipoLocalID <= 0 || req.EquipoVisitanteID <= 0 {
		http.Error(w, "Los IDs de liga y equipos son obligatorios", http.StatusBadRequest)
		return
	}
	if req.Estado == "" {
		req.Estado = "pendiente" // Valor por defecto definido en schema.sql
	}

	// Mapeo de los tipos simples de Go (int32, string) a los tipos sql.Null... que generó sqlc
	params := DataBase.CreatePartidoParams{
		LigaID:            req.LigaID,
		EquipoLocalID:     req.EquipoLocalID,
		EquipoVisitanteID: req.EquipoVisitanteID,
		GolesLocal:        sql.NullInt32{Int32: req.GolesLocal, Valid: true},
		GolesVisitante:    sql.NullInt32{Int32: req.GolesVisitante, Valid: true},
		Estado:            sql.NullString{String: req.Estado, Valid: true},
	}

	nuevoPartido, err := app.Repo.CreatePartido(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(nuevoPartido)
}

// ------------------------------------------------------------------
// RUTAS ESPECÍFICAS POR ID: /partidos/{id}
// ------------------------------------------------------------------

func (app *APIHandler) PartidoHandler(w http.ResponseWriter, r *http.Request) {
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

	partidoID := int32(id)

	switch r.Method {
	case http.MethodGet:
		app.getPartido(w, r, partidoID)
	case http.MethodPut:
		app.updatePartido(w, r, partidoID)
	case http.MethodDelete:
		app.deletePartido(w, r, partidoID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *APIHandler) getPartido(w http.ResponseWriter, r *http.Request, id int32) {
	partido, err := app.Repo.GetPartido(r.Context(), id)
	if err != nil {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partido)
}

func (app *APIHandler) deletePartido(w http.ResponseWriter, r *http.Request, id int32) {
	err := app.Repo.DeletePartido(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al eliminar el partido", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (app *APIHandler) updatePartido(w http.ResponseWriter, r *http.Request, id int32) {
	var req UpdatePartidoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.Estado == "" {
		req.Estado = "pendiente"
	}

	// Mapeo hacia sqlc utilizando sql.NullInt32 y sql.NullString para la actualización
	params := DataBase.UpdatePartidoParams{
		ID:             id,
		GolesLocal:     sql.NullInt32{Int32: req.GolesLocal, Valid: true},
		GolesVisitante: sql.NullInt32{Int32: req.GolesVisitante, Valid: true},
		Estado:         sql.NullString{String: req.Estado, Valid: true},
	}

	err := app.Repo.UpdatePartido(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(params)
}
