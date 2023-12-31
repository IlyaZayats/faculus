package repository

import (
	"context"
	"github.com/IlyaZayats/faculus/internal/entity"
	"github.com/IlyaZayats/faculus/internal/interfaces"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStudentRepository struct {
	db *pgxpool.Pool
}

func NewPostgresStudentRepository(db *pgxpool.Pool) (interfaces.StudentRepository, error) {
	return &PostgresStudentRepository{
		db: db,
	}, nil
}

func (r *PostgresStudentRepository) GetStudents() ([]entity.Student, error) {
	var students []entity.Student
	q := "SELECT id, group_id, firstname, lastname, middlename, birthdate, phone, sex FROM Students"
	rows, err := r.db.Query(context.Background(), q)
	if err != nil && err.Error() != "no rows in result set" {
		return students, err
	}
	//students, err =
	return r.parseRowsToSlice(rows)

}

func (r *PostgresStudentRepository) InsertStudent(student entity.Student) error {
	q := "INSERT INTO Students (lastname, firstname, middlename, birthdate, group_id, phone, sex) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if _, err := r.db.Exec(context.Background(), q, student.LastName, student.FirstName, student.MiddleName, student.BirthDate, student.GroupId, student.PhoneNumber, student.Gender); err != nil {
		return err
	}
	return nil
}

func (r *PostgresStudentRepository) UpdateStudent(student entity.Student) error {
	q := "UPDATE Students SET (lastname, firstname, middlename, birthdate, phone, sex) = ($1, $2, $3, $4, $5, $6) WHERE id=$7"
	if _, err := r.db.Exec(context.Background(), q, student.LastName, student.FirstName, student.MiddleName, student.BirthDate, student.PhoneNumber, student.Gender, student.Id); err != nil {
		return err
	}
	return nil
}

func (r *PostgresStudentRepository) DeleteStudent(id int) error {
	q := "DELETE FROM Students WHERE id=$1"
	if _, err := r.db.Exec(context.Background(), q, id); err != nil {
		return err
	}
	return nil
}

func (r *PostgresStudentRepository) parseRowsToSlice(rows pgx.Rows) ([]entity.Student, error) {
	var slice []entity.Student
	defer rows.Close()
	for rows.Next() {
		var id, groupId, gender int
		var firstName, lastName, middleName, phone string
		var birthDate pgtype.Date
		if err := rows.Scan(&id, &groupId, &firstName, &lastName, &middleName, &birthDate, &phone, &gender); err != nil {
			return slice, err
		}
		birthDateString := birthDate.Time.Format("2006.01.02")
		slice = append(slice, entity.Student{Id: id, GroupId: groupId, LastName: lastName, FirstName: firstName, MiddleName: middleName, BirthDate: birthDateString, PhoneNumber: phone, Gender: gender})
	}
	return slice, nil
}
