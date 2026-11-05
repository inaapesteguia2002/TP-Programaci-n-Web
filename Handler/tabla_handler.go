package Handler

import (
	"log"
	"net/http"
	"strconv"

	"DeRabona/DataBase"
	"DeRabona/Logic"
)

type TablaHandler struct {
	Service *Logic.TablaService
}

// Create atiende POST /tabla: inscribe un equipo en una liga
// y redirige a la página principal (patrón Post/Redirect/Get).
func (h *TablaHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	if _, err := h.Service.RegistrarEquipo(r.Context(), arg); err != nil {
		log.Printf("inscribir equipo en liga: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *TablaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ligaID, err1 := strconv.Atoi(r.FormValue("liga_id"))
	equipoID, err2 := strconv.Atoi(r.FormValue("equipo_id"))
	if err1 != nil || err2 != nil || ligaID <= 0 || equipoID <= 0 {
		http.Error(w, "datos inválidos", http.StatusBadRequest)
		return
	}
	// AJUSTAR: campos reales de DeleteEquipoDeTablaParams
	arg := DataBase.DeleteEquipoDeTablaParams{LigaID: int32(ligaID), EquipoID: int32(equipoID)}
	if err := h.Service.RemoverEquipo(r.Context(), arg); err != nil {
		log.Printf("sacar equipo de tabla: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
