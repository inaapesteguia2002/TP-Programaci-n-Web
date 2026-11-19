package Handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"DeRabona/DataBase"
	"DeRabona/Logic"
	"DeRabona/Views"
)

type PartidosHandler struct {
	Service *Logic.PartidosService
	Equipos *Logic.EquiposService
	Ligas *Logic.LigasService
	Tabla *Logic.TablaService
}

// renderPartidos devuelve el fragmento actualizado de la lista de partidos junto con el estado de error (si lo hubiera).
func (h *PartidosHandler) renderPartidos(w http.ResponseWriter, r *http.Request, errorMsg string) {
	ctx := r.Context()

	partidos, err := h.Service.ListPartidos(ctx)
	if err != nil {
		log.Printf("listar partidos: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var equipos []DataBase.Equipo
	if h.Equipos != nil {
		equipos, err = h.Equipos.ListEquipos(ctx)
		if err != nil {
			log.Printf("listar equipos para partidos: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	if err := Views.PartidoListConEstado(partidos, equipos, errorMsg).Render(ctx, w); err != nil {
		log.Printf("renderizar PartidoList: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// Create atiende POST /partidos: valida y guarda el partido, vuelve a consultar la lista
// y devuelve únicamente el fragmento de la lista de partidos para HTMX.
func (h *PartidosHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario inválido", http.StatusBadRequest)
		return
	}

	ligaID, err1 := strconv.Atoi(r.FormValue("liga_id"))
	localID, err2 := strconv.Atoi(r.FormValue("equipo_local_id"))
	visitanteID, err3 := strconv.Atoi(r.FormValue("equipo_visitante_id"))
	golesLocal, err4 := strconv.Atoi(r.FormValue("goles_local"))
	golesVisitante, err5 := strconv.Atoi(r.FormValue("goles_visitante"))
	estado := r.FormValue("estado")

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || ligaID <= 0 || localID <= 0 || visitanteID <= 0 {
		h.renderPartidos(w, r, "los datos del partido no son válidos")
		return
	}
	if golesLocal < 0 || golesVisitante < 0 {
		h.renderPartidos(w, r, "los goles no pueden ser negativos")
		return
	}
	if localID == visitanteID {
		h.renderPartidos(w, r, "un equipo no puede jugar contra sí mismo")
		return
	}
	if estado != "pendiente" && estado != "finalizado" {
		h.renderPartidos(w, r, "estado inválido")
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
		if errors.Is(err, Logic.ErrEquipoFueraDeLiga) {
			h.renderPartidos(w, r, err.Error())
			return
		}
		log.Printf("crear partido: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.renderPartidos(w, r, int32(ligaID))
}

// Delete atiende DELETE /partidos/{id}: obtiene el ID de la URL, elimina el partido
// y devuelve una respuesta vacía con 200 OK para que HTMX quite la fila del DOM.
func (h *PartidosHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		idStr = r.FormValue("id")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeletePartido(r.Context(), int32(id)); err != nil {
		log.Printf("eliminar partido: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PartidosHandler) renderPartidosConTabla(w http.ResponseWriter, r *http.Request, ligaID int32) {
	ctx := r.Context()

	partidos, err := h.Service.ListPartidos(ctx)
	if err != nil {
		log.Printf("listar partidos: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var equipos []DataBase.Equipo
	if h.Equipos != nil {
		equipos, err = h.Equipos.ListEquipos(ctx)
		if err != nil {
			log.Printf("listar equipos para partidos: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Si no se inyectaron Ligas/Tabla (ej. en algún test unitario), cae al render clásico
	if h.Ligas == nil || h.Tabla == nil {
		if err := Views.PartidoListConEstado(partidos, equipos, "").Render(ctx, w); err != nil {
			log.Printf("renderizar PartidoList: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	ligas, err := h.Ligas.ListLigas(ctx)
	if err != nil {
		log.Printf("listar ligas para tabla: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tabla, err := h.Tabla.ObtenerTabla(ctx, ligaID)
	if err != nil {
		log.Printf("obtener tabla tras crear partido: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := Views.PartidoListConTabla(partidos, equipos, ligas, tabla, ligaID).Render(ctx, w); err != nil {
		log.Printf("renderizar PartidoListConTabla: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
