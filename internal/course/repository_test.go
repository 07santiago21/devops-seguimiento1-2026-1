package course

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockCourseDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	dialector := postgres.New(postgres.Config{Conn: dbMock})
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)
	return db, mock
}

// ========== CREATE TESTS ==========
func TestCourseRepository_Create_Success(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "courses"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c := &Course{Name: "Math", Description: "Mathematics", Credits: 4, Capacity: 30}
	err := repo.Create(c)
	assert.NoError(t, err)
	assert.NotEmpty(t, c.ID)
}

func TestCourseRepository_Create_GeneratesUUID(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "courses"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	c1 := &Course{Name: "Math", Description: "Mathematics", Credits: 4, Capacity: 30}
	c2 := &Course{Name: "Physics", Description: "Physics", Credits: 3, Capacity: 25}

	repo.Create(c1)
	repo.Create(c2)

	assert.NotEqual(t, c1.ID, c2.ID, "IDs should be different")
}

// ========== GETALL TESTS ==========
func TestCourseRepository_GetAll_Success(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "credits", "capacity"}).
		AddRow("1", "Math", "Mathematics", 4, 30).
		AddRow("2", "Physics", "Physics", 3, 25)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" ORDER BY created_at desc`)).
		WillReturnRows(rows)

	courses, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(courses))
	assert.Equal(t, "Math", courses[0].Name)
}

func TestCourseRepository_GetAll_Empty(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "credits", "capacity"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" ORDER BY created_at desc`)).
		WillReturnRows(rows)

	courses, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(courses))
}

func TestCourseRepository_GetAll_DatabaseError(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" ORDER BY created_at desc`)).
		WillReturnError(gorm.ErrInvalidDB)

	_, err := repo.GetAll()
	assert.Error(t, err)
}

// ========== GET TESTS ==========
func TestCourseRepository_Get_Success(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "credits", "capacity"}).
		AddRow("1", "Math", "Mathematics", 4, 30)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" WHERE "courses"."id" = $1`)).
		WithArgs("1", sqlmock.AnyArg()).
		WillReturnRows(rows)

	course, err := repo.Get("1")
	assert.NoError(t, err)
	assert.Equal(t, "Math", course.Name)
	assert.Equal(t, int32(4), course.Credits)
}

func TestCourseRepository_Get_NotFound(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "credits", "capacity"})
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "courses" WHERE "courses"."id" = $1`)).
		WithArgs("999", sqlmock.AnyArg()).
		WillReturnRows(rows)

	_, err := repo.Get("999")
	assert.Error(t, err)
}

// ========== DELETE TESTS ==========
func TestCourseRepository_Delete_Success(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "courses" WHERE "courses"."id" = $1`)).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete("1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCourseRepository_Delete_Error(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "courses"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	err := repo.Delete("1")
	assert.Error(t, err)
}

// ========== PATCH TESTS ==========

func TestCourseRepository_Patch_Error(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "courses"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	name := "Test"
	err := repo.Patch("1", &name, nil, nil, nil)
	assert.Error(t, err)
}

// ========== PUT TESTS ==========
func TestCourseRepository_Put_Success(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "courses"`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Put("1", "Physics", "Physics course", 3, 25)
	assert.NoError(t, err)
}

func TestCourseRepository_Put_NotFound(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "courses"`)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Put("999", "Physics", "Physics course", 3, 25)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCourseRepository_Put_DatabaseError(t *testing.T) {
	db, mock := setupMockCourseDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "courses"`)).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	err := repo.Put("1", "Physics", "Physics course", 3, 25)
	assert.Error(t, err)
}

// ========== REPOSITORY FACTORY TEST ==========
func TestNewRepository_CreatesRepository(t *testing.T) {
	db, _ := setupMockCourseDB(t)
	repo := NewRepository(db)

	assert.NotNil(t, repo)
}
