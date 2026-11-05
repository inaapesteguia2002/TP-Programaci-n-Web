package Handler

import (
	"log"
	"net/http"
	"strconv"

	"DeRabona/DataBase"
	"DeRabona/Logic"
	"DeRabona/Views"
)

// HomeHandler atiende la página principal: junta los datos de todas las entidades
type HomeHandler struct {
	Ligas    *Logic.LigasService
	Equipos  *Logic.EquiposService
	Partidos *Logic.PartidosService
	Tabla    *Logic.TablaService
}

// Index atiende GET /: obtiene todos los datos y renderiza la página principal.
// La liga cuya tabla se muestra se elige con /?liga_id=N
func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ligas, err := h.Ligas.ListLigas(ctx)
	if err != nil {
		log.Printf("listar ligas: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	equipos, err := h.Equipos.ListEquipos(ctx)
	if err != nil {
		log.Printf("listar equipos: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	partidos, err := h.Partidos.ListPartidos(ctx)
	if err != nil {
		log.Printf("listar partidos: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Liga seleccionada para la tabla de posiciones (0 = ninguna)
	var ligaSeleccionada int32
	var tabla []DataBase.GetTablaPorLigaRow
	if v := r.URL.Query().Get("liga_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil || id <= 0 {
			http.Error(w, "liga_id inválido", http.StatusBadRequest)
			return
		}
		ligaSeleccionada = int32(id)

		tabla, err = h.Tabla.ObtenerTabla(ctx, ligaSeleccionada)
		if err != nil {
			log.Printf("obtener tabla de liga %d: %v", ligaSeleccionada, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	if err := Views.IndexPage(ligas, equipos, partidos, tabla, ligaSeleccionada).Render(ctx, w); err != nil {
		log.Printf("renderizar index: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
