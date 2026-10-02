// gdrive-auth runs locally once; it is not part of the production HTTP server.
package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	id, secret := os.Getenv("BACKUP_GDRIVE_CLIENT_ID"), os.Getenv("BACKUP_GDRIVE_CLIENT_SECRET")
	if id == "" || secret == "" {
		return fmt.Errorf("set BACKUP_GDRIVE_CLIENT_ID and BACKUP_GDRIVE_CLIENT_SECRET for a Desktop OAuth client")
	}
	parent, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(parent, 10*time.Minute)
	defer cancel()
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("could not bind local OAuth callback")
	}
	cfg := &oauth2.Config{ClientID: id, ClientSecret: secret, Endpoint: google.Endpoint, RedirectURL: "http://" + listener.Addr().String() + "/callback", Scopes: []string{"https://www.googleapis.com/auth/drive"}}
	state, verifier := oauth2.GenerateVerifier(), oauth2.GenerateVerifier()
	result := make(chan *oauth2.Token, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != http.MethodGet || subtle.ConstantTimeCompare([]byte(r.URL.Query().Get("state")), []byte(state)) != 1 {
			http.Error(w, "Invalid authorization state", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("error") != "" || r.URL.Query().Get("code") == "" {
			http.Error(w, "Authorization was not granted. Restart the helper to try again.", http.StatusBadRequest)
			return
		}
		exchangeCtx, done := context.WithTimeout(ctx, 30*time.Second)
		defer done()
		token, err := cfg.Exchange(exchangeCtx, r.URL.Query().Get("code"), oauth2.VerifierOption(verifier))
		if err != nil || token.RefreshToken == "" {
			http.Error(w, "Could not obtain offline authorization. Restart the helper and grant consent.", http.StatusBadRequest)
			return
		}
		select {
		case result <- token:
			_, _ = fmt.Fprintln(w, "Authorization complete. Return to your terminal and close this tab.")
		default:
			http.Error(w, "Already authorized", http.StatusConflict)
		}
	})
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 40 * time.Second}
	defer func() { _ = srv.Close() }()
	go func() { _ = srv.Serve(listener) }()
	fmt.Println("Open this URL in your browser and authorize the Google account owning the backup folder:")
	fmt.Println(cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce, oauth2.S256ChallengeOption(verifier)))
	select {
	case <-ctx.Done():
		return fmt.Errorf("authorization cancelled or timed out")
	case token := <-result:
		fmt.Println("Store the following secret in your server environment; do not commit or share it:")
		fmt.Printf("BACKUP_GDRIVE_REFRESH_TOKEN=%s\n", token.RefreshToken)
		return nil
	}
}
