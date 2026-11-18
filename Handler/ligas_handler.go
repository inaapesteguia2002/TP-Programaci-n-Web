package Handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
	"DeRabona/Logic"
	"DeRabona/Views"
)

type LigasHandler struct {
	Service *Logic.LigasService
}

// Create atiende POST /ligas: crea la liga, consulta la lista actualizada
// y devuelve el fragmento HTML de la lista de ligas (más los selectores OOB) para HTMX.
func (h *LigasHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
	if _, err := h.Service.CreateLiga(ctx, arg); err != nil {
		log.Printf("crear liga: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ligas, err := h.Service.ListLigas(ctx)
	if err != nil {
		log.Printf("listar ligas tras crear: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := Views.LigaListConOOB(ligas).Render(ctx, w); err != nil {
		log.Printf("renderizar LigaList: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Delete atiende DELETE /ligas/{id}: obtiene el ID de la URL, elimina la liga
// y devuelve una respuesta vacía con 200 OK para que HTMX quite la fila del DOM.
func (h *LigasHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		idStr = r.FormValue("id")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteLiga(r.Context(), int32(id)); err != nil {
		log.Printf("eliminar liga: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
