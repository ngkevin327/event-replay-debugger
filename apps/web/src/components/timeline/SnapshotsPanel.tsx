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
  return (
    <table>
      <caption>Consumer offset snapshots</caption>
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
            <td>{s.consumer_group}</td>
            <td>{s.topic}</td>
            <td>{s.partition}</td>
            <td>{s.offset}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
