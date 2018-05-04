package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)
import color "github.com/fatih/color"

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func simulateSpeed(ctx context.Context, timeout float64, w http.ResponseWriter, r *http.Request, done chan<- bool) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("Expected http.ResponseWriter to be an http.Flusher")
	}
	multiplier := time.Duration(int32(timeout * 1000))
	select {
	case <-time.After(time.Millisecond * multiplier):
	case <-ctx.Done():
		done <- false // Cancel job.
		return
	}

	enableCors(&w)

	http.ServeFile(w, r, r.URL.Path[1:])
	flusher.Flush()
	done <- true
}

func handler(w http.ResponseWriter, r *http.Request) {
	done := make(chan bool)
	closeNotifier, ok := w.(http.CloseNotifier)
	if !ok {
		panic("Expected http.ResponseWriter to be an http.CloseNotifier")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	speed := float64(0)
	speedQueryParam, ok := r.URL.Query()["speed"]
	if ok || len(speedQueryParam) > 1 {
		speedFactor, err := strconv.ParseFloat(speedQueryParam[0], 64)
		if err == nil {
			yellow := color.New(color.FgHiYellow).PrintfFunc()
			path := r.URL
			yellow("🐌: %s by %.2f second\n", path, speedFactor)

			speed = speedFactor
		}
	}

	go simulateSpeed(ctx, speed, w, r, done)
	select {
	case <-done:
	case <-time.After(time.Second * 60):
		cancel()
		if !<-done {
			fmt.Fprint(w, "Server is busy.")
		}
	case <-closeNotifier.CloseNotify():
		cancel()
		fmt.Println("Client has disconnected.")
		<-done
	}
}

func localIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func Run(port string, sslPort string, ssl map[string]string) chan error {
	ip := localIP()
	errs := make(chan error)
	// Starting HTTP server
	go func() {
		color.Cyan("🚀 Local HTTP on http://localhost%s", port)
		color.Cyan("🚀 External HTTP on http://%s%s", ip, port)
		if err := http.ListenAndServe(port, nil); err != nil {
			errs <- err
		}

	}()

	// Starting HTTPS server
	go func() {
		// just for ordering
		time.Sleep(200 * time.Millisecond)
		color.Green("🚀 Local HTTPS on https://localhost%s", sslPort)
		color.Green("🚀 External HTTPS on https://%s%s", ip, sslPort)
		if err := http.ListenAndServeTLS(sslPort, ssl["cert"], ssl["key"], nil); err != nil {
			errs <- err
		}
	}()

	return errs
}

func main() {
	http.HandleFunc("/", handler)

	httpCLF := flag.Int("http", 8080, "http port")
	httpsCLF := flag.Int("https", 8443, "https port")
	flag.Parse()

	httpPort := ":" + strconv.Itoa(*httpCLF)
	httpsPort := ":" + strconv.Itoa(*httpsCLF)

	errs := Run(httpPort, httpsPort, map[string]string{
		"cert": "server.crt",
		"key":  "server.key",
	})

	// This will run forever until channel receives error
	select {
	case err := <-errs:
		red := color.New(color.FgRed).PrintfFunc()
		red("🛑 Error: Could not start serving: %s", err)
	}
}
