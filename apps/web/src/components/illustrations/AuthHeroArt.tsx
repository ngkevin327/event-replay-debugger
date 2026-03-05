/** Background art for auth hero panel */
export function AuthHeroArt() {
  return (
    <div className="auth-hero-art" aria-hidden>
      <svg viewBox="0 0 400 320" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle cx="200" cy="160" r="120" stroke="rgba(56,189,248,0.15)" strokeWidth="1" />
        <circle cx="200" cy="160" r="80" stroke="rgba(99,102,241,0.2)" strokeWidth="1" />
        <path
          d="M80 200 Q200 80 320 200"
          stroke="url(#hero-line)"
          strokeWidth="2"
          fill="none"
          strokeDasharray="8 6"
        />
        <rect x="120" y="120" width="160" height="80" rx="12" fill="rgba(20,28,46,0.6)" stroke="rgba(148,163,184,0.2)" />
        <rect x="136" y="140" width="80" height="8" rx="4" fill="rgba(56,189,248,0.4)" />
        <rect x="136" y="156" width="128" height="6" rx="3" fill="rgba(148,163,184,0.25)" />
        <rect x="136" y="170" width="96" height="6" rx="3" fill="rgba(148,163,184,0.15)" />
        <circle cx="320" cy="100" r="24" fill="rgba(34,211,238,0.12)" />
        <circle cx="100" cy="220" r="16" fill="rgba(99,102,241,0.15)" />
        <defs>
          <linearGradient id="hero-line" x1="80" y1="200" x2="320" y2="200">
            <stop stopColor="#22d3ee" />
            <stop offset="1" stopColor="#a855f7" />
          </linearGradient>
        </defs>
      </svg>
    </div>
  );
}
