import { useState } from "react";
import { TopicAllowlistEditor } from "@/components/TopicAllowlistEditor";
import { useApiKeys } from "@/api/hooks";

function ApiKeyList({ projectId }: { projectId: string }) {
  const { data } = useApiKeys(projectId);
  return (
    <ul>
      {(data?.keys ?? []).map((k) => (
        <li key={k.id}>{k.prefix}••••</li>
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
      <p>Rotated key shown once: replay_live_••••</p>
      <button type="button" onClick={onClose}>
        Done
      </button>
    </dialog>
  );
}

export function Settings() {
  const projectId = localStorage.getItem("replay_project_id") ?? "";
  const [rotateOpen, setRotateOpen] = useState(false);
  return (
    <div>
      <h1>Project settings</h1>
      <section>
        <h2>API keys</h2>
        <ApiKeyList projectId={projectId} />
        <button type="button" onClick={() => setRotateOpen(true)}>
          Rotate key
        </button>
        <RotateKeyModal open={rotateOpen} onClose={() => setRotateOpen(false)} />
      </section>
      <section>
        <h2>Topic allowlist</h2>
        <TopicAllowlistEditor projectId={projectId} />
      </section>
    </div>
  );
}

export default Settings;
