import { useState } from "react";
import { TopicAllowlistEditor } from "@/components/TopicAllowlistEditor";
import { PageHeader } from "@/components/PageHeader";
import { useApiKeys } from "@/api/hooks";

function ApiKeyList({ projectId }: { projectId: string }) {
  const { data } = useApiKeys(projectId);
  const keys = data?.keys ?? [];
  if (!keys.length) {
    return <p style={{ color: "var(--color-text-secondary)" }}>No API keys yet.</p>;
  }
  return (
    <ul style={{ listStyle: "none", padding: 0, margin: 0 }}>
      {keys.map((k) => (
        <li
          key={k.id}
          style={{
            fontFamily: "var(--font-mono)",
            fontSize: "var(--text-sm)",
            padding: "0.5rem 0",
          }}
        >
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
  if (!open) return null;
  return (
    <dialog open>
      <div className="modal-card">
        <h2>API key rotated</h2>
        <p>Copy your new key now — it won&apos;t be shown again.</p>
        <p style={{ fontFamily: "var(--font-mono)", fontSize: "var(--text-sm)" }}>
          replay_live_••••
        </p>
        <button type="button" className="btn btn--primary" onClick={onClose}>
          Done
        </button>
      </div>
    </dialog>
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
        <button
          type="button"
          className="btn btn--secondary"
          style={{ marginTop: "1rem" }}
          onClick={() => setRotateOpen(true)}
        >
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
