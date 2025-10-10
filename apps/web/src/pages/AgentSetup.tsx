import { CopyBlock } from "@/components/CopyBlock";

const HELM_CMD = `helm install replay-agent ./deploy/helm/replay-agent \\
  -f agent-values.yaml`;

export function HelmCopyBlock() {
  return <CopyBlock text={HELM_CMD} label="helm command" />;
}

export function AgentSetupWizard() {
  return (
    <div>
      <h1>Agent setup</h1>
      <p>Install the capture agent with Helm.</p>
      <HelmCopyBlock />
    </div>
  );
}

export default function AgentSetup() {
  return <AgentSetupWizard />;
}
