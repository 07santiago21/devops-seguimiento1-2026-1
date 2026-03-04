package student

import (
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
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("1", "Juan")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students" WHERE "students"."id" = $1`)).
		WithArgs("1", sqlmock.AnyArg()).
		WillReturnRows(rows)

	res, err := repo.Get("1")
	assert.NoError(t, err)
	assert.Equal(t, "Juan", res.Name)
}

func TestRepository_Put_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "students"`)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Put("99", "N", "L", 20)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestRepository_Create_GetAll_Delete_Patch(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	// create
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "students"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	s := &Student{Name: "A", LastName: "B", Age: 10}
	err := repo.Create(s)
	assert.NoError(t, err)
	assert.NotEmpty(t, s.ID)

	// getall
	rows := sqlmock.NewRows([]string{"id", "name", "last_name", "age"}).
		AddRow("1", "A", "B", 10)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "students" ORDER BY created_at desc`)).
		WillReturnRows(rows)

	_, err = repo.GetAll()
	assert.NoError(t, err)

	// patch (partial update)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "students"`)).
		WithArgs("Jane", "1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	name := "Jane"
	err = repo.Patch("1", &name, nil, nil)
	assert.NoError(t, err)

	// delete
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "students"`)).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.Delete("1")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
