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

	// Dejamos que GORM maneje transacciones normalmente
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

		mock.ExpectQuery(`SELECT \* FROM "enrollments"`).WillReturnRows(rows)

		mock.ExpectQuery(`SELECT \* FROM "courses" WHERE "courses"\."id" = \$1`).
			WithArgs("c1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("c1", "Course 1"))

		mock.ExpectQuery(`SELECT \* FROM "students" WHERE "students"\."id" = \$1`).
			WithArgs("s1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow("s1", "Student 1"))

		res, err := r.GetAll()
		assert.NoError(t, err)
		assert.NotNil(t, res)
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

	t.Run("Get - Not Found", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)

		// GORM suele añadir LIMIT 1 al final
		mock.ExpectQuery(`SELECT \* FROM "enrollments" WHERE id = \$1`).
			WithArgs("invalid-id").
			WillReturnError(gorm.ErrRecordNotFound)

		res, err := r.Get("invalid-id")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Put - Success", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)

		mock.ExpectBegin()
		// Ajustado al orden exacto que mostró tu log de error: course_id, student_id, total_amount
		mock.ExpectExec(`UPDATE "enrollments" SET "course_id"=\$1,"student_id"=\$2,"total_amount"=\$3 WHERE id = \$4`).
			WithArgs("c-new", "s-new", 200.0, "uuid-existing").
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

	t.Run("Delete - Not Found", func(t *testing.T) {
		db, mock := setupTestDB(t)
		r := NewRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "enrollments" WHERE id = \$1`).
			WithArgs("missing-id").
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := r.Delete("missing-id")
		assert.Error(t, err)
	})
}
