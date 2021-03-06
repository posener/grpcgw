package grpcgw

import (
	"crypto/tls"
	"log"
	"mime"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"crypto/x509"

	"github.com/justinas/alice"
	"io/ioutil"
	"path/filepath"
	"github.com/philips/go-bindata-assetfs"
)

type server struct {
	Service
	Address    string
	Middleware alice.Chain
	Swagger    map[string]string
	KeyFile    string
	CertFile   string
	SwaggersPath string
}

func NewServer(service Service) *server {
	return &server{Service: service, Middleware: alice.Chain{}}
}

func Serve(s *server, ctx context.Context) {
	s.checkSecure()

	var err error
	certPool := s.createCertPool()

	mainMux := http.NewServeMux()
	grpcHandler := createGrpcHandler(s, certPool)

	prefix := "/swagger-ui/"
	mainMux.Handle(prefix, http.StripPrefix(prefix, handleSwaggerUI()))
	prefix = "/swaggers/"
	mainMux.Handle(prefix, http.StripPrefix(prefix, handleSwaggerJson(s.SwaggersPath)))

	mainMux.Handle("/", createGateway(s, ctx, certPool))

	conn, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Panic(err)
	}

	mainHandler := s.Middleware.Append(gatewayMiddleware(grpcHandler)).Then(mainMux)

	certificates := []tls.Certificate{s.createCertificate()}
	tlsConfig := tls.Config{Certificates: certificates, NextProtos: []string{"h2"}}
	srv := &http.Server{Addr: s.Address, Handler: mainHandler, TLSConfig: &tlsConfig}

	listener := tls.NewListener(conn, srv.TLSConfig)
	err = srv.Serve(listener)

	if err != nil {
		log.Panicf("ListenAndServe failed: %s", err)
	}

	log.Printf("Grpc is ready on: %s", s.Address)
	go func() {
		<- ctx.Done()
		grpcHandler.Stop()
		conn.Close()
		listener.Close()
	}()

	return
}

func createGrpcHandler(s *server, certPool *x509.CertPool) *grpc.Server {
	options := []grpc.ServerOption{grpc.Creds(credentials.NewClientTLSFromCert(certPool, s.Address))}
	grpcHandler := grpc.NewServer(options...)
	s.RegisterGRPC(grpcHandler)
	return grpcHandler
}

func createGateway(s *server, ctx context.Context, certPool *x509.CertPool) http.Handler {
	gwMux := runtime.NewServeMux()

	dialCreds := credentials.NewTLS(&tls.Config{ServerName: s.Address, RootCAs: certPool})
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(dialCreds)}
	err := s.RegisterGatewayEndpoints(ctx, gwMux, s.Address, dialOptions)
	if err != nil {
		log.Panicf("Failed registering: %v", err)
	}
	return gwMux
}

// construct a gateway middleware.
// According to the request it chooses if to use the gateway handler
// or if to pass it on.
func gatewayMiddleware(grpcHandler *grpc.Server) alice.Constructor {
	return func (handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcHandler.ServeHTTP(w, r)
			} else {
				handler.ServeHTTP(w, r)
			}
		})
	}
}

func handleSwaggerUI() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	return http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "swagger-ui",
	})
}

func handleSwaggerJson(swaggersPath string) http.Handler {
	path, err := filepath.Abs(swaggersPath)
	if err != nil {
		log.Panic("Failed calculating absoulute path of swagger directory")
	}
	return http.FileServer(http.Dir(path))
}

func (s *server)checkSecure() {
	if s.CertFile == "" || s.KeyFile == "" {
		log.Fatal("Must provide a key and certificate to run server")
	}
}

func (s *server)createCertPool() *x509.CertPool{
	cert, err := ioutil.ReadFile(s.CertFile)
	if err != nil {
		log.Fatalf("Failed reading cert from %s: %s", s.CertFile, err)
	}
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(cert)
	if !ok {
		log.Panic("Could not add cert to certPool")
	}
	return certPool
}

func (s *server)createCertificate() tls.Certificate {
	keyPair, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
	if err != nil {
		log.Panicf("Failed loading key-pair from files %s, %s: %s", s.CertFile, s.KeyFile, err)
	}
	return keyPair
}