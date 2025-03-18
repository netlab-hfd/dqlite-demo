package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"

	data "github.com/netlab-hfd/dqlite-demo/db"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

const (
	query  = "SELECT value FROM model WHERE key = ?"
	update = "INSERT OR REPLACE INTO model(key, value) VALUES(?, ?)"
)

var (
	api      string
	db       string
	join     *[]string
	dir      string
	verbose  bool
	diskMode bool
	crt      string
	key      string
)

var rootCmd = &cobra.Command{
	Use:   "dqldemo",
	Short: "dqldemo showcases distributed systems concepts with dqlite",
	Long:  "A small API that allows setting key value pairs in a distributed sqlite database.",
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()

		dbOpts := data.DQLConfig{
			Db:       db,
			Cluster:  join,
			Diskmode: diskMode,
			Crt:      crt,
			Key:      key,
			Dir:      dir,
			Ctx:      ctx,
		}

		db, err := data.NewDqlLiteConnection(&dbOpts)
		if err != nil {
			panic(err)
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			key := strings.TrimLeft(r.URL.Path, "/")
			result := ""
			err = nil

			switch r.Method {
			case "GET":
				row := db.QueryRow(query, key)
				err = row.Scan(&result)
			case "PUT":
				result = "done"
				value, _ := io.ReadAll(r.Body)
				_, err = db.Exec(update, key, string(value[:]))
			default:
				err = fmt.Errorf("unsupported method")
			}

			if err == nil {
				fmt.Fprint(w, result)
			} else {
				http.Error(w, err.Error(), 500)
			}
		})

		listener, err := net.Listen("tcp", api)
		if err != nil {
			panic(err)
		}

		go http.Serve(listener, nil)

		ch := make(chan os.Signal, 32)
		signal.Notify(ch, unix.SIGINT)
		signal.Notify(ch, unix.SIGQUIT)
		signal.Notify(ch, unix.SIGTERM)

		<-ch

		listener.Close()
		db.Close()

		// TODO: handover and close app
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&api, "api", "", "api endpoint to use for serving")
	rootCmd.PersistentFlags().StringVar(&db, "db", "", "database name")
	rootCmd.PersistentFlags().StringVar(&dir, "dir", "", "path for storing sqlite data")
	rootCmd.PersistentFlags().StringVar(&crt, "crt", "", "certificate file for TLS")
	rootCmd.PersistentFlags().StringVar(&key, "key", "", "key file for certificate")

	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&diskMode, "disk-mode", false, "use diskmode")

	join = rootCmd.PersistentFlags().StringSlice("join", []string{}, "List of nodes to form a cluster")

	if err := rootCmd.MarkPersistentFlagRequired("api"); err != nil {
		panic(err)
	}

	if err := rootCmd.MarkPersistentFlagRequired("db"); err != nil {
		panic(err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
