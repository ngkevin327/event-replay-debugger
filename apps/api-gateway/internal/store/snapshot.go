package store

import "context"

// InsertOffsetSnapshot records a consumer offset snapshot for an incident.
func (s *Store) InsertOffsetSnapshot(ctx context.Context, incidentID, consumerGroup, topic string, partition int, offset int64) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO offset_snapshots (incident_id, consumer_group, topic, partition, offset)
		 VALUES ($1, $2, $3, $4, $5)`,
		incidentID, consumerGroup, topic, partition, offset,
	)
	return err
}
