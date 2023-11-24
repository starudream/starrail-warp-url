package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/core/v2/utils/osutil"
	"github.com/starudream/go-lib/core/v2/utils/signalutil"

	"github.com/starudream/starrail-warp-url/web"
)

func startServer(data string) {
	mux := http.NewServeMux()
	mux.Handle("/", web.Handler(data))

	ln := osutil.Must1(net.Listen("tcp", ":0"))

	addr := fmt.Sprintf("http://%s:%d", localIP(), ln.Addr().(*net.TCPAddr).Port)
	slog.Info("server started at %s", addr)
	fmt.Printf("\n" + genQRCode(addr) + "\n")

	hs := &http.Server{Handler: mux}

	go func() {
		<-time.After(10 * time.Minute)
		signalutil.Cancel()
	}()
	go func() {
		<-signalutil.Defer(func() { _ = hs.Shutdown(context.Background()) }).Done()
	}()
	go func() {
		gh.Silently(hs.Serve(ln))
	}()
}
