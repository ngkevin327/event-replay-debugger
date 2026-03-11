import { useState } from "react";
import { TopicAllowlistEditor } from "@/components/TopicAllowlistEditor";
import { PageHeader } from "@/components/PageHeader";
import { Modal } from "@/components/Modal";
import { useApiKeys } from "@/api/hooks";

function ApiKeyList({ projectId }: { projectId: string }) {
  const { data } = useApiKeys(projectId);
  const keys = data?.keys ?? [];
  if (!keys.length) {
    return <p className="text-muted">No API keys yet.</p>;
  }
  return (
    <ul className="list-plain">
      {keys.map((k) => (
        <li key={k.id} className="text-mono">
          {k.prefix}••••
        </li>
      ))}
    </ul>
  );
}

export function RotateKeyModal({
  open,
  onClose,
}: {
  open: boolean;
  onClose: () => void;
}) {
  return (
    <Modal
      open={open}
      onClose={onClose}
      title="API key rotated"
      footer={
        <button type="button" className="btn btn--primary" onClick={onClose}>
          Done
        </button>
      }
    >
      <p>Copy your new key now — it won&apos;t be shown again.</p>
      <p className="text-mono mt-4">replay_live_••••••••••••</p>
    </Modal>
  );
}

export function Settings() {
  const projectId = localStorage.getItem("replay_project_id") ?? "";
  const [rotateOpen, setRotateOpen] = useState(false);
  return (
    <div>
      <PageHeader
        title="Settings"
        description="Manage API keys, topic allowlists, and project configuration."
      />

      <section className="section card">
        <h2 className="section__title">API keys</h2>
        <ApiKeyList projectId={projectId} />
        <button type="button" className="btn btn--secondary mt-4" onClick={() => setRotateOpen(true)}>
          Rotate key
        </button>
        <RotateKeyModal open={rotateOpen} onClose={() => setRotateOpen(false)} />
      </section>

      <section className="section card">
        <h2 className="section__title">Topic allowlist</h2>
        <TopicAllowlistEditor projectId={projectId} />
      </section>
    </div>
  );
}

export default Settings;
