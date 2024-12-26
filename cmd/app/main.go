package main

import (
	"context"
	"embed"
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sankalp-r/url-shortner/internal/handlers"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

var (
	domain      = flag.String("domain", "", "your ZITADEL instance domain (in the form: https://<instance>.zitadel.cloud or https://<yourdomain>)")
	key         = flag.String("key", "", "encryption key")
	clientID    = flag.String("clientID", "", "clientID provided by ZITADEL")
	redirectURI = flag.String("redirectURI", "", "redirectURI registered at ZITADEL")
	keyFile     = flag.String("keyFile", "", "path to your key.json")
	//go:embed "templates/*.html"
	templates embed.FS
)

func main() {
	flag.Parse()
	ctx := context.Background()

	t, err := template.New("").ParseFS(templates, "templates/*.html")
	if err != nil {
		slog.Error("unable to parse template", "error", err)
		os.Exit(1)
	}

	// initialize authenticator
	authN, err := authentication.New(ctx, zitadel.New(*domain), *key,
		oidc.DefaultAuthentication(*clientID, *redirectURI, *key),
	)
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	// initialize authorizer
	authZ, err := authorization.New(ctx, zitadel.New(*domain), oauth.DefaultAuthorization(*keyFile))
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	router := http.NewServeMux()

	// initialize handler
	handler := handlers.New(handlers.WithAuthenticator(authN), handlers.WithAuthorizer(authZ), handlers.WithATemplate(t))
	handler.RegisterRoutes(router)

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("error creating server")
		}
	}()

	// Block until a signal is received
	<-stop

	// Create a context with a timeout for the shutdown
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exiting")

}
