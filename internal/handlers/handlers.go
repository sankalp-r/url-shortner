package handlers

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	storage "github.com/sankalp-r/url-shortner/pkg/storage"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	oidc2 "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
)

type Handler struct {
	store         storage.Store
	authenticator *authentication.Authenticator[*oidc2.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo]]
	authorizer    *authorization.Authorizer[*oauth.IntrospectionContext]
	template      *template.Template
}

// Option function that allows injecting dependencies into the handler.
type Option func(*options)

type options struct {
	authenticator *authentication.Authenticator[*oidc2.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo]]
	authorizer    *authorization.Authorizer[*oauth.IntrospectionContext]
	template      *template.Template
}

// WithAuthenticator returns an Option injected with Authenticator.
func WithAuthenticator(authn *authentication.Authenticator[*oidc2.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo]]) Option {
	return func(o *options) {
		o.authenticator = authn
	}
}

// WithAuthorizer returns an Option injected with Authorizer.
func WithAuthorizer(authz *authorization.Authorizer[*oauth.IntrospectionContext]) Option {
	return func(o *options) {
		o.authorizer = authz
	}
}

// WithTemplate returns an Option injected with template.
func WithATemplate(template *template.Template) Option {
	return func(o *options) {
		o.template = template
	}
}

func New(opts ...Option) *Handler {
	defaultOpts := &options{}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	return &Handler{
		store:         storage.NewStore(),
		authenticator: defaultOpts.authenticator,
		authorizer:    defaultOpts.authorizer,
		template:      defaultOpts.template,
	}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("error decoding request", "error", err)
		http.Error(w, "error generating short URL", http.StatusBadRequest)
		return
	}

	shortURLCode, err := h.store.Create(req.URL)
	if err != nil {
		slog.Error("error creating short-url code", "error", err)
		http.Error(w, "error generating short URL", http.StatusBadRequest)
		return
	}
	res := ShortenResponse{ShortURL: shortURLCode}

	json.NewEncoder(w).Encode(res)
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	originalURL, err := h.store.Get(shortURL)
	if err != nil {
		slog.Error("error getting url", "error", err)
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	if h.authenticator != nil {
		router.Handle("/auth/", h.authenticator)

		authenticationMiddleware := authentication.Middleware(h.authenticator)

		router.Handle("/profile", authenticationMiddleware.RequireAuthentication()(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			authCtx := authenticationMiddleware.Context(req.Context())
			_, err := json.MarshalIndent(authCtx.UserInfo, "", "	")
			if err != nil {
				slog.Error("error marshalling profile response", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data := struct {
				Username string
				Token    string
			}{
				Username: authCtx.UserInfo.PreferredUsername,
				Token:    authCtx.GetTokens().AccessToken,
			}
			err = h.template.ExecuteTemplate(w, "profile.html", data)
			if err != nil {
				slog.Error("error writing profile response", "error", err)
			}
		})))

		router.Handle("/", authenticationMiddleware.CheckAuthentication()(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if authentication.IsAuthenticated(req.Context()) {
				http.Redirect(w, req, "/profile", http.StatusFound)
				return
			}
			err := h.template.ExecuteTemplate(w, "home.html", nil)
			if err != nil {
				slog.Error("error writing home page response", "error", err)
			}
		})))
	}

	r := mux.NewRouter().UseEncodedPath()

	if h.authorizer != nil {
		authorizationMiddleware := middleware.New(h.authorizer)
		r.Handle("/short", authorizationMiddleware.RequireAuthorization()(http.HandlerFunc(h.ShortenURL))).Methods("POST").Name("Create short URL")
	} else {
		r.Handle("/short", http.HandlerFunc(h.ShortenURL)).Methods("POST").Name("Create short URL")
	}

	r.Handle("/{shortURL}", http.HandlerFunc(h.RedirectURL)).Methods("GET").Name("Redirect short URL")

	router.Handle("/v1/", http.StripPrefix("/v1", r))

}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url_code"`
}
