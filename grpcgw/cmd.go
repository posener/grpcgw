package grpcgw

import (
	"log"

	"github.com/posener/grpcgw/middleware"

	"github.com/spf13/cobra"
)

var (
	Address           string
	noAPICallsLogging bool
)

const (
	defaultAddress = "localhost:10000"
)

func AddCommands(rootCmd *cobra.Command, service Service) {
	s := NewServer(service)
	serveCmd := newServeCommand(s)
	serveCmd.Flags().StringVarP(&s.Address, "address", "a", defaultAddress, "Listen address")
	serveCmd.Flags().BoolVar(&noAPICallsLogging, "no-api-log", false, "Don't log API calls")
	serveCmd.Flags().StringVar(&s.KeyFile, "key", "certs/server.key", "Private key file")
	serveCmd.Flags().StringVar(&s.CertFile, "cert", "certs/server.pem", "Certificate file")
	rootCmd.AddCommand(serveCmd)

	SendCmd.Flags().StringVarP(&Address, "address", "a", defaultAddress, "Listen address")
	rootCmd.AddCommand(SendCmd)
}

// serveCmd represents the serve command
func newServeCommand(s *server) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run grpcgw example server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if !noAPICallsLogging {
				s.Middleware = s.Middleware.Append(middleware.APILoggerMiddleware)
			}
			Serve(s)
		},
	}
}

// sendCmd represents the send command
var SendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a GRPC message",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("Must specify which message to send")
	},
}
