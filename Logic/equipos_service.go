package Logic

import (
	"DeRabona/DataBase"
	"context"
	"database/sql"
	"errors"
)

// EquiposService maneja la lógica de negocio para los equipos
type EquiposService struct {
	Repo *DataBase.Queries
	DB   *sql.DB // se usa para abrir transacciones
}

// CreateEquipo valida y guarda un nuevo equipo
func (s *EquiposService) CreateEquipo(ctx context.Context, nombre string) (DataBase.Equipo, error) {
	if nombre == "" {
		return DataBase.Equipo{}, errors.New("el nombre del equipo no puede estar vacío")
	}
	return s.Repo.CreateEquipo(ctx, nombre)
}

// ListEquipos devuelve todos los equipos
func (s *EquiposService) ListEquipos(ctx context.Context) ([]DataBase.Equipo, error) {
	return s.Repo.ListEquipos(ctx)
}

// GetEquipo obtiene un equipo por su ID
func (s *EquiposService) GetEquipo(ctx context.Context, id int32) (DataBase.Equipo, error) {
	return s.Repo.GetEquipo(ctx, id)
}

// UpdateEquipo actualiza los datos de un equipo existente
func (s *EquiposService) UpdateEquipo(ctx context.Context, arg DataBase.UpdateEquipoParams) error {
	if arg.Nombre == "" {
		return errors.New("el nuevo nombre no puede estar vacío")
	}
	return s.Repo.UpdateEquipo(ctx, arg)
}

// DeleteEquipo elimina un equipo y todo lo que depende de él:
//  1. por cada partido que jugó, revierte el resultado en la tabla del rival y borra el partido
//  2. borra el equipo (su fila en las tablas de posiciones se elimina sola, por ON DELETE CASCADE)

// Todo va en una transacción: si algo falla, no se borra nada.
func (s *EquiposService) DeleteEquipo(ctx context.Context, id int32) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // si ya se hizo Commit, no tiene efecto

	q := s.Repo.WithTx(tx)

	partidos, err := q.ListPartidosEquipo(ctx, id)
	if err != nil {
		return err
	}
	for _, p := range partidos {
		if err := revertirEnTabla(ctx, q, p); err != nil {
			return err
		}
		if err := q.DeletePartido(ctx, p.ID); err != nil {
			return err
		}
	}

	if err := q.DeleteEquipo(ctx, id); err != nil {
		return err
	}
	return tx.Commit()
}
