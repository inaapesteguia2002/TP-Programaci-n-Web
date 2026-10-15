package Handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
	"DeRabona/Logic"
)

// LigasHandler se encarga de las peticiones HTTP e inyecta el Service
type LigasHandler struct {
	Service *Logic.LigasService
}

// HandleLigas maneja GET /ligas y POST /ligas
func (h *LigasHandler) HandleLigas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ligas, err := h.Service.ListLigas(r.Context())
		if err != nil {
			http.Error(w, "Error al obtener las ligas", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ligas)

	case http.MethodPost:
		var req struct {
			Nombre string `json:"nombre"`
			Pais   string `json:"pais"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		arg := DataBase.CreateLigaParams{
			Nombre: req.Nombre,
			Pais:   req.Pais,
		}

		nuevaLiga, err := h.Service.CreateLiga(r.Context(), arg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Cae acá si el Service tira error de validación
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201 Created
		json.NewEncoder(w).Encode(nuevaLiga)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandleLigaByID maneja GET, PUT y DELETE en /ligas/{id}
func (h *LigasHandler) HandleLigaByID(w http.ResponseWriter, r *http.Request) {
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

	ligaID := int32(id)

	switch r.Method {
	case http.MethodGet:
		liga, err := h.Service.GetLiga(r.Context(), ligaID)
		if err != nil {
			http.Error(w, "Liga no encontrada", http.StatusNotFound) // 404 Not Found
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(liga)

	case http.MethodPut:
		var req struct {
			Nombre string `json:"nombre"`
			Pais   string `json:"pais"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		arg := DataBase.UpdateLigaParams{
			ID:     ligaID,
			Nombre: req.Nombre,
			Pais:   req.Pais,
		}

		if err := h.Service.UpdateLiga(r.Context(), arg); err != nil {
			http.Error(w, "Error al actualizar o liga no encontrada", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"mensaje": "Liga actualizada con éxito"})

	case http.MethodDelete:
		if err := h.Service.DeleteLiga(r.Context(), ligaID); err != nil {
			http.Error(w, "Error al eliminar o liga no encontrada", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
