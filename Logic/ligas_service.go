package Logic

import (
	"context"
	"errors"

	"DeRabona/DataBase"
)

// LigasService maneja las reglas de negocio para las ligas
type LigasService struct {
	Repo *DataBase.Queries
}

// CreateLiga valida y guarda una nueva liga
func (s *LigasService) CreateLiga(ctx context.Context, arg DataBase.CreateLigaParams) (DataBase.Liga, error) {
	if arg.Nombre == "" || arg.Pais == "" {
		return DataBase.Liga{}, errors.New("el nombre y el país de la liga son obligatorios")
	}
	return s.Repo.CreateLiga(ctx, arg)
}

// ListLigas devuelve todas las ligas registradas
func (s *LigasService) ListLigas(ctx context.Context) ([]DataBase.Liga, error) {
	return s.Repo.ListLigas(ctx)
}

// GetLiga obtiene una liga específica por su ID
func (s *LigasService) GetLiga(ctx context.Context, id int32) (DataBase.Liga, error) {
	return s.Repo.GetLiga(ctx, id)
}

// UpdateLiga actualiza los datos de una liga existente
func (s *LigasService) UpdateLiga(ctx context.Context, arg DataBase.UpdateLigaParams) error {
	if arg.Nombre == "" || arg.Pais == "" {
		return errors.New("el nombre y el país no pueden estar vacíos al actualizar")
	}
	return s.Repo.UpdateLiga(ctx, arg)
}

// DeleteLiga elimina una liga por su ID
func (s *LigasService) DeleteLiga(ctx context.Context, id int32) error {
	return s.Repo.DeleteLiga(ctx, id)
}
