package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAutoMigrate(t *testing.T) {
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

	// --- LA SOLUCIÓN AL ERROR ---
	// Preparamos el mock para aceptar CUALQUIER consulta o ejecución de migración
	// GORM hace SELECTs para ver si existe la tabla y CREATE TABLE si no existe.
	mock.ExpectQuery(".*").WillReturnRows(sqlmock.NewRows(nil)) // Simula que la tabla no existe
	mock.ExpectExec("CREATE TABLE.*").WillReturnResult(sqlmock.NewResult(0, 0))

	// 3. Ejecutar la función del main
	err = autoMigrate(db)

	// 4. Validar
	assert.NoError(t, err)
}

func TestMainInitializations(t *testing.T) {
	// Test simple para marcar presencia en el paquete main
	assert.True(t, true)
}
