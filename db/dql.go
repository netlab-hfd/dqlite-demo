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
	db       string
	cluster  *[]string
	diskmode bool
	crt      string
	key      string
	dir      string
	ctx      context.Context
}

func NewDqlLiteConnection(config *DQLConfig) (*sql.DB, error) {
	options := []app.Option{
		app.WithAddress(config.db),
		app.WithCluster(*config.cluster),
		app.WithLogFunc(logging.LogFunc),
		app.WithDiskMode(config.diskmode),
	}

	if (config.crt != "" && config.key == "") || (config.key != "" && config.crt == "") {
		return nil, fmt.Errorf("both tls certificate and a key must be given")
	}

	if config.crt != "" {
		cert, err := tls.LoadX509KeyPair(config.crt, config.key)
		if err != nil {
			return nil, err
		}

		data, err := os.ReadFile(config.crt)
		if err != nil {
			return nil, err
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("bad certificate passed")
		}

		options = append(options, app.WithTLS(app.SimpleTLSConfig(cert, pool)))
	}

	app, err := app.New(config.dir)
	if err != nil {
		return nil, err
	}

	if err := app.Ready(config.ctx); err != nil {
		return nil, err
	}

	db, err := app.Open(config.ctx, "demo")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, err
}
