package Handler

import (
	"log"
	"net/http"
	"strconv"

	"DeRabona/DataBase"
	"DeRabona/Logic"
	"DeRabona/Views"
)

type TablaHandler struct {
	Service *Logic.TablaService
	Ligas   *Logic.LigasService
}

// Create atiende POST /tabla: inscribe un equipo en una liga, consulta la tabla actualizada
// de esa liga y devuelve el fragmento TablaList para que HTMX lo muestre inmediatamente.
func (h *TablaHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario inválido", http.StatusBadRequest)
		return
	}

	ligaID, err1 := strconv.Atoi(r.FormValue("liga_id"))
	equipoID, err2 := strconv.Atoi(r.FormValue("equipo_id"))
	if err1 != nil || err2 != nil || ligaID <= 0 || equipoID <= 0 {
		http.Error(w, "liga y equipo son obligatorios", http.StatusBadRequest)
		return
	}

	arg := DataBase.CreateTablaParams{
		LigaID:   int32(ligaID),
		EquipoID: int32(equipoID),
	}
	if _, err := h.Service.RegistrarEquipo(ctx, arg); err != nil {
		log.Printf("inscribir equipo en liga: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tabla, err := h.Service.ObtenerTabla(ctx, int32(ligaID))
	if err != nil {
		log.Printf("obtener tabla tras inscribir: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var ligas []DataBase.Liga
	if h.Ligas != nil {
		ligas, err = h.Ligas.ListLigas(ctx)
		if err != nil {
			log.Printf("listar ligas para tabla: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	if err := Views.TablaList(ligas, tabla, int32(ligaID)).Render(ctx, w); err != nil {
		log.Printf("renderizar TablaList: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Delete atiende DELETE /tabla/{liga_id}/{equipo_id}: remueve al equipo de la liga
// y devuelve una respuesta vacía con 200 OK para que HTMX elimine la fila del DOM.
func (h *TablaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ligaStr := r.PathValue("liga_id")
	equipoStr := r.PathValue("equipo_id")
	if ligaStr == "" {
		ligaStr = r.FormValue("liga_id")
	}
	if equipoStr == "" {
		equipoStr = r.FormValue("equipo_id")
	}

	ligaID, err1 := strconv.Atoi(ligaStr)
	equipoID, err2 := strconv.Atoi(equipoStr)
	if err1 != nil || err2 != nil || ligaID <= 0 || equipoID <= 0 {
		http.Error(w, "datos inválidos", http.StatusBadRequest)
		return
	}

	arg := DataBase.DeleteEquipoDeTablaParams{
		LigaID:   int32(ligaID),
		EquipoID: int32(equipoID),
	}
	if err := h.Service.RemoverEquipo(r.Context(), arg); err != nil {
		log.Printf("sacar equipo de tabla: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
