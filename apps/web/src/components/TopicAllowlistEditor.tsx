import { useState } from "react";
import { useUpdateAllowlist } from "@/api/hooks";

export function TopicAllowlistEditor({ projectId }: { projectId: string }) {
  const [topics, setTopics] = useState<string[]>([]);
  const [input, setInput] = useState("");
  const update = useUpdateAllowlist();

  function addTopic() {
    if (!input.trim()) return;
    setTopics((t) => [...t, input.trim()]);
    setInput("");
  }

  function save() {
    update.mutate(topics);
  }

  return (
    <div>
      <div className="chips">
        {topics.map((t) => (
          <span key={t} className="chip">
            {t}
          </span>
        ))}
      </div>
      <input
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && addTopic()}
        aria-label="Add topic"
      />
      <button type="button" onClick={save}>
        Save allowlist
      </button>
      <input type="hidden" data-project={projectId} />
    </div>
  );
}
