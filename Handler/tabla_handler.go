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

// TablaHandler se encarga de las peticiones HTTP e inyecta el Service
type TablaHandler struct {
	Service *Logic.TablaService
}

// HandleTabla maneja POST, PUT y DELETE en /tabla
func (h *TablaHandler) HandleTabla(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req struct {
			LigaID   int32 `json:"liga_id"`
			EquipoID int32 `json:"equipo_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		arg := DataBase.CreateTablaParams{
			LigaID:   req.LigaID,
			EquipoID: req.EquipoID,
		}

		nuevaEntrada, err := h.Service.RegistrarEquipo(r.Context(), arg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201 Created
		json.NewEncoder(w).Encode(nuevaEntrada)

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

		arg := DataBase.UpdateTablaParams{
			LigaID:            req.LigaID,
			EquipoID:          req.EquipoID,
			Puntos:            sql.NullInt32{Int32: req.Puntos, Valid: true},
			PartidosGanados:   sql.NullInt32{Int32: req.PartidosGanados, Valid: true},
			PartidosEmpatados: sql.NullInt32{Int32: req.PartidosEmpatados, Valid: true},
			PartidosPerdidos:  sql.NullInt32{Int32: req.PartidosPerdidos, Valid: true},
		}

		if err := h.Service.ActualizarEstadisticas(r.Context(), arg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"liga_id":            arg.LigaID,
			"equipo_id":          arg.EquipoID,
			"puntos":             arg.Puntos.Int32,
			"partidos_ganados":   arg.PartidosGanados.Int32,
			"partidos_empatados": arg.PartidosEmpatados.Int32,
			"partidos_perdidos":  arg.PartidosPerdidos.Int32,
		})

	case http.MethodDelete:
		var req struct {
			LigaID   int32 `json:"liga_id"`
			EquipoID int32 `json:"equipo_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		arg := DataBase.DeleteEquipoDeTablaParams{
			LigaID:   req.LigaID,
			EquipoID: req.EquipoID,
		}

		if err := h.Service.RemoverEquipo(r.Context(), arg); err != nil {
			http.Error(w, "Error al eliminar el equipo de la tabla", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandleTablaByLiga maneja el GET en /tabla/{liga_id}
func (h *TablaHandler) HandleTablaByLiga(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	ligaID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "El ID de la liga debe ser numérico", http.StatusBadRequest)
		return
	}

	tabla, err := h.Service.ObtenerTabla(r.Context(), int32(ligaID))
	if err != nil {
		http.Error(w, "Error al obtener la tabla de posiciones", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tabla)
}
