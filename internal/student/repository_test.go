package student

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	dialector := postgres.New(postgres.Config{Conn: dbMock})
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)
	return db, mock
}

func TestRepository_Get(t *testing.T) {
	t.Run("Get - Success", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewRepository(db)

		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("1", "Juan")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students" WHERE "students"."id" = $1`)).
			WithArgs("1", 1).
			WillReturnRows(rows)

		res, err := repo.Get("1")
		assert.NoError(t, err)
		assert.Equal(t, "Juan", res.Name)
	})

	t.Run("Get - Error", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students"`)).
			WillReturnError(errors.New("db error"))

		res, err := repo.Get("1")
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestRepository_Errors(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	t.Run("GetAll - Error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students"`)).
			WillReturnError(errors.New("query error"))

		res, err := repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Delete - Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "students"`)).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		err := repo.Delete("1")
		assert.Error(t, err)
	})
}

func TestRepository_Create_GetAll_Delete_Patch(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	// Create
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "students"`)).
		WithArgs(sqlmock.AnyArg(), "A", "B", 10, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	s := &Student{Name: "A", LastName: "B", Age: 10}
	err := repo.Create(s)
	assert.NoError(t, err)

	// GetAll
	rows := sqlmock.NewRows([]string{"id", "name", "last_name", "age"}).
		AddRow("1", "A", "B", 10)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students" ORDER BY created_at desc`)).
		WillReturnRows(rows)

	_, err = repo.GetAll()
	assert.NoError(t, err)

	// Patch - All fields
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "students" SET`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	name := "Jane"
	last := "Doe"
	age := int32(25)
	err = repo.Patch("1", &name, &last, &age)
	assert.NoError(t, err)

	// Delete
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "students"`)).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.Delete("1")
	assert.NoError(t, err)
}

func TestRepository_Put_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	// GORM incluye SET "id", "name", "last_name", "age" Y EL WHERE id
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "students" SET`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "99").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Put("99", "N", "L", 20)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
