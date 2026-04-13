package enrollment

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %s", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn: dbMock,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening gorm: %s", err)
	}

	return db, mock
}

func TestRepository_Queries(t *testing.T) {

	t.Run("GetAll - Success", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)

		rows := sqlmock.NewRows([]string{"id", "student_id", "course_id", "total_amount"}).
			AddRow("uuid-1", "s1", "c1", 500.0)

		// GORM usa ORDER BY created_at desc según tu implementación
		mock.ExpectQuery(`SELECT \* FROM "enrollments" ORDER BY created_at desc`).WillReturnRows(rows)

		// Mocks para Preloads
		mock.ExpectQuery(`SELECT \* FROM "courses"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c1"))
		mock.ExpectQuery(`SELECT \* FROM "students"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("s1"))

		res, err := r.GetAll()
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("Get - Success", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)

		rows := sqlmock.NewRows([]string{"id", "student_id", "course_id"}).AddRow("uuid-1", "s1", "c1")

		// GORM .First() añade LIMIT $2. Esperamos el ID y el valor 1.
		mock.ExpectQuery(`SELECT \* FROM "enrollments" WHERE id = \$1`).
			WithArgs("uuid-1", 1).
			WillReturnRows(rows)

		mock.ExpectQuery(`SELECT \* FROM "courses"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c1"))
		mock.ExpectQuery(`SELECT \* FROM "students"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("s1"))

		res, err := r.Get("uuid-1")
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Get - Not Found", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)

		// Importante: .WithArgs("invalid-id", 1) por el LIMIT de .First()
		mock.ExpectQuery(`SELECT \* FROM "enrollments" WHERE id = \$1`).
			WithArgs("invalid-id", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		res, err := r.Get("invalid-id")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Create - Success", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)
		e := &Enrollment{StudentID: "s1", CourseID: "c1", TotalAmount: 100.0}

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "enrollments"`).
			WithArgs(sqlmock.AnyArg(), e.StudentID, e.CourseID, e.TotalAmount, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := r.Create(e)
		assert.NoError(t, err)
	})

	t.Run("Put - Success", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)

		mock.ExpectBegin()
		// GORM puede variar el orden, usamos AnyArg si el orden falla,
		// pero aquí seguimos el patrón de tu error log
		mock.ExpectExec(`UPDATE "enrollments" SET`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "uuid-existing").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := r.Put("uuid-existing", "s-new", "c-new", 200.0)
		assert.NoError(t, err)
	})

	t.Run("Patch - Success", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)
		amt := 300.0

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "enrollments" SET "total_amount"=\$1 WHERE id = \$2`).
			WithArgs(amt, "uuid-1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := r.Patch("uuid-1", &amt)
		assert.NoError(t, err)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "enrollments" WHERE id = \$1`).
			WithArgs("uuid-to-delete").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := r.Delete("uuid-to-delete")
		assert.NoError(t, err)
	})
}
