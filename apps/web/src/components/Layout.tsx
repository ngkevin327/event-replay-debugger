import { NavLink, Outlet } from "react-router-dom";
import { AppHeader } from "./AppHeader";

export function Layout() {
  return (
    <div className="layout">
      <a href="#main" className="skip-link">
        Skip to content
      </a>
      <AppHeader />
      <nav className="nav" aria-label="Main">
        <NavLink to="/" end>
          Dashboard
        </NavLink>
        <NavLink to="/incidents">Incidents</NavLink>
        <NavLink to="/settings">Settings</NavLink>
        <NavLink to="/agents/setup">Agent setup</NavLink>
      </nav>
      <main id="main">
        <Outlet />
      </main>
    </div>
  );
}
