package restapi

import (
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type RestAPI struct {
	port        string
	siteHandler http.Handler
}

const defPort = "8083"

func Init() *RestAPI {
	port := os.Getenv("PORT")
	if port == "" {
		port = defPort
	}

	restApi := &RestAPI{
		port: port,
	}

	router := mux.NewRouter()

	// ATTENTION: Next handlers without api prefix must be before ROOT index handler!
	// Pprof
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	router.HandleFunc("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	router.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	router.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	router.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	router.HandleFunc("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
	router.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler())

	restApi.siteHandler = router

	return restApi
}

func (restApi *RestAPI) Start(logger *zap.SugaredLogger) {
	logger.Infof("Start RestAPI at :%v", restApi.port)
	go func() {
		logger.Fatalf("ListenAndServe: %v", http.ListenAndServe(":"+restApi.port, restApi.siteHandler))
	}()
}
