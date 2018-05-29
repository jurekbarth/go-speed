package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	color "github.com/fatih/color"
)

var rsaKeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAoHPLSLplMH1YKQNmp/kxiC8LX4gNTz2RRHtx/tRGJKddwtua
+QTLpQ965rD2winaNOv7zcAtjMsQmepG02vxXoy+cvy3rXgHZKh1N4+A+99TRj5P
cY1eOMEzEdUVnWNKiL5x17fPUOxOc/iJQI1zth7+LI5IbfEGcxgxhONYu3g55hAS
iNVkCTx7aLGDGcuSipH4fVB0M9MkO8CO0EiMeHUYd3ieqrP2FYDsoZtmANrogrUl
ysRdP3I623jr9wF0RshT0qANdnF90aGx6glIzMRpUTCcD3tNoGPc/Ow+Nr2SpPbG
8zrEMriSzwov8edU7HeGur1gssEH5ArwQuYPvQIDAQABAoIBAQCUMP5SyJzGwS3Y
i2SXxUbjIZgefnjUc+ekWWM62fGCzvWBD/S9A5nWdEqtoEn3oFIByOaC7HjlbXOC
xGbvw+Vkzxbi+tfmJlKlvBSu4SJe/q9Z1BjppoicYIv7b1OMTnU7gLGCbCjU87ut
zqFtdnelgFB+9Fae/BpZ2MF7m8KLOZOWE0Am2lgUx5nyLtx9bdYEiIog5pOgATJ1
Yc9HtVoHIfnqZzIrLAoJRndWIqnU2ZeISEK2MyY+9NOd643tIJSa8VL6awexkO/8
VugLWtuk2+zyZBbT5ufZueZGYp1mYRJ48Amh07WY4azk1aLOtNXnPJq6WEQkJGtg
mZa8oYlhAoGBANJWZxEbyCoOGF5ddkwXYIRtU2L1afLxooP1UZJI25bMCrHkzXiS
OFUw+GN/uPPeqCPS5toJDF5eqeyzmOdfpPkeaR4Nyk8eASqpm6DjhEvpr6JhbOKL
+FUVL80+LyxP19QkLm/5vIZBKWZs8J+X+vAO0XQguxn6flZhiEGhAkrZAoGBAMNJ
ACr4gT0UpC4ffafSmGEl5k9Z1unH610jx/csI+pKSjGaQ5gkuvZozhS6mD3VzSBU
VW/v5UpVfNI4GjtTMbkkKOWOO8+Fzfx9scX8mzJBKDMABJTMbEgD7G20JyOvvo4s
hdN+BOy5byzC/h725OA2CU1utWNqE/6g0X7EgHWFAoGAUa1dnoYkRzhr/BDdBBU7
1JDDhbT47G8qhYV4pI6IPtmC+at4om5dU6+NdM2/G2wF7MtT+6zx0Z9+6ryfDpHU
dSx680G1ot1q5I8yMNrIn9Xh7vNYHezuhNOSWWfhV5q1m9pk8fSPYa7iDbUWB1M0
DY4jha3EGgVsk8yR5bJJOpkCgYEAna1vyUJld6AXAHbEyqCsEKS9VQzBDnoxfD7L
0rN9PEtHpM1eDpZ5r0PoQax4CFV9DsGJSpx0kpR7+HD8HTKLT2X274LsoB71twz2
YVoZJXaesq8tA8gbFfq1B88SWyonvjwMwjtaVplTPt0iunW3T6HR2Qeuxdp80nef
L7AR2NECgYB0afYyJ3XDjKeUnf6EhLRFSqswNIj+l6u3k9sRYoyC3VBXSKAT5T+W
r4s0r8xvAg3air34lqnfVxonbsat0TFAO4SrJYwnwnSoIdvesNkrKmr7NzwUtwpx
SiFNZrKeYsaGp978l0vp0f4s7YZqHisVlkV4rxIYUaTByUAyWMuuHw==
-----END RSA PRIVATE KEY-----`)

var rsaCertPEM = []byte(`-----BEGIN CERTIFICATE-----
MIIDiDCCAnACCQC9tV/d9CmTijANBgkqhkiG9w0BAQsFADCBhTELMAkGA1UEBhMC
REUxDzANBgNVBAgMBkJheWVybjEPMA0GA1UEBwwGTXVuaWNoMQswCQYDVQQKDAJW
STEMMAoGA1UECwwDREVWMRYwFAYDVQQDDA1qdXJla2JhcnRoLmRlMSEwHwYJKoZI
hvcNAQkBFhJwb3N0QGp1cmVrYmFydGguZGUwHhcNMTgwNTA0MDU1MTM3WhcNMjgw
NTAxMDU1MTM3WjCBhTELMAkGA1UEBhMCREUxDzANBgNVBAgMBkJheWVybjEPMA0G
A1UEBwwGTXVuaWNoMQswCQYDVQQKDAJWSTEMMAoGA1UECwwDREVWMRYwFAYDVQQD
DA1qdXJla2JhcnRoLmRlMSEwHwYJKoZIhvcNAQkBFhJwb3N0QGp1cmVrYmFydGgu
ZGUwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCgc8tIumUwfVgpA2an
+TGILwtfiA1PPZFEe3H+1EYkp13C25r5BMulD3rmsPbCKdo06/vNwC2MyxCZ6kbT
a/FejL5y/LeteAdkqHU3j4D731NGPk9xjV44wTMR1RWdY0qIvnHXt89Q7E5z+IlA
jXO2Hv4sjkht8QZzGDGE41i7eDnmEBKI1WQJPHtosYMZy5KKkfh9UHQz0yQ7wI7Q
SIx4dRh3eJ6qs/YVgOyhm2YA2uiCtSXKxF0/cjrbeOv3AXRGyFPSoA12cX3RobHq
CUjMxGlRMJwPe02gY9z87D42vZKk9sbzOsQyuJLPCi/x51Tsd4a6vWCywQfkCvBC
5g+9AgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJnEK26Yu1qLQld9knhCa1fWjBBk
NtZWRNxfykkLU+aeA5yQzr+rMRpIazIP5KcJ80eCqXue0h7N9PYarY33WSkvLEBC
8Tc3Hm69vfMguqKWo/oqQlsMSG1o3HrwU7Sw5d/smFpj0SHet6/aIVMQUEaqez/u
3DywGlYIKe64gvtHqCMgXkAFaxm/Er2l85hyPdWAxiR0ejOGd1+psHeEH2rqCMoT
XvUg+Qw5Eep6XDyq43MaNFywBqcZYai1YZnacJ2Cc6fmraKDWPtVMvwh4Jj0LBGb
F5Eyba7Xn8syaOD8U1dhOa8A4Q3rMe0hA3LWI34O6goGbUzpBeXjWfBbnhU=
-----END CERTIFICATE-----`)

type serverConfig struct {
	httpPort     string
	httpsPort    string
	rsaKeyPath   string
	rsaCertPath  string
	rootDir      string
	defaultSpeed float64
}

// helpers

func localIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getFileExtension(s string) string {
	if idx := strings.LastIndex(s, "."); idx != -1 {
		return s[idx:]
	}
	return ""
}

func simulateSpeed(ctx context.Context, timeout float64, serverConfig serverConfig, w http.ResponseWriter, r *http.Request, done chan<- bool) {

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
	typ := mime.TypeByExtension(getFileExtension(r.URL.Path))
	fs := http.FileServer(http.Dir(serverConfig.rootDir))
	switch {
	case strings.HasPrefix(typ, "text/"):
		fs = gziphandler.GzipHandler(fs)
	case typ == "application/xml":
		fs = gziphandler.GzipHandler(fs)
	case typ == "application/javascript":
		fs = gziphandler.GzipHandler(fs)
	case typ == "":
		fs = gziphandler.GzipHandler(fs)
	}
	fs.ServeHTTP(w, r)
	flusher.Flush()
	done <- true
}

func handler(w http.ResponseWriter, r *http.Request, serverConfig serverConfig) {
	// Print body to console, for DEV
	if r.Method != "GET" {
		color.Blue("###############################")
		contentType := r.Header.Get("Content-type")
		color.Blue(r.Method + "-request with content-type: " + contentType)
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			color.Blue(string(body))
		}
		color.Blue("###############################")
	}
	done := make(chan bool)
	closeNotifier, ok := w.(http.CloseNotifier)
	if !ok {
		panic("Expected http.ResponseWriter to be an http.CloseNotifier")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	speed := serverConfig.defaultSpeed
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

	go simulateSpeed(ctx, speed, serverConfig, w, r, done)
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

func run(serverConfig serverConfig) chan error {
	ip := localIP()
	errs := make(chan error)
	// Starting HTTP server
	go func() {
		color.Cyan("🚀 Local HTTP on http://localhost%s", serverConfig.httpPort)
		color.Cyan("🚀 External HTTP on http://%s%s", ip, serverConfig.httpPort)
		if err := http.ListenAndServe(serverConfig.httpPort, nil); err != nil {
			errs <- err
		}

	}()

	// Starting HTTPS server
	go func() {
		// just for ordering
		time.Sleep(200 * time.Millisecond)
		color.Green("🚀 Local HTTPS on https://localhost%s", serverConfig.httpsPort)
		color.Green("🚀 External HTTPS on https://%s%s", ip, serverConfig.httpsPort)
		var cert tls.Certificate
		var certErr error
		if (serverConfig.rsaKeyPath != "") && (serverConfig.rsaCertPath != "") {
			cert, certErr = tls.LoadX509KeyPair(serverConfig.rsaCertPath, serverConfig.rsaKeyPath)
		} else {
			cert, certErr = tls.X509KeyPair(rsaCertPEM, rsaKeyPEM)
		}

		if certErr != nil {
			errs <- certErr
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		server := http.Server{
			// Other options
			Addr:      serverConfig.httpsPort,
			TLSConfig: tlsConfig,
		}

		err := server.ListenAndServeTLS("", "")
		if err != nil {
			errs <- err
		}
	}()

	return errs
}

func main() {
	// httpPort     string
	// httpsPort    string
	// rsaKeyPath   string
	// rsaCertPath  string
	// rootDir      string
	// defaultSpeed float64
	var serverConfig serverConfig
	httpCLF := flag.Int("http", 8080, "http port")
	httpsCLF := flag.Int("https", 8443, "https port")
	rsaKeyPathCLF := flag.String("key", "", "key path")
	rsaCertPathCLF := flag.String("cert", "", "cert path")
	rootDirCLF := flag.String("root", "", "root directory")
	defaultSpeedCLF := flag.Float64("defaultSpeed", 0, "default speed")
	flag.Parse()
	serverConfig.httpPort = ":" + strconv.Itoa(*httpCLF)
	serverConfig.httpsPort = ":" + strconv.Itoa(*httpsCLF)
	serverConfig.rsaKeyPath = *rsaKeyPathCLF
	serverConfig.rsaCertPath = *rsaCertPathCLF
	if *rootDirCLF != "" {
		serverConfig.rootDir = *rootDirCLF + "/"
	}
	serverConfig.defaultSpeed = *defaultSpeedCLF

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, serverConfig)
	})
	errs := run(serverConfig)

	// This will run forever until channel receives error
	select {
	case err := <-errs:
		red := color.New(color.FgRed).PrintfFunc()
		red("🛑 Error: Could not start serving: %s", err)
	}
}
