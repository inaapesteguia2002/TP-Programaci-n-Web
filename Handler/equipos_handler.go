package Handler

import (
	"net/http"
	"strconv"
	"strings"

	"DeRabona/Logic"
	"DeRabona/views"
)

type EquiposHandler struct {
	Service *Logic.EquiposService
}

func (h *EquiposHandler) HandleEquipos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		equipos, err := h.Service.ListEquipos(ctx)
		if err != nil {
			http.Error(w, "Error al obtener equipos", http.StatusInternalServerError)
			return
		}
		// Renderiza la página completa usando el layout y componentes templ
		views.IndexEquipos(equipos).Render(ctx, w)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulario inválido", http.StatusBadRequest)
			return
		}
		nombre := r.FormValue("nombre")

		// CORREGIDO: Se pasa directamente el string 'nombre' acorde a EquiposService.CreateEquipo
		_, err := h.Service.CreateEquipo(ctx, nombre)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		equipos, err := h.Service.ListEquipos(ctx)
		if err != nil {
			http.Error(w, "Error al listar equipos", http.StatusInternalServerError)
			return
		}

		// Si la petición viene de HTMX, devolvemos solo el fragmento de la lista
		if r.Header.Get("HX-Request") == "true" {
			views.EquipoList(equipos).Render(ctx, w)
			return
		}

		// Si es una petición clásica, aplicamos Post/Redirect/Get (PRG)
		http.Redirect(w, r, "/equipos", http.StatusSeeOther)

	case http.MethodDelete:
		// Espera una ruta tipo /equipos/{id}
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

		err = h.Service.DeleteEquipo(ctx, int32(id))
		if err != nil {
			http.Error(w, "Error al eliminar equipo", http.StatusInternalServerError)
			return
		}

		// Respuesta vacía con 200 OK para que HTMX remueva el elemento con outerHTML
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
