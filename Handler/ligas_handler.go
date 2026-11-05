package Handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
	"DeRabona/Logic"
)

type LigasHandler struct {
	Service *Logic.LigasService
}

// Create atiende POST /ligas: parsea el formulario, guarda la liga
// y redirige a la página principal (patrón Post/Redirect/Get).
func (h *LigasHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario inválido", http.StatusBadRequest)
		return
	}

	nombre := strings.TrimSpace(r.FormValue("nombre"))
	pais := strings.TrimSpace(r.FormValue("pais"))
	if nombre == "" || pais == "" {
		http.Error(w, "el nombre y el país de la liga son obligatorios", http.StatusBadRequest)
		return
	}

	arg := DataBase.CreateLigaParams{Nombre: nombre, Pais: pais}
	if _, err := h.Service.CreateLiga(r.Context(), arg); err != nil {
		log.Printf("crear liga: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *LigasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteLiga(r.Context(), int32(id)); err != nil {
		log.Printf("eliminar liga: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
