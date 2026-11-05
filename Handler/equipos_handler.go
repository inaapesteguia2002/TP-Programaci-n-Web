package Handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/Logic"
)

type EquiposHandler struct {
	Service *Logic.EquiposService
}

// Create atiende POST /equipos: parsea el formulario, guarda el equipo
// y redirige a la página principal (patrón Post/Redirect/Get).
func (h *EquiposHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario inválido", http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(r.FormValue("nombre"))
	if nombre == "" {
		http.Error(w, "el nombre del equipo es obligatorio", http.StatusBadRequest)
		return
	}

	if _, err := h.Service.CreateEquipo(r.Context(), nombre); err != nil {
		log.Printf("crear equipo: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *EquiposHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteEquipo(r.Context(), int32(id)); err != nil {
		log.Printf("eliminar equipo: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
