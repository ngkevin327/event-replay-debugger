export function SnapshotsPanel({
  snapshots,
}: {
  snapshots: {
    consumer_group: string;
    topic: string;
    partition: number;
    offset: number;
  }[];
}) {
  if (!snapshots.length) {
    return (
      <div className="empty-state">
        <p>No offset snapshots captured for this incident.</p>
      </div>
    );
  }

  return (
    <div className="data-table-wrap">
      <table className="data-table">
        <caption className="table-caption">Consumer offset snapshots</caption>
        <thead>
          <tr>
            <th>Group</th>
            <th>Topic</th>
            <th>Partition</th>
            <th>Offset</th>
          </tr>
        </thead>
        <tbody>
          {snapshots.map((s, i) => (
            <tr key={i}>
              <td className="text-mono">{s.consumer_group}</td>
              <td className="text-mono">{s.topic}</td>
              <td>{s.partition}</td>
              <td className="text-mono">{s.offset}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
