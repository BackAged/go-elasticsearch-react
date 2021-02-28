package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BackAged/go-elasticsearch-react/backend/config"
	"github.com/BackAged/go-elasticsearch-react/backend/infra"
	"github.com/BackAged/go-elasticsearch-react/backend/repo"
	"github.com/BackAged/go-elasticsearch-react/backend/rest"
	"github.com/BackAged/go-elasticsearch-react/backend/search"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"
)

var serveRestCmd = &cobra.Command{
	Use:   "serve-rest",
	Short: "start a rest server",
	RunE:  serveRest,
}

func serveRest(cmd *cobra.Command, args []string) error {
	cnf := config.GetApp()
	fmt.Printf("loaded config => %+v\n", cnf)

	log.Println("connecting elasticSearch")
	es, err := infra.NewEsClient(cnf.ElasticSearchURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected elasticSearch")

	prdRepo := repo.NewProductRepo(es, repo.RepoNameProduct)
	brndRepo := repo.NewBrandRepo(es, repo.RepoNameBrand)
	shpRepo := repo.NewShopRepo(es, repo.RepoNameShop)

	// initiating services
	svc := search.NewService(prdRepo, brndRepo, shpRepo)

	// inittiating handler
	brndHndlr := rest.NewBrandHandler(svc)
	shpHndlr := rest.NewShopHandler(svc)
	prdHndlr := rest.NewProductHandler(svc)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Use(cors.AllowAll().Handler)

	r.Mount("/api/v1/search/brand", brndHndlr.Router())
	r.Mount("/api/v1/search/shop", shpHndlr.Router())
	r.Mount("/api/v1/search/product", prdHndlr.Router())

	timeout := 30 * time.Second
	srvr := http.Server{
		Addr:              fmt.Sprintf(":%d", cnf.Port),
		Handler:           r,
		ReadTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       timeout,
	}

	errCh := make(chan error)

	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, os.Interrupt}

	graceful := func() error {
		log.Println("Shutting down server gracefully with in", timeout)
		log.Println("To shutdown immediately press again")

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		return srvr.Shutdown(ctx)
	}

	forced := func() error {
		log.Println("Shutting down server forcefully")
		return srvr.Close()
	}

	go func() {
		log.Println("Starting server on", srvr.Addr)
		if err := srvr.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	go func() {
		errCh <- HandleSignals(sigs, graceful, forced)
	}()

	return <-errCh
}

// HandleSignals listen on the registered signals and fires the gracefulHandler for the
// first signal and the forceHandler (if any) for the next this function blocks and
// return any error that returned by any of the handlers first
func HandleSignals(sigs []os.Signal, gracefulHandler, forceHandler func() error) error {
	sigCh := make(chan os.Signal)
	errCh := make(chan error, 1)

	signal.Notify(sigCh, sigs...)
	defer signal.Stop(sigCh)

	grace := true

	select {
	case err := <-errCh:
		return err
	case <-sigCh:
		if grace {
			grace = false
			go func() {
				errCh <- gracefulHandler()
			}()
		} else if forceHandler != nil {
			err := forceHandler()
			errCh <- err
		}
	}

	return nil
}
