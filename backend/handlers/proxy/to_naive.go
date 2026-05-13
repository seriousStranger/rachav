//nolint:all
package proxy

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// WARNING
// ai generated code
// todo: rewrite

func toNaive(w http.ResponseWriter, r *http.Request, tr *http.Transport, upstreamAddr string) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Hour)
	defer cancel()

	rc := http.NewResponseController(w)

	pr, pw := io.Pipe()
	outReq, _ := http.NewRequestWithContext(ctx, http.MethodConnect, "http://"+upstreamAddr, pr)
	outReq.Host = r.Host

	resp, err := tr.RoundTrip(outReq)
	if err != nil {
		slog.Error(
			"Upstream error (%s): %v",
			slog.String("Host", r.Host),
			slog.String("Error", err.Error()),
		)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	rc.Flush()

	toUpstream := pw
	toClient := w

	errChan := make(chan error, 2)

	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, rerr := r.Body.Read(buf)
			if n > 0 {
				_, werr := toUpstream.Write(buf[:n])
				if werr != nil {
					errChan <- werr
					return
				}
			}
			if rerr != nil {
				pw.Close()
				if rerr != io.EOF {
					errChan <- rerr
				} else {
					errChan <- nil
				}
				return
			}
		}
	}()

	go func() {
		buf := make([]byte, 32*1024)
		for {
			n, rerr := resp.Body.Read(buf)
			if n > 0 {
				_, werr := toClient.Write(buf[:n])
				if werr != nil {
					errChan <- werr
					return
				}
				_ = rc.Flush()
			}
			if rerr != nil {
				if rerr != io.EOF {
					errChan <- rerr
				} else {
					errChan <- nil
				}
				return
			}
		}
	}()

	err = <-errChan
	if err != nil && err != context.Canceled {
		slog.Error("Stream finished with error", slog.String("Error:", err.Error()))
	}
}
