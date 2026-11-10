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
	DB   *sql.DB // se usa para abrir transacciones
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

// revertirEnTabla deshace el efecto de un partido finalizado en la tabla de posiciones.
// Si el partido estaba pendiente nunca impactó en la tabla, así que no hace nada.
func revertirEnTabla(ctx context.Context, q *DataBase.Queries, p DataBase.Partido) error {
	if p.Estado.String != "finalizado" {
		return nil
	}
	gl, gv := p.GolesLocal.Int32, p.GolesVisitante.Int32

	if err := restarATabla(ctx, q, p.LigaID, p.EquipoLocalID, gl, gv); err != nil {
		return err
	}
	return restarATabla(ctx, q, p.LigaID, p.EquipoVisitanteID, gv, gl)
}

// restarATabla resta a la fila de un equipo lo que le había sumado un partido.
func restarATabla(ctx context.Context, q *DataBase.Queries, ligaID, equipoID, propios, rival int32) error {
	pts, g, e, perd := resultado(propios, rival)
	n := func(v int32) sql.NullInt32 { return sql.NullInt32{Int32: v, Valid: true} }

	return q.RestarTabla(ctx, DataBase.RestarTablaParams{
		LigaID:            ligaID,
		EquipoID:          equipoID,
		Puntos:            n(pts),
		PartidosGanados:   n(g),
		PartidosEmpatados: n(e),
		PartidosPerdidos:  n(perd),
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

// DeletePartido elimina un partido y, si estaba finalizado, revierte su efecto en la tabla de posiciones.
// Va en una transacción: o se borra el partido y se corrige la tabla, o no se hace nada.
func (s *PartidosService) DeletePartido(ctx context.Context, id int32) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // si ya se hizo Commit, no tiene efecto

	q := s.Repo.WithTx(tx)

	partido, err := q.GetPartido(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil // ya no existe, no hay nada que borrar
	}
	if err != nil {
		return err
	}

	if err := revertirEnTabla(ctx, q, partido); err != nil {
		return err
	}
	if err := q.DeletePartido(ctx, id); err != nil {
		return err
	}
	return tx.Commit()
}
