package Handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/Logic"
	"DeRabona/Views"
)

type EquiposHandler struct {
	Service *Logic.EquiposService
}

// Create atiende POST /equipos: crea el equipo, consulta la lista actualizada
// y devuelve el fragmento HTML de la lista de equipos (más los selectores OOB) para HTMX.
func (h *EquiposHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario inválido", http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(r.FormValue("nombre"))
	if nombre == "" {
		http.Error(w, "el nombre del equipo es obligatorio", http.StatusBadRequest)
		return
	}

	if _, err := h.Service.CreateEquipo(ctx, nombre); err != nil {
		log.Printf("crear equipo: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	equipos, err := h.Service.ListEquipos(ctx)
	if err != nil {
		log.Printf("listar equipos tras crear: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := Views.EquipoListConOOB(equipos).Render(ctx, w); err != nil {
		log.Printf("renderizar EquipoList: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Delete atiende DELETE /equipos/{id}: obtiene el ID de la URL, elimina el equipo
// y devuelve una respuesta vacía con 200 OK para que HTMX quite la fila del DOM.
func (h *EquiposHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		idStr = r.FormValue("id")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteEquipo(ctx, int32(id)); err != nil {
		log.Printf("eliminar equipo: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		equipos, err := h.Service.ListEquipos(ctx)
		if err != nil {
			log.Printf("listar equipos tras eliminar: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if err := Views.SelectsEquiposOOB(equipos).Render(ctx, w); err != nil {
			log.Printf("renderizar SelectsEquiposOOB: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
