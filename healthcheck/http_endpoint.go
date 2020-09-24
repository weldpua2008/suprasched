package healthcheck

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func StartHealthCheck(listenAddr string, uri string) *http.Server {
	log.Tracef("Start healthcheck at %v [%v]", uri, listenAddr)
	srv := &http.Server{
		Addr:         listenAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	http.HandleFunc(uri, healthHandler)
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%w %v", ErrServerListenError, err)
		}
	}()
	return srv
}

func WaitForShutdown(ctx context.Context, srv *http.Server) {
	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Warningf("Shutdown %v", err)

	}
}
