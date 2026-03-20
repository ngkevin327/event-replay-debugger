import type { ReactNode } from "react";
import { BrandLogo } from "./BrandLogo";
import { AuthHeroArt } from "./illustrations/AuthHeroArt";

export function AuthLayout({
  title,
  subtitle,
  children,
  footer,
}: {
  title: string;
  subtitle: string;
  children: ReactNode;
  footer: ReactNode;
}) {
  return (
    <div className="auth-page">
      <aside className="auth-page__hero" aria-hidden>
        <div className="auth-page__hero-inner">
          <div className="auth-page__brand">
            <BrandLogo size="lg" />
          </div>
          <h2 className="auth-page__tagline">
            See how production incidents actually unfolded
          </h2>
          <div className="auth-page__desc-highlight">
            <p className="auth-page__desc">
              Replay captures Kafka behavior, reconstructs timelines, and runs
              deterministic replays so your team debugs async failures with
              confidence—not guesswork.
            </p>
          </div>
          <ul className="auth-page__features">
            <li>Virtualized incident timelines</li>
            <li>Workflow graph visualization</li>
            <li>Deterministic replay & divergence reports</li>
          </ul>
          <AuthHeroArt />
        </div>
      </aside>
      <main className="auth-page__main">
        <div className="auth-card">
          <header className="auth-card__header">
            <div className="auth-card__brand">
              <BrandLogo size="lg" />
            </div>
            <h1>{title}</h1>
            <p className="auth-card__subtitle">{subtitle}</p>
          </header>
          {children}
          <footer className="auth-card__footer">{footer}</footer>
        </div>
      </main>
    </div>
  );
}
