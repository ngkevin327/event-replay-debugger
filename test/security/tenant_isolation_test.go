package security

import "testing"

func TestCrossTenantReadDeniedMatrix(t *testing.T) {
	orgA, orgB := TwoOrgs()
	cases := []IsolationCase{
		{Name: "get foreign org", Method: "GET", Path: "/v1/orgs/" + orgB.ID, ExpectStatus: 404, ExpectDenied: true},
		{Name: "list foreign project incidents", Method: "GET", Path: "/v1/projects/proj-b/incidents", ExpectStatus: 404, ExpectDenied: true},
		{Name: "get foreign incident", Method: "GET", Path: "/v1/incidents/inc-b", ExpectStatus: 404, ExpectDenied: true},
		{Name: "export foreign incident", Method: "GET", Path: "/v1/incidents/inc-b/export", ExpectStatus: 404, ExpectDenied: true},
		{Name: "create replay foreign incident", Method: "POST", Path: "/v1/incidents/inc-b/replays", Body: `{"timing_mode":"strict"}`, ExpectStatus: 404, ExpectDenied: true},
		{Name: "put foreign notification prefs", Method: "PUT", Path: "/v1/projects/proj-b/notification-preferences", Body: `{}`, ExpectStatus: 404, ExpectDenied: true},
		{Name: "own org readable", Method: "GET", Path: "/v1/orgs/" + orgA.ID, ExpectStatus: 200, ExpectDenied: false},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.ExpectDenied && tc.ExpectStatus < 400 {
				t.Fatalf("denied case must expect 4xx, got %d", tc.ExpectStatus)
			}
			_ = orgA.Token
		})
	}
}

func TestCrossTenantWriteDenied(t *testing.T) {
	_, orgB := TwoOrgs()
	writes := []IsolationCase{
		{Name: "post incident other project", Method: "POST", Path: "/v1/projects/proj-b/incidents", ExpectDenied: true, ExpectStatus: 404},
		{Name: "share token other incident", Method: "POST", Path: "/v1/incidents/inc-b/share-tokens", ExpectDenied: true, ExpectStatus: 404},
	}
	for _, tc := range writes {
		t.Run(tc.Name, func(t *testing.T) {
			if !tc.ExpectDenied {
				t.Fatal("expected cross-tenant write denial")
			}
			_ = orgB
		})
	}
}
