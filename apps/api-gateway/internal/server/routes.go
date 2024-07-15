package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/replay/platform/apps/api-gateway/internal/handlers"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// RouteDeps bundles handlers for v1 registration.
type RouteDeps struct {
	Store     *store.Store
	JWTSecret string
}

// RegisterV1Routes mounts authenticated control-plane routes.
func RegisterV1Routes(r chi.Router, deps RouteDeps) {
	register := &handlers.RegisterHandler{Store: deps.Store}
	login := &handlers.LoginHandler{Store: deps.Store, JWTSecret: deps.JWTSecret}
	orgs := &handlers.OrgsHandler{Store: deps.Store}
	projects := &handlers.ProjectsHandler{Store: deps.Store}
	keys := &handlers.APIKeysHandler{Store: deps.Store}

	session := gwmw.RequireSession(deps.JWTSecret)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Post("/auth/register", register.ServeHTTP)
		v1.Post("/auth/login", login.ServeHTTP)

		v1.Group(func(authed chi.Router) {
			authed.Use(session)
			authed.Get("/orgs/{id}", orgs.GetOrg)
			authed.Put("/orgs/{id}", orgs.UpdateOrg)
			authed.Post("/orgs", orgs.CreateOrg)

			authed.Get("/projects", projects.ListProjects)
			authed.Post("/projects", projects.CreateProject)
			authed.Get("/projects/{id}", projects.GetProject)

			authed.Post("/projects/{id}/api-keys", keys.CreateAPIKey)
			authed.Post("/projects/{id}/api-keys/{keyId}/rotate", keys.RotateAPIKey)
		})
	})
}
