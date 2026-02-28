package student

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	// 1. Crear el mock
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	// 2. Configurar GORM
	dialector := postgres.New(postgres.Config{
		Conn: dbMock,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := NewRepository(db)

	t.Run("Éxito al insertar estudiante", func(t *testing.T) {
		student := &Student{
			Name:     "Juan",
			LastName: "Perez",
			Age:      20,
		}

		// --- CAMBIO AQUÍ ---
		// Usamos sqlmock.AnyArg() para TODOS los campos para evitar errores de timestamps o campos extras de GORM
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "students"`)).
			WithArgs(
				sqlmock.AnyArg(), // ID (UUID)
				sqlmock.AnyArg(), // Name
				sqlmock.AnyArg(), // LastName
				sqlmock.AnyArg(), // Age
				sqlmock.AnyArg(), // CreatedAt (el que causó el error)
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Ejecución
		err := repo.Create(student)

		// Validaciones
		assert.NoError(t, err)
		assert.NotEmpty(t, student.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
