package Logic

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
)

// DTO para registrar un equipo en la tabla de una liga
type CreateTablaRequest struct {
	LigaID   int32 `json:"liga_id"`
	EquipoID int32 `json:"equipo_id"`
}

// DTO para actualizar estadísticas sumando valores
type UpdateTablaRequest struct {
	Puntos            int32 `json:"puntos"`
	PartidosGanados   int32 `json:"partidos_ganados"`
	PartidosEmpatados int32 `json:"partidos_empatados"`
	PartidosPerdidos  int32 `json:"partidos_perdidos"`
}

// ------------------------------------------------------------------
// RUTAS GENERALES: /tabla
// ------------------------------------------------------------------

func (app *APIHandler) TablasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.createTabla(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *APIHandler) createTabla(w http.ResponseWriter, r *http.Request) {
	var req CreateTablaRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.LigaID <= 0 || req.EquipoID <= 0 {
		http.Error(w, "Los IDs de liga y equipo son obligatorios", http.StatusBadRequest)
		return
	}

	params := DataBase.CreateTablaParams{
		LigaID:   req.LigaID,
		EquipoID: req.EquipoID,
	}

	nuevaPosicion, err := app.Repo.CreateTabla(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al registrar el equipo en la tabla", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevaPosicion)
}

// ------------------------------------------------------------------
// RUTA PARA LISTAR POSICIONES DE UNA LIGA: /tabla/{liga_id}
// ------------------------------------------------------------------

func (app *APIHandler) TablaLigaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "El ID de la liga debe ser numérico", http.StatusBadRequest)
		return
	}

	ligaID := int32(id)

	// Llama al ListTabla que hace el JOIN con equipos y ordena por puntos
	posiciones, err := app.Repo.ListTabla(r.Context(), ligaID)
	if err != nil {
		http.Error(w, "Error al obtener la tabla de posiciones", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posiciones)
}

// ------------------------------------------------------------------
// RUTA PARA MODIFICAR O ELIMINAR: /tabla/equipo
// ------------------------------------------------------------------

func (app *APIHandler) TablaAccionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		var req struct {
			LigaID            int32 `json:"liga_id"`
			EquipoID          int32 `json:"equipo_id"`
			Puntos            int32 `json:"puntos"`
			PartidosGanados   int32 `json:"partidos_ganados"`
			PartidosEmpatados int32 `json:"partidos_empatados"`
			PartidosPerdidos  int32 `json:"partidos_perdidos"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		params := DataBase.UpdateTablaParams{
			LigaID:            req.LigaID,
			EquipoID:          req.EquipoID,
			Puntos:            sql.NullInt32{Int32: req.Puntos, Valid: true},
			PartidosGanados:   sql.NullInt32{Int32: req.PartidosGanados, Valid: true},
			PartidosEmpatados: sql.NullInt32{Int32: req.PartidosEmpatados, Valid: true},
			PartidosPerdidos:  sql.NullInt32{Int32: req.PartidosPerdidos, Valid: true},
		}

		err := app.Repo.UpdateTabla(r.Context(), params)
		if err != nil {
			http.Error(w, "Error al actualizar la tabla", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "actualizado"})

	case http.MethodDelete:
		var req struct {
			LigaID   int32 `json:"liga_id"`
			EquipoID int32 `json:"equipo_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		params := DataBase.DeleteEquipoDeTablaParams{
			LigaID:   req.LigaID,
			EquipoID: req.EquipoID,
		}

		err := app.Repo.DeleteEquipoDeTabla(r.Context(), params)
		if err != nil {
			http.Error(w, "Error al eliminar de la tabla", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
