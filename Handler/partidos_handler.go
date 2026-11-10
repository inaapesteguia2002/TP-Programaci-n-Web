package Handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"DeRabona/DataBase"
	"DeRabona/Logic"
)

type PartidosHandler struct {
	Service *Logic.PartidosService
}

// errorEnForm vuelve a la página principal (sección partidos) mostrando el mensaje en la propia interfaz, en vez de mostrar una página de error aparte.
func errorEnForm(w http.ResponseWriter, r *http.Request, msg string) {
	http.Redirect(w, r, "/?error="+url.QueryEscape(msg)+"#seccion-partidos", http.StatusSeeOther)
}

// Create atiende POST /partidos: parsea y valida el formulario, guarda el partido
// y redirige a la página principal (patrón Post/Redirect/Get).
func (h *PartidosHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		errorEnForm(w, r, "formulario inválido")
		return
	}

	ligaID, err1 := strconv.Atoi(r.FormValue("liga_id"))
	localID, err2 := strconv.Atoi(r.FormValue("equipo_local_id"))
	visitanteID, err3 := strconv.Atoi(r.FormValue("equipo_visitante_id"))
	golesLocal, err4 := strconv.Atoi(r.FormValue("goles_local"))
	golesVisitante, err5 := strconv.Atoi(r.FormValue("goles_visitante"))
	estado := r.FormValue("estado")

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		errorEnForm(w, r, "los datos del partido no son válidos")
		return
	}
	if golesLocal < 0 || golesVisitante < 0 {
		errorEnForm(w, r, "los goles no pueden ser negativos")
		return
	}
	if localID == visitanteID {
		errorEnForm(w, r, "un equipo no puede jugar contra sí mismo")
		return
	}
	if estado != "pendiente" && estado != "finalizado" {
		errorEnForm(w, r, "estado inválido")
		return
	}

	arg := DataBase.CreatePartidoParams{
		LigaID:            int32(ligaID),
		EquipoLocalID:     int32(localID),
		EquipoVisitanteID: int32(visitanteID),
		GolesLocal:        sql.NullInt32{Int32: int32(golesLocal), Valid: true},
		GolesVisitante:    sql.NullInt32{Int32: int32(golesVisitante), Valid: true},
		Estado:            sql.NullString{String: estado, Valid: true},
	}

	if _, err := h.Service.CreatePartido(r.Context(), arg); err != nil {

		// Regla de negocio: los dos equipos tienen que estar inscriptos en la liga
		if errors.Is(err, Logic.ErrEquipoFueraDeLiga) {
			errorEnForm(w, r, err.Error())
			return
		}
		log.Printf("crear partido: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PartidosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeletePartido(r.Context(), int32(id)); err != nil {
		log.Printf("eliminar partido: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
