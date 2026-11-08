package repo

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	postgresImage = "postgres:16-alpine"
	mysqlImage    = "mysql:8.0"
)

var (
	startedContainers   []testcontainers.Container
	startedContainersMu sync.Mutex

	postgresOnce sync.Once
	postgresDB   *sqlx.DB
	postgresErr  error

	mysqlOnce sync.Once
	mysqlDB   *sqlx.DB
	mysqlErr  error
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	startedContainersMu.Lock()
	for _, container := range startedContainers {
		_ = container.Terminate(context.Background())
	}
	startedContainersMu.Unlock()
	os.Exit(exitCode)
}

func rememberContainer(container testcontainers.Container) {
	startedContainersMu.Lock()
	defer startedContainersMu.Unlock()
	startedContainers = append(startedContainers, container)
}

func requireContainers(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping container test in short mode")
	}
	testcontainers.SkipIfProviderIsNotHealthy(t)
}

func postgresDatabase(t *testing.T) *sqlx.DB {
	t.Helper()
	requireContainers(t)
	postgresOnce.Do(func() {
		ctx := context.Background()
		container, err := postgres.Run(ctx, postgresImage,
			postgres.WithDatabase("test"), postgres.WithUsername("user"), postgres.WithPassword("password"),
			postgres.BasicWaitStrategies(),
		)
		if err != nil {
			postgresErr = err
			return
		}
		rememberContainer(container)
		connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			postgresErr = err
			return
		}
		postgresDB, postgresErr = sqlx.ConnectContext(ctx, "pgx", connectionString)
	})
	if postgresErr != nil {
		t.Fatal(postgresErr)
	}
	return postgresDB
}

func mysqlDatabase(t *testing.T) *sqlx.DB {
	t.Helper()
	requireContainers(t)
	mysqlOnce.Do(func() {
		ctx := context.Background()
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        mysqlImage,
				ExposedPorts: []string{"3306/tcp"},
				Env:          map[string]string{"MYSQL_ROOT_PASSWORD": "password", "MYSQL_DATABASE": "test"},
				WaitingFor:   wait.ForListeningPort("3306/tcp").WithStartupTimeout(2 * time.Minute),
			},
			Started: true,
		})
		if err != nil {
			mysqlErr = err
			return
		}
		rememberContainer(container)
		endpoint, err := container.PortEndpoint(ctx, "3306/tcp", "")
		if err != nil {
			mysqlErr = err
			return
		}
		dsn := fmt.Sprintf("root:password@tcp(%s)/test?parseTime=true&multiStatements=false", endpoint)
		mysqlDB, mysqlErr = connectWithRetry(ctx, "mysql", dsn, 90*time.Second)
	})
	if mysqlErr != nil {
		t.Fatal(mysqlErr)
	}
	return mysqlDB
}

func connectWithRetry(ctx context.Context, driver, dsn string, timeout time.Duration) (*sqlx.DB, error) {
	deadline := time.Now().Add(timeout)
	for {
		db, err := sqlx.ConnectContext(ctx, driver, dsn)
		if err == nil {
			return db, nil
		}
		if time.Now().After(deadline) {
			return nil, err
		}
		time.Sleep(time.Second)
	}
}
