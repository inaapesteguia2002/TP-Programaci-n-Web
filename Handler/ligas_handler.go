package Handler

import (
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
	"DeRabona/Logic"
	"DeRabona/views"
)

type LigasHandler struct {
	Service *Logic.LigasService
}

func (h *LigasHandler) HandleLigas(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		ligas, err := h.Service.ListLigas(ctx)
		if err != nil {
			http.Error(w, "Error al obtener ligas", http.StatusInternalServerError)
			return
		}
		views.IndexLigas(ligas).Render(ctx, w)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulario inválido", http.StatusBadRequest)
			return
		}
		nombre := r.FormValue("nombre")
		pais := r.FormValue("pais")

		_, err := h.Service.CreateLiga(ctx, DataBase.CreateLigaParams{Nombre: nombre, Pais: pais})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ligas, err := h.Service.ListLigas(ctx)
		if err != nil {
			http.Error(w, "Error al listar ligas", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			views.LigaList(ligas).Render(ctx, w)
			return
		}

		http.Redirect(w, r, "/ligas", http.StatusSeeOther)

	case http.MethodDelete:
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "URL inválida", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		err = h.Service.DeleteLiga(ctx, int32(id))
		if err != nil {
			http.Error(w, "Error al eliminar liga", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
