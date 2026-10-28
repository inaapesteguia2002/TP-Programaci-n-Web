package Handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
	"DeRabona/Logic"
)

// EquiposHandler se encarga de las peticiones HTTP e inyecta el Service
type EquiposHandler struct {
	Service *Logic.EquiposService
}

// HandleEquipos maneja la colección entera: GET /equipos y POST /equipos
func (h *EquiposHandler) HandleEquipos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		equipos, err := h.Service.ListEquipos(r.Context())
		if err != nil {
			http.Error(w, "Error al obtener los equipos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(equipos)

	case http.MethodPost:
		var req struct {
			Nombre string `json:"nombre"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		nuevoEquipo, err := h.Service.CreateEquipo(r.Context(), req.Nombre)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Error de validación del Service
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // Estado 201: Creado
		json.NewEncoder(w).Encode(nuevoEquipo)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandleEquipoByID maneja recursos individuales: GET, PUT, DELETE en /equipos/{id}
func (h *EquiposHandler) HandleEquipoByID(w http.ResponseWriter, r *http.Request) {
	// Extraemos el ID de la URL (ej: /equipos/1)
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

	equipoID := int32(id)

	switch r.Method {
	case http.MethodGet:
		equipo, err := h.Service.GetEquipo(r.Context(), equipoID)
		if err != nil {
			http.Error(w, "Equipo no encontrado", http.StatusNotFound) // Estado 404
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(equipo)

	case http.MethodPut:
		var req struct {
			Nombre string `json:"nombre"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		arg := DataBase.UpdateEquipoParams{
			ID:     equipoID,
			Nombre: req.Nombre,
		}

		if err := h.Service.UpdateEquipo(r.Context(), arg); err != nil {
			http.Error(w, "Error al actualizar o equipo no encontrado", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     arg.ID,
			"nombre": arg.Nombre,
		})

	case http.MethodDelete:
		if err := h.Service.DeleteEquipo(r.Context(), equipoID); err != nil {
			http.Error(w, "Error al eliminar o equipo no encontrado", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // Estado 204: Sin contenido (Eliminación exitosa)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
