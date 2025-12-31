package security

// OrgFixture represents a tenant for isolation tests.
type OrgFixture struct {
	ID    string
	Name  string
	Token string
}

// TwoOrgs returns distinct org fixtures for cross-tenant matrix tests.
func TwoOrgs() (OrgA, OrgB OrgFixture) {
	return OrgFixture{
			ID:    "org-a-00000000-0000-0000-0000-000000000001",
			Name:  "Tenant Alpha",
			Token: "token-org-a",
		}, OrgFixture{
			ID:    "org-b-00000000-0000-0000-0000-000000000002",
			Name:  "Tenant Beta",
			Token: "token-org-b",
		}
}

// IsolationCase describes an API action and expected denial when crossing tenants.
type IsolationCase struct {
	Name           string
	Method         string
	Path           string
	Body           string
	ExpectStatus   int
	ExpectDenied   bool
}
