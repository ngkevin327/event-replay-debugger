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
            <h2 className="setup-step__title">Create an API key</h2>
            <p className="setup-step__desc">
              Generate a project API key under Settings and store it in a Kubernetes secret.
            </p>
          </div>
        </div>
        <div className="setup-step">
          <span className="setup-step__num">2</span>
          <div>
            <h2 className="setup-step__title">Configure values</h2>
            <p className="setup-step__desc">
              Set <code className="code-inline">projectId</code>, ingest URL, and topic allowlist in{" "}
              <code className="code-inline">agent-values.yaml</code>.
            </p>
          </div>
        </div>
        <div className="setup-step">
          <span className="setup-step__num">3</span>
          <div className="setup-step__body">
            <h2 className="setup-step__title">Install with Helm</h2>
            <HelmCopyBlock />
          </div>
        </div>
      </div>

      <div className="empty-state empty-state--illustrated mt-6">
        <OpsIllustration variant="agent" />
        <p>After install, agent heartbeats appear on your Dashboard.</p>
      </div>
    </div>
  );
}

export default function AgentSetup() {
  return <AgentSetupWizard />;
}
