/** Decorative SVG for empty states and marketing panels */
export function OpsIllustration({ variant = "timeline" }: { variant?: "timeline" | "agent" }) {
  if (variant === "agent") {
    return (
      <svg
        className="ops-illustration"
        viewBox="0 0 200 160"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        aria-hidden
      >
        <rect x="24" y="40" width="152" height="88" rx="12" stroke="url(#ops-g)" strokeWidth="2" fill="var(--color-surface)" />
        <circle cx="48" cy="64" r="8" fill="var(--color-accent-muted)" />
        <rect x="64" y="58" width="96" height="8" rx="4" fill="var(--color-border-strong)" />
        <rect x="64" y="74" width="72" height="6" rx="3" fill="var(--color-surface-muted)" />
        <path d="M40 120 L80 96 L120 108 L160 80" stroke="url(#ops-g)" strokeWidth="2" strokeLinecap="round" />
        <defs>
          <linearGradient id="ops-g" x1="0" y1="0" x2="200" y2="160">
            <stop stopColor="var(--color-accent)" />
            <stop offset="1" stopColor="#6366f1" />
          </linearGradient>
        </defs>
      </svg>
    );
  }

  return (
    <svg
      className="ops-illustration"
      viewBox="0 0 200 160"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden
    >
      <line x1="40" y1="130" x2="40" y2="30" stroke="var(--color-border-strong)" strokeWidth="2" />
      {[0, 1, 2, 3].map((i) => (
        <g key={i}>
          <circle cx="40" cy={40 + i * 28} r="6" fill="url(#ops-g2)" />
          <rect
            x="56"
            y={32 + i * 28}
            width={120 - i * 12}
            height="20"
            rx="6"
            fill="var(--color-surface)"
            stroke="var(--color-border-strong)"
          />
        </g>
      ))}
      <defs>
        <linearGradient id="ops-g2" x1="34" y1="30" x2="46" y2="50">
          <stop stopColor="#22d3ee" />
          <stop offset="1" stopColor="#6366f1" />
        </linearGradient>
      </defs>
    </svg>
  );
}
