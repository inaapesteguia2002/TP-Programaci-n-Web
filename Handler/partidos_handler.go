package Handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"DeRabona/DataBase"
	"DeRabona/Logic"
	"DeRabona/views"
)

type PartidosHandler struct {
	PartidoService *Logic.PartidosService
	LigaService    *Logic.LigasService
	EquipoService  *Logic.EquiposService
}

func (h *PartidosHandler) HandlePartidos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		partidos, err := h.PartidoService.ListPartidos(ctx)
		if err != nil {
			http.Error(w, "Error al listar partidos", http.StatusInternalServerError)
			return
		}
		ligas, _ := h.LigaService.ListLigas(ctx)
		equipos, _ := h.EquipoService.ListEquipos(ctx)

		views.IndexPartidos(partidos, ligas, equipos).Render(ctx, w)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulario inválido", http.StatusBadRequest)
			return
		}

		ligaID, _ := strconv.Atoi(r.FormValue("liga_id"))
		localID, _ := strconv.Atoi(r.FormValue("equipo_local_id"))
		visitanteID, _ := strconv.Atoi(r.FormValue("equipo_visitante_id"))
		golesLocal, _ := strconv.Atoi(r.FormValue("goles_local"))
		golesVisitante, _ := strconv.Atoi(r.FormValue("goles_visitante"))
		estado := r.FormValue("estado")

		arg := DataBase.CreatePartidoParams{
			LigaID:            int32(ligaID),
			EquipoLocalID:     int32(localID),
			EquipoVisitanteID: int32(visitanteID),
			GolesLocal:        sql.NullInt32{Int32: int32(golesLocal), Valid: true},
			GolesVisitante:    sql.NullInt32{Int32: int32(golesVisitante), Valid: true},
			Estado:            sql.NullString{String: estado, Valid: estado != ""},
		}

		_, err := h.PartidoService.CreatePartido(ctx, arg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		partidos, err := h.PartidoService.ListPartidos(ctx)
		if err != nil {
			http.Error(w, "Error al listar partidos", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			views.PartidoList(partidos).Render(ctx, w)
			return
		}

		http.Redirect(w, r, "/partidos", http.StatusSeeOther)

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

		err = h.PartidoService.DeletePartido(ctx, int32(id))
		if err != nil {
			http.Error(w, "Error al eliminar partido", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
