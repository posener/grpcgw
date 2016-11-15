package grpcgw

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	assetfs "github.com/philips/go-bindata-assetfs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"crypto/x509"

	"github.com/justinas/alice"
	"io/ioutil"
)

type server struct {
	Service
	Address    string
	Middleware alice.Chain
	Swagger    map[string]string
	KeyFile string
	CertFile string
}

func NewServer(service Service) *server {
	return &server{Service: service, Middleware: alice.Chain{}}
}

func Serve(s *server) {

	log.Print("Loading server certificates...")

	keyPair, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
	if err != nil {
		log.Panicf("Failed loading key-pair from files %s, %s: %s", s.CertFile, s.KeyFile, err)
	}

	cert, err := ioutil.ReadFile(s.CertFile)
	if err != nil {
		log.Panicf("Failed reading cert from %s: %s", s.CertFile, err)
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(cert)
	if !ok {
		log.Panic("Could not add cert to certPool")
	}

	log.Print("Initializing grpc server...")
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewClientTLSFromCert(certPool, s.Address)),
	)
	ctx := context.Background()

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: s.Address,
		RootCAs:    certPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	log.Print("Initializing static pages server...")
	mux := http.NewServeMux()
	for name, body := range s.Swagger {
		mux.HandleFunc(fmt.Sprintf("/%s.json", name), func(w http.ResponseWriter, req *http.Request) {
			io.Copy(w, strings.NewReader(body))
		})
	}

	gwmux := runtime.NewServeMux()

	log.Print("Registering server endpoints...")
	err = register(s, ctx, grpcServer, gwmux, s.Address, dopts)
	if err != nil {
		log.Printf("serve: %v", err)
		return
	}

	mux.Handle("/", gwmux)
	serveSwagger(mux)

	log.Print("Starting to listen...")
	conn, err := net.Listen("tcp", s.Address)
	if err != nil {
		panic(err)
	}

	handler := s.Middleware.Then(grpcHandlerFunc(s, grpcServer, mux))

	srv := &http.Server{
		Addr:    s.Address,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			NextProtos:   []string{"h2"},
		},
	}

	log.Printf("Grpc is ready on: %s", s.Address)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		log.Fatal("ListenAndServe failed: ", err)
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

func serveSwagger(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in swagger-ui/ dir on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
