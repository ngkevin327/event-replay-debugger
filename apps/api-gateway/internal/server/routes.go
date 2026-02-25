package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/replay/platform/apps/api-gateway/internal/audit"
	"github.com/replay/platform/apps/api-gateway/internal/handlers"
	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
	"github.com/replay/platform/apps/api-gateway/internal/store"
)

// RouteDeps bundles handlers for v1 registration.
type RouteDeps struct {
	Store     *store.Store
	Audit     *audit.Logger
	JWTSecret string
}

// RegisterV1Routes mounts authenticated control-plane routes.
func RegisterV1Routes(r chi.Router, deps RouteDeps) {
	register := &handlers.RegisterHandler{Store: deps.Store}
	login := &handlers.LoginHandler{Store: deps.Store, JWTSecret: deps.JWTSecret}
	orgs := &handlers.OrgsHandler{Store: deps.Store}
	projects := &handlers.ProjectsHandler{Store: deps.Store}
	keys := &handlers.APIKeysHandler{Store: deps.Store, Audit: deps.Audit}
	members := &handlers.MembersHandler{Store: deps.Store, Audit: deps.Audit}
	agents := &handlers.AgentsHandler{Store: deps.Store}
	incidents := &handlers.IncidentsHandler{Store: deps.Store}
	timeline := &handlers.TimelineHandler{Loader: handlers.ResolveTimelineLoader()}
	graph := &handlers.GraphHandler{Loader: handlers.ResolveGraphLoader()}
	replays := &handlers.ReplaysHandler{Store: deps.Store}
	snapshots := &handlers.SnapshotsHandler{Store: deps.Store}
	notifications := &handlers.NotificationsHandler{Store: deps.Store}
	exportH := &handlers.ExportHandler{Store: deps.Store}
	share := &handlers.ShareHandler{Store: deps.Store}

	session := gwmw.RequireSession(deps.JWTSecret)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Post("/auth/register", register.ServeHTTP)
		v1.Post("/auth/login", login.ServeHTTP)

		admin := gwmw.RequireRole(deps.Store, store.RoleAdmin)
		member := gwmw.RequireRole(deps.Store, store.RoleMember)

		v1.Group(func(authed chi.Router) {
			authed.Use(session)
			authed.Get("/orgs/{id}", orgs.GetOrg)
			authed.With(admin).Put("/orgs/{id}", orgs.UpdateOrg)
			authed.With(member).Post("/orgs", orgs.CreateOrg)

			authed.Get("/projects", projects.ListProjects)
			authed.With(member).Post("/projects", projects.CreateProject)
			authed.Get("/projects/{id}", projects.GetProject)
			authed.Get("/projects/{id}/notification-preferences", notifications.GetPrefs)
			authed.With(admin).Put("/projects/{id}/notification-preferences", notifications.PutPrefs)

			authed.With(member).Post("/projects/{projectId}/incidents", incidents.CreateIncident)
			authed.Get("/projects/{projectId}/incidents", incidents.ListIncidents)
			authed.Get("/incidents/{incidentId}", incidents.GetIncident)
			authed.Get("/incidents/{incidentId}/export", exportH.ExportIncidentSummary)
			authed.With(member).Post("/incidents/{incidentId}/share-tokens", share.CreateShareToken)
			authed.Get("/incidents/{incidentId}/timeline", timeline.GetTimeline)
			authed.Get("/incidents/{incidentId}/graph", graph.GetGraph)
			authed.With(member).Post("/incidents/{incidentId}/replays", replays.CreateReplay)
			authed.Get("/replays/{replayId}", replays.GetReplay)
			authed.Delete("/replays/{replayId}", replays.CancelReplay)
			authed.With(member).Post("/incidents/{incidentId}/snapshots", snapshots.IngestSnapshot)

			authed.With(admin).Post("/projects/{id}/api-keys", keys.CreateAPIKey)
			authed.With(admin).Post("/projects/{id}/api-keys/{keyId}/rotate", keys.RotateAPIKey)

			authed.With(admin).Post("/orgs/{id}/members", members.InviteMember)
			authed.With(admin).Put("/orgs/{id}/members/{userId}", members.UpdateRole)
			authed.With(admin).Delete("/orgs/{id}/members/{userId}", members.RemoveMember)

			authed.With(member).Post("/agents/register", agents.RegisterAgent)
			authed.Post("/agents/heartbeat", agents.HeartbeatAgent)
			authed.Get("/agents", agents.ListAgents)
		})

		v1.Route("/shared", func(shared chi.Router) {
			shared.Use(gwmw.ShareTokenAuth(handlers.LookupShareToken))
			shared.Use(gwmw.RequireReadOnly)
			shared.Get("/incidents/{incidentId}", share.GetSharedIncident)
		})
	})
}
