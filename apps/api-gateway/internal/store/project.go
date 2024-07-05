package store

import "context"

// CreateProject inserts a project under an org.
func (s *Store) CreateProject(ctx context.Context, orgID, name string) (Project, error) {
	var p Project
	err := s.pool.QueryRow(ctx,
		`INSERT INTO projects (org_id, name) VALUES ($1, $2)
		 RETURNING id, org_id, name, created_at`,
		orgID, name,
	).Scan(&p.ID, &p.OrgID, &p.Name, &p.CreatedAt)
	return p, err
}

// ListProjectsByOrg returns projects for an organization.
func (s *Store) ListProjectsByOrg(ctx context.Context, orgID string) ([]Project, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, org_id, name, created_at FROM projects WHERE org_id = $1 ORDER BY created_at DESC`,
		orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.OrgID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetProject loads a project by id.
func (s *Store) GetProject(ctx context.Context, id string) (Project, error) {
	var p Project
	err := s.pool.QueryRow(ctx,
		`SELECT id, org_id, name, created_at FROM projects WHERE id = $1`, id,
	).Scan(&p.ID, &p.OrgID, &p.Name, &p.CreatedAt)
	return p, err
}
