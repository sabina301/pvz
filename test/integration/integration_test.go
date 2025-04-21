package integration

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"os"
	"path/filepath"
	"pvz/internal/models/product"
	"pvz/internal/models/pvz"
	"pvz/internal/models/reception"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"pvz/internal/logger"
	"pvz/internal/models/auth"
	"pvz/internal/repositories"
	"pvz/internal/services"
)

func TestFullFlow(t *testing.T) {
	logger.Init("debug")

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)
	connStr := "postgres://user:password@localhost:" + port.Port() + "/testdb?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	wd, _ := os.Getwd()
	sourceURL := "file://" + filepath.Join(wd, "..", "..", "internal", "migrations")

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	require.NoError(t, err)
	require.NoError(t, m.Up())

	userRepo := repositories.NewUserRepository(db)
	pvzRepo := repositories.NewPvzRepository(db)
	productRepo := repositories.NewProductRepository(db)
	receptionRepo := repositories.NewReceptionRepository(db)

	authService := services.NewAuthService(userRepo, db)
	pvzService := services.NewPvzService(pvzRepo, productRepo, receptionRepo, db)
	productService := services.NewProductService(productRepo, pvzRepo, receptionRepo, db)
	receptionService := services.NewReceptionService(receptionRepo, pvzRepo, db)

	t.Run("full user flow", func(t *testing.T) {
		_, err := authService.Register(auth.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Role:     "moderator",
		})
		require.NoError(t, err)

		pvzResp, err := pvzService.Create(pvz.CreateRequest{
			City: "Москва",
		})
		require.NoError(t, err)
		require.NotEmpty(t, pvzResp.Id)

		receptionResp, err := receptionService.Create(reception.CreateRequest{
			PvzId: pvzResp.Id,
		})
		require.NoError(t, err)
		require.Equal(t, reception.InProgressStatus, receptionResp.Status)

		productResp, err := productService.AddInReception(product.AddInReceptionRequest{
			PvzId: pvzResp.Id,
			Type:  "электроника",
		})
		require.NoError(t, err)
		require.NotEmpty(t, productResp.Id)

		closedReception, err := pvzService.CLoseLastReception(pvzResp.Id)
		require.NoError(t, err)
		require.Equal(t, reception.CloseStatus, closedReception.Status)

		start := time.Now().Add(-24 * time.Hour)
		end := time.Now()
		list, err := pvzService.ListWithFilterDate(pvz.ListRequest{
			StartDate: &start,
			EndDate:   &end,
			Page:      1,
			Limit:     10,
		})
		require.NoError(t, err)
		require.Len(t, list, 1)
		require.Len(t, list[0].Receptions, 1)
		require.Len(t, list[0].Receptions[0].Products, 1)
	})
}
