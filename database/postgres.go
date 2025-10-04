package database

import (
	"context"
	"database/sql"

	"github.com/farinas09/go-grpc/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (repo *PostgresRepository) SetStudent(ctx context.Context, student *models.Student) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO students (id, name, age) VALUES ($1, $2, $3)", student.Id, student.Name, student.Age)
	return err
}

func (repo *PostgresRepository) GetStudent(ctx context.Context, id string) (*models.Student, error) {
	var student models.Student

	err := repo.db.QueryRowContext(ctx, "SELECT id, name, age FROM students WHERE id = $1", id).
		Scan(&student.Id, &student.Name, &student.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &student, nil
}

func (repo *PostgresRepository) SetTest(ctx context.Context, test *models.Test) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO tests (id, name) VALUES ($1, $2)", test.Id, test.Name)
	return err
}

func (repo *PostgresRepository) GetTest(ctx context.Context, id string) (*models.Test, error) {
	var test models.Test

	err := repo.db.QueryRowContext(ctx, "SELECT id, name FROM tests WHERE id = $1", id).
		Scan(&test.Id, &test.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &test, nil
}

func (repo *PostgresRepository) SetQuestion(ctx context.Context, question *models.Question) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO questions (id, question, answer, test_id) VALUES ($1, $2, $3, $4)", question.Id, question.Question, question.Answer, question.TestId)
	return err
}

func (repo *PostgresRepository) SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO enrollments (student_id, test_id) VALUES ($1, $2)", enrollment.StudentId, enrollment.TestId)
	return err
}

func (repo *PostgresRepository) GetStudentsPerTest(ctx context.Context, testId string) ([]*models.Student, error) {
	var students []*models.Student

	rows, err := repo.db.QueryContext(ctx, "SELECT students.id, students.name, students.age FROM students JOIN enrollments ON students.id = enrollments.student_id WHERE enrollments.test_id = $1", testId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, &student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (repo *PostgresRepository) GetQuestionsPerTest(ctx context.Context, testId string) ([]*models.Question, error) {
	var questions []*models.Question

	rows, err := repo.db.QueryContext(ctx, "SELECT id, question FROM questions WHERE test_id = $1", testId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var question models.Question
		err := rows.Scan(&question.Id, &question.Question)
		if err != nil {
			return nil, err
		}
		questions = append(questions, &question)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}
