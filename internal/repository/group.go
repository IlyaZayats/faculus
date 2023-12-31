package repository

import (
	"context"
	"github.com/IlyaZayats/faculus/internal/entity"
	"github.com/IlyaZayats/faculus/internal/interfaces"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresGroupRepository struct {
	db *pgxpool.Pool
}

func NewPostgresGroupRepository(db *pgxpool.Pool) (interfaces.GroupRepository, error) {
	return &PostgresGroupRepository{
		db: db,
	}, nil
}

func (r *PostgresGroupRepository) GetGroups() ([]entity.Group, error) {
	var groups []entity.Group
	q := "SELECT id, faculty_id, name FROM Groups"
	rows, err := r.db.Query(context.Background(), q)
	if err != nil && err.Error() != "no rows in result set" {
		return groups, err
	}
	return r.parseRowsToSlice(rows)

}

func (r *PostgresGroupRepository) InsertGroup(group entity.Group) error {
	q := "INSERT INTO Groups (faculty_id, name) VALUES ($1, $2)"
	if _, err := r.db.Exec(context.Background(), q, group.FacultyId, group.Name); err != nil {
		return err
	}
	return nil
}

func (r *PostgresGroupRepository) UpdateGroup(group entity.Group) error {
	q := "UPDATE Groups SET name=$1 WHERE id=$2"
	if _, err := r.db.Exec(context.Background(), q, group.Name, group.Id); err != nil {
		return err
	}
	return nil
}

func (r *PostgresGroupRepository) DeleteGroup(id int) error {
	q := "DELETE FROM Groups WHERE id=$1"
	if _, err := r.db.Exec(context.Background(), q, id); err != nil {
		return err
	}
	return nil
}

func (r *PostgresGroupRepository) parseRowsToSlice(rows pgx.Rows) ([]entity.Group, error) {
	var slice []entity.Group
	defer rows.Close()
	for rows.Next() {
		var id, facultyId int
		var name string
		if err := rows.Scan(&id, &facultyId, &name); err != nil {
			return slice, err
		}
		slice = append(slice, entity.Group{Id: id, FacultyId: facultyId, Name: name})
	}
	return slice, nil
}
