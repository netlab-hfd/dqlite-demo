package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"

	"github.com/canonical/go-dqlite/v3/app"
	"github.com/netlab-hfd/dqlite-demo/logging"
)

const schema = "CREATE TABLE IF NOT EXISTS model (key TEXT, value TEXT, UNIQUE(key))"

type DQLConfig struct {
	Db       string
	Cluster  *[]string
	Diskmode bool
	Crt      string
	Key      string
	Dir      string
	Ctx      context.Context
}

func NewDqlLiteConnection(config *DQLConfig) (*sql.DB, error) {
	options := []app.Option{
		app.WithAddress(config.Db),
		app.WithCluster(*config.Cluster),
		app.WithLogFunc(logging.LogFunc),
		app.WithDiskMode(config.Diskmode),
	}

	if (config.Crt != "" && config.Key == "") || (config.Key != "" && config.Crt == "") {
		return nil, fmt.Errorf("both tls certificate and a key must be given")
	}

	if config.Crt != "" {
		cert, err := tls.LoadX509KeyPair(config.Crt, config.Key)
		if err != nil {
			fmt.Println("Certificate error")
			return nil, err
		}

		data, err := os.ReadFile(config.Crt)
		if err != nil {
			fmt.Println("Cert could not be read")
			return nil, err
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("bad certificate passed")
		}

		options = append(options, app.WithTLS(app.SimpleTLSConfig(cert, pool)))
	}

	app, err := app.New(config.Dir, options...)
	if err != nil {
		fmt.Println("failed starting app with cluster")
		return nil, err
	}

	if err := app.Ready(config.Ctx); err != nil {
		fmt.Println("app didnt get ready")
		return nil, err
	}

	db, err := app.Open(config.Ctx, "demo")
	if err != nil {
		fmt.Println("failed connecting to database")
		return nil, err
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, err
}
