package containers

import (
	"github.com/amirrmonfared/testMicroServices/project/pkg/postgresql"
	"github.com/ory/dockertest/v3"
	"gorm.io/gorm"
)

type PostgreSqlContainer struct {
	pool      *dockertest.Pool
	resource  *dockertest.Resource
	imagename string
	opts      postgresql.Opts
}

type IPostgreSqlContainer interface {
	C() PostgreSqlContainer
	create() error
	Connect() *gorm.DB
	AutoMigrate(db *gorm.DB)
	Flush(db *gorm.DB)
}

func NewPostgresqlContainer(pool *dockertest.Pool) PostgreSqlContainer {
	opts := postgresql.Opts{
		Host:     "localhost",
		User:     "testcontainer",
		Password: "Aa123456.",
		Database: "testcontainer",
		Port:     5432,
	}

	return PostgreSqlContainer{pool: pool, opts: opts, imagename: "postgresql-testcontainer"}
}

// func (container PostgreSqlContainer) C() PostgreSqlContainer {
// 	return container
// }

// func (container PostgreSqlContainer) Create() error {
// 	if IsRunning(*container.pool, container.imagename) {
// 		return nil
// 	}

// 	dockerOpts := dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "latest",
// 		Env: []string{
// 			"POSTGRES_USER=" + container.opts.User,
// 			"POSTGRES_PASSWORD=" + container.opts.Password,
// 			"POSTGRES_DB=" + container.opts.Database,
// 		},
// 		ExposedPorts: []string{strconv.Itoa(container.opts.Port)},
// 		PortBindings: map[docker.Port][]docker.PortBinding{
// 			docker.Port(strconv.Itoa(container.opts.Port)): {{HostIP: "0.0.0.0", HostPort: strconv.Itoa(container.opts.Port)}},
// 		},
// 		Name: container.imagename,
// 	}

// 	resource, err := container.pool.RunWithOptions(&dockerOpts)
// 	if err != nil {
// 		log.Fatalf("Could not start resource (Postgresql Test Container): %s", err.Error())
// 		return err
// 	}

// 	container.resource = resource
// 	return nil
// }
