package healthcheck
import (
  "net/http"
  "fmt"
  "context"
  "time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "OK")
}

func StartHealthCheck(listenAddr string, uri string) (*http.Server) {
    log.Tracef("Start healthcheck at %v [%v]", uri, listenAddr)
     srv:=&http.Server{
	Addr:           listenAddr,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
    }
      http.HandleFunc(uri, healthHandler)
go func() {
    // log.Warningf(" !!!")

      if err := srv.ListenAndServe(); err != http.ErrServerClosed {
          log.Fatalf("%w %v", ErrServerListenError, err)
      }
      log.Warningf(" !!!")

  }()
  return srv
}


func WaitForShutdown(ctx context.Context, srv *http.Server) {
    srv.Shutdown(ctx)
}
