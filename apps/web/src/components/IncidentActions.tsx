import { useState } from "react";
import { useCreateShare, useExportIncident } from "@/api/hooks";

export function ExportButton({ incidentId }: { incidentId: string }) {
  const exportMutation = useExportIncident(incidentId);
  return (
    <button
      type="button"
      onClick={() => exportMutation.mutate()}
      disabled={exportMutation.isPending}
    >
      {exportMutation.isPending ? "Exporting…" : "Export JSON"}
    </button>
  );
}

export function ShareLinkModal({ incidentId }: { incidentId: string }) {
  const [open, setOpen] = useState(false);
  const [link, setLink] = useState("");
  const createShare = useCreateShare(incidentId);

  return (
    <>
      <button type="button" onClick={() => setOpen(true)}>
        Share read-only link
      </button>
      {open && (
        <dialog open>
          <h2>Share incident</h2>
          <p>Generate an expiring read-only link for partners.</p>
          <button
            type="button"
            onClick={async () => {
              const res = await createShare.mutateAsync(72);
              setLink(window.location.origin + res.url);
            }}
          >
            Create link
          </button>
          {link && (
            <p>
              <input readOnly value={link} aria-label="Share link" />
              <button type="button" onClick={() => navigator.clipboard.writeText(link)}>
                Copy link
              </button>
            </p>
          )}
          <button type="button" onClick={() => setOpen(false)}>
            Close
          </button>
        </dialog>
      )}
    </>
  );
}

export function IncidentActions({ incidentId }: { incidentId: string }) {
  return (
    <div className="incident-actions">
      <ExportButton incidentId={incidentId} />
      <ShareLinkModal incidentId={incidentId} />
    </div>
  );
}
