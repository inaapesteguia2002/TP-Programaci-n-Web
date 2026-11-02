package Handler

import (
	"net/http"
	"strconv"

	"DeRabona/DataBase"
	"DeRabona/Logic"
	"DeRabona/views"
)

type TablaHandler struct {
	TablaService  *Logic.TablaService
	LigaService   *Logic.LigasService
	EquipoService *Logic.EquiposService
}

func (h *TablaHandler) HandleTabla(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		// Lee el parámetro ?liga_id= de la URL (por defecto 1 si no viene especificado)
		ligaIDStr := r.URL.Query().Get("liga_id")
		var ligaID int32 = 1
		if ligaIDStr != "" {
			parsed, err := strconv.Atoi(ligaIDStr)
			if err == nil {
				ligaID = int32(parsed)
			}
		}

		filas, err := h.TablaService.ObtenerTabla(ctx, ligaID)
		if err != nil {
			http.Error(w, "Error al obtener la tabla de posiciones", http.StatusInternalServerError)
			return
		}

		ligas, _ := h.LigaService.ListLigas(ctx)
		equipos, _ := h.EquipoService.ListEquipos(ctx)

		views.IndexTabla(ligaID, filas, ligas, equipos).Render(ctx, w)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulario inválido", http.StatusBadRequest)
			return
		}

		ligaID, _ := strconv.Atoi(r.FormValue("liga_id"))
		equipoID, _ := strconv.Atoi(r.FormValue("equipo_id"))

		arg := DataBase.CreateTablaParams{
			LigaID:   int32(ligaID),
			EquipoID: int32(equipoID),
		}

		_, err := h.TablaService.RegistrarEquipo(ctx, arg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filas, err := h.TablaService.ObtenerTabla(ctx, int32(ligaID))
		if err != nil {
			http.Error(w, "Error al actualizar la tabla", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			views.TablaList(filas).Render(ctx, w)
			return
		}

		http.Redirect(w, r, "/tabla?liga_id="+strconv.Itoa(ligaID), http.StatusSeeOther)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
