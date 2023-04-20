package cmd

import (
	"context"
	"github.com/masinger/courl/internal/authprovider"
	"github.com/masinger/courl/internal/clientconfig"
	"github.com/masinger/courl/internal/config"
	"github.com/masinger/courl/internal/request"
	"github.com/masinger/courl/internal/response"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

var authConfig = authprovider.Config{}
var clientRequest = request.Request{
	Method: http.MethodGet,
}
var clientConfig = clientconfig.Config{
	FollowRedirects:   false,
	IgnoreMissingAuth: false,
}
var clientResponse = response.Response{}
var verbose = false

var rootCmd = &cobra.Command{
	Use:   "courl <URL>",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetOutput(os.Stderr)
		log.SetFormatter(&log.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
		clientResponse.Verbose = verbose
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		var err error
		tokenSource := authConfig.TokenSource(context.Background())
		if tokenSource == nil {
			cfg, err := config.ReadDefault()
			if err != nil {
				return err
			}
			err = cfg.Apply(&authConfig, url)
			if err != nil {
				return err
			}
			tokenSource = authConfig.TokenSource(context.Background())

		}

		client, err := clientConfig.Create(tokenSource)
		if err != nil {
			return err
		}

		req, err := clientRequest.CreateRequest(url)
		if err != nil {
			return err
		}

		return clientResponse.Handle(client.Do(req))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Auth
	rootCmd.PersistentFlags().StringVar(&authConfig.TokenURL, "token-url", "", "Set the token url.")
	rootCmd.PersistentFlags().StringVar(&authConfig.ClientID, "client-id", "", "Set the client id.")
	rootCmd.PersistentFlags().StringVar(&authConfig.ClientSecret, "client-secret", "", "Set the client secret.")

	// Client behavior
	rootCmd.PersistentFlags().BoolVarP(&clientConfig.FollowRedirects, "location", "L", false, "Set if redirects should be followed.")
	rootCmd.PersistentFlags().BoolVarP(&clientConfig.InsecureSkipVerify, "insecure", "k", false, "This option makes courl skip the certificate verification step and proceed without checking")
	rootCmd.PersistentFlags().BoolVar(&clientConfig.IgnoreMissingAuth, "ignore-missing-auth", false, "If set, courl will send the request even if it can not perform authentication.")

	// Request
	rootCmd.PersistentFlags().StringArrayVarP(&clientRequest.Headers, "header", "H", []string{}, "Set headers in form of 'Header-Name: value'")
	rootCmd.PersistentFlags().StringVarP(&clientRequest.Method, "request", "X", "", "Set request method.")
	rootCmd.PersistentFlags().StringArrayVarP(&clientRequest.LiteralDataExpressions, "data", "d", []string{}, "Sends the specified data in a form POST request to the HTTP server (application/x-www-form-urlencoded).")
	rootCmd.PersistentFlags().StringArrayVar(&clientRequest.UnencodedDataExpressions, "data-urlencode", []string{}, "This posts data, similar to the other -d, --data options with the exception that this performs URL-encoding.")
	rootCmd.PersistentFlags().StringArrayVar(&clientRequest.DataBinaryExpressions, "data-binary", []string{}, "his posts data exactly as specified with no extra processing whatsoever.")

	// Response
	rootCmd.PersistentFlags().StringVarP(&clientResponse.OutputPath, "output", "o", "", "Write output to <file> instead of stdout")
	rootCmd.PersistentFlags().BoolVarP(&clientResponse.ErrorOnFail, "fail", "f", false, "Fail fast with no output at all on server errors.")
	rootCmd.PersistentFlags().BoolVar(&clientResponse.AllowBinaryToStdout, "binary-stdout", false, "Allow binary outputs to be written to stdout.")

	// Courl
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Set to enable verbose output.")

}
