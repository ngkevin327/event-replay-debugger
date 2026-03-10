import { CopyBlock } from "@/components/CopyBlock";
import { PageHeader } from "@/components/PageHeader";
import { OpsIllustration } from "@/components/illustrations/OpsIllustration";

const HELM_CMD = `helm upgrade --install replay-agent ./deploy/helm/replay-agent \\
  --namespace replay --create-namespace \\
  -f agent-values.yaml`;

export function HelmCopyBlock() {
  return <CopyBlock text={HELM_CMD} label="helm command" />;
}

export function AgentSetupWizard() {
  return (
    <div>
      <PageHeader
        title="Agent setup"
        description="Install the capture agent in your Kubernetes cluster to ship Kafka events to Replay."
      />

      <div className="card setup-steps stagger-children">
        <div className="setup-step">
          <span className="setup-step__num">1</span>
          <div>
            <h2 style={{ fontSize: "var(--text-base)", color: "var(--color-text)" }}>
              Create an API key
            </h2>
            <p style={{ margin: 0, color: "var(--color-text-secondary)", fontSize: "var(--text-sm)" }}>
              Generate a project API key under Settings and store it in a Kubernetes secret.
            </p>
          </div>
        </div>
        <div className="setup-step">
          <span className="setup-step__num">2</span>
          <div>
            <h2 style={{ fontSize: "var(--text-base)", color: "var(--color-text)" }}>
              Configure values
            </h2>
            <p style={{ margin: 0, color: "var(--color-text-secondary)", fontSize: "var(--text-sm)" }}>
              Set <code style={{ fontFamily: "var(--font-mono)" }}>projectId</code>, ingest URL, and
              topic allowlist in <code style={{ fontFamily: "var(--font-mono)" }}>agent-values.yaml</code>.
            </p>
          </div>
        </div>
        <div className="setup-step">
          <span className="setup-step__num">3</span>
          <div style={{ flex: 1 }}>
            <h2 style={{ fontSize: "var(--text-base)", color: "var(--color-text)" }}>
              Install with Helm
            </h2>
            <HelmCopyBlock />
          </div>
        </div>
      </div>

      <div className="empty-state empty-state--illustrated" style={{ marginTop: "2rem" }}>
        <OpsIllustration variant="agent" />
        <p>After install, agent heartbeats appear on your Dashboard.</p>
      </div>
    </div>
  );
}

export default function AgentSetup() {
  return <AgentSetupWizard />;
}
