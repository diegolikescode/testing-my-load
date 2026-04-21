package pkg

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/phuslu/log"
)

const pprofServerPort = ":6060"

func StartPprofServer() {
	log.Info().Msg("starting Pprof server on port=" + pprofServerPort)
	http.ListenAndServe(pprofServerPort, nil)
}
