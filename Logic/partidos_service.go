package Logic

import (
	"context"
	"database/sql"
	"errors"

	"DeRabona/DataBase"
)

// ErrEquipoFueraDeLiga se devuelve cuando alguno de los equipos del partido no está inscripto en la liga elegida.
var ErrEquipoFueraDeLiga = errors.New("ambos equipos deben estar inscriptos en la liga seleccionada")

// PartidosService maneja las reglas de negocio para los partidos
type PartidosService struct {
	Repo *DataBase.Queries
}

// CreatePartido valida y guarda un nuevo partido.
// Si el partido ya está finalizado, también actualiza la tabla de posiciones.
func (s *PartidosService) CreatePartido(ctx context.Context, arg DataBase.CreatePartidoParams) (DataBase.Partido, error) {

	// Regla de negocio: Un equipo no puede jugar contra sí mismo
	if arg.EquipoLocalID == arg.EquipoVisitanteID {
		return DataBase.Partido{}, errors.New("un equipo no puede jugar contra sí mismo")
	}
	if arg.LigaID <= 0 {
		return DataBase.Partido{}, errors.New("el ID de la liga es inválido")
	}

	// Regla de negocio: ambos equipos deben pertenecer a la liga del partido
	inscriptos, err := s.Repo.GetTablaPorLiga(ctx, arg.LigaID)
	if err != nil {
		return DataBase.Partido{}, err
	}
	var localEnLiga, visitanteEnLiga bool
	for _, fila := range inscriptos {
		if fila.EquipoID == arg.EquipoLocalID {
			localEnLiga = true
		}
		if fila.EquipoID == arg.EquipoVisitanteID {
			visitanteEnLiga = true
		}
	}
	if !localEnLiga || !visitanteEnLiga {
		return DataBase.Partido{}, ErrEquipoFueraDeLiga
	}

	partido, err := s.Repo.CreatePartido(ctx, arg)
	if err != nil {
		return DataBase.Partido{}, err
	}

	if partido.Estado.String == "finalizado" {
		gl, gv := partido.GolesLocal.Int32, partido.GolesVisitante.Int32

		if err := s.sumarATabla(ctx, partido.LigaID, partido.EquipoLocalID, gl, gv); err != nil {
			return DataBase.Partido{}, err
		}
		if err := s.sumarATabla(ctx, partido.LigaID, partido.EquipoVisitanteID, gv, gl); err != nil {
			return DataBase.Partido{}, err
		}
	}

	return partido, nil
}

// resultado devuelve puntos, ganados, empatados y perdidos para un equipo
// según los goles que hizo y los que le hicieron.
func resultado(propios, rival int32) (puntos, ganados, empatados, perdidos int32) {
	switch {
	case propios > rival:
		return 3, 1, 0, 0
	case propios < rival:
		return 0, 0, 0, 1
	default:
		return 1, 0, 1, 0
	}
}

// sumarATabla impacta el resultado de un partido en la fila de un equipo.
func (s *PartidosService) sumarATabla(ctx context.Context, ligaID, equipoID, propios, rival int32) error {
	pts, g, e, p := resultado(propios, rival)
	n := func(v int32) sql.NullInt32 { return sql.NullInt32{Int32: v, Valid: true} }

	return s.Repo.UpdateTabla(ctx, DataBase.UpdateTablaParams{
		LigaID:            ligaID,
		EquipoID:          equipoID,
		Puntos:            n(pts),
		PartidosGanados:   n(g),
		PartidosEmpatados: n(e),
		PartidosPerdidos:  n(p),
	})
}

// ListPartidos devuelve todos los partidos
func (s *PartidosService) ListPartidos(ctx context.Context) ([]DataBase.Partido, error) {
	return s.Repo.ListPartidos(ctx)
}

// GetPartido obtiene un partido por su ID
func (s *PartidosService) GetPartido(ctx context.Context, id int32) (DataBase.Partido, error) {
	return s.Repo.GetPartido(ctx, id)
}

// UpdatePartido actualiza los datos de un partido
func (s *PartidosService) UpdatePartido(ctx context.Context, arg DataBase.UpdatePartidoParams) error {
	if !arg.Estado.Valid || arg.Estado.String == "" {
		arg.Estado = sql.NullString{String: "pendiente", Valid: true}
	}
	return s.Repo.UpdatePartido(ctx, arg)
}

// DeletePartido elimina un partido por su ID
func (s *PartidosService) DeletePartido(ctx context.Context, id int32) error {
	return s.Repo.DeletePartido(ctx, id)
}
