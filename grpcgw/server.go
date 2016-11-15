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

func Serve(s *server) {

	var err error

	log.Print("Loading server certificates...")

	s.checkSecure()

	ctx := context.Background()
	certPool := s.createCertPool()

	log.Print("Starting secure GRPC server")
	options := []grpc.ServerOption{grpc.Creds(credentials.NewClientTLSFromCert(certPool, s.Address))}
	grpcServer := grpc.NewServer(options...)

	log.Print("Initializing static pages server...")
	gwmux := runtime.NewServeMux()

	log.Print("Registering server endpoints...")


	dialCreds := credentials.NewTLS(&tls.Config{ServerName: s.Address, RootCAs: certPool})
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(dialCreds)}
	err = register(s, ctx, grpcServer, gwmux, s.Address, dialOptions)
	if err != nil {
		log.Panicf("Failed registering: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	addSwaggerUIHandlers(s, mux)

	log.Print("Starting to listen...")
	conn, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Panic(err)
	}

	handler := s.Middleware.Then(grpcHandlerFunc(s, grpcServer, mux))

	certificates := []tls.Certificate{s.createCertificate()}
	tlsConfig := tls.Config{Certificates: certificates, NextProtos: []string{"h2"}}
	srv := &http.Server{Addr: s.Address, Handler: handler, TLSConfig: &tlsConfig}

	log.Printf("Grpc is ready on: %s", s.Address)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		log.Panicf("ListenAndServe failed: %s", err)
	}

	return
}

func register(s Service, ctx context.Context, grpcServer *grpc.Server, gwmux *runtime.ServeMux, grpcEndpointAddr string, opts []grpc.DialOption) error {
	s.RegisterGRPC(grpcServer)
	err := s.RegisterGatewayEndpoints(ctx, gwmux, grpcEndpointAddr, opts)
	return err
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(s *server, grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := getHandler(r, grpcServer, otherHandler)
		handler.ServeHTTP(w, r)
	})
}

func getHandler(r *http.Request, grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		return grpcServer
	} else {
		return otherHandler
	}
}

func addSwaggerUIHandlers(s *server, mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")
    uiServer := http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, uiServer))

	path, err := filepath.Abs(s.SwaggersPath)
	if err != nil {
		log.Panic("Failed calculating absoulute path of swagger directory")
	}
	jsonServer := http.FileServer(http.Dir(path))
	prefix = "/swaggers/"
	mux.Handle(prefix, http.StripPrefix(prefix, jsonServer))
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