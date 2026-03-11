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
    <div className="allowlist-editor">
      <div className="chips">
        {topics.map((t) => (
          <span key={t} className="chip">
            {t}
          </span>
        ))}
      </div>
      <div className="form-field">
        <label htmlFor="topic-input">Add topic</label>
        <input
          id="topic-input"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && (e.preventDefault(), addTopic())}
          placeholder="payments.settlement"
        />
      </div>
      <div className="toolbar-inline">
        <button type="button" className="btn btn--ghost" onClick={addTopic}>
          Add topic
        </button>
        <button type="button" className="btn btn--primary" onClick={save} disabled={update.isPending}>
          {update.isPending ? "Saving…" : "Save allowlist"}
        </button>
      </div>
      <input type="hidden" data-project={projectId} />
    </div>
  );
}
