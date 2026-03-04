export function BrandLogo({ size = "md" }: { size?: "sm" | "md" }) {
  return (
    <div className={`brand-logo brand-logo--${size}`} aria-hidden>
      <span className="brand-logo__mark">
        <svg viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="16" cy="16" r="14" stroke="url(#replay-grad)" strokeWidth="2" />
          <path
            d="M11 16h4l2-5 2 10 2-5h4"
            stroke="url(#replay-grad)"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <defs>
            <linearGradient id="replay-grad" x1="4" y1="4" x2="28" y2="28">
              <stop stopColor="#22d3ee" />
              <stop offset="1" stopColor="#6366f1" />
            </linearGradient>
          </defs>
        </svg>
      </span>
      <span className="brand-logo__text">Replay</span>
    </div>
  );
}
