package Handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
	"DeRabona/Logic"
)

// PartidosHandler se encarga de las peticiones HTTP e inyecta el Service
type PartidosHandler struct {
	Service *Logic.PartidosService
}

// HandlePartidos maneja GET /partidos y POST /partidos
func (h *PartidosHandler) HandlePartidos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		partidos, err := h.Service.ListPartidos(r.Context())
		if err != nil {
			http.Error(w, "Error al obtener los partidos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(partidos)

	case http.MethodPost:
		var req struct {
			LigaID            int32  `json:"liga_id"`
			EquipoLocalID     int32  `json:"equipo_local_id"`
			EquipoVisitanteID int32  `json:"equipo_visitante_id"`
			GolesLocal        int32  `json:"goles_local"`
			GolesVisitante    int32  `json:"goles_visitante"`
			Estado            string `json:"estado"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		arg := DataBase.CreatePartidoParams{
			LigaID:            req.LigaID,
			EquipoLocalID:     req.EquipoLocalID,
			EquipoVisitanteID: req.EquipoVisitanteID,
			GolesLocal:        sql.NullInt32{Int32: req.GolesLocal, Valid: true},
			GolesVisitante:    sql.NullInt32{Int32: req.GolesVisitante, Valid: true},
			Estado:            sql.NullString{String: req.Estado, Valid: req.Estado != ""},
		}

		nuevoPartido, err := h.Service.CreatePartido(r.Context(), arg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201 Created
		json.NewEncoder(w).Encode(nuevoPartido)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandlePartidoByID maneja GET, PUT y DELETE en /partidos/{id}
func (h *PartidosHandler) HandlePartidoByID(w http.ResponseWriter, r *http.Request) {
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
		partido, err := h.Service.GetPartido(r.Context(), partidoID)
		if err != nil {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(partido)

	case http.MethodPut:
		var req struct {
			GolesLocal     int32  `json:"goles_local"`
			GolesVisitante int32  `json:"goles_visitante"`
			Estado         string `json:"estado"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		arg := DataBase.UpdatePartidoParams{
			ID:             partidoID,
			GolesLocal:     sql.NullInt32{Int32: req.GolesLocal, Valid: true},
			GolesVisitante: sql.NullInt32{Int32: req.GolesVisitante, Valid: true},
			Estado:         sql.NullString{String: req.Estado, Valid: req.Estado != ""},
		}

		if err := h.Service.UpdatePartido(r.Context(), arg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"mensaje": "Partido actualizado con éxito"})

	case http.MethodDelete:
		if err := h.Service.DeletePartido(r.Context(), partidoID); err != nil {
			http.Error(w, "Error al eliminar o partido no encontrado", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
