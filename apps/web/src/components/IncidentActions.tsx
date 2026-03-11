import { useState } from "react";
import { useCreateShare, useExportIncident } from "@/api/hooks";
import { Modal } from "@/components/Modal";

export function ExportButton({ incidentId }: { incidentId: string }) {
  const exportMutation = useExportIncident(incidentId);
  return (
    <button
      type="button"
      className="btn btn--secondary"
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
      <button type="button" className="btn btn--secondary" onClick={() => setOpen(true)}>
        Share read-only link
      </button>
      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title="Share incident"
        footer={
          <button type="button" className="btn btn--ghost" onClick={() => setOpen(false)}>
            Close
          </button>
        }
      >
        <p>Generate an expiring read-only link for partners.</p>
        <button
          type="button"
          className="btn btn--primary"
          onClick={async () => {
            const res = await createShare.mutateAsync(72);
            setLink(window.location.origin + res.url);
          }}
          disabled={createShare.isPending}
        >
          {createShare.isPending ? "Creating…" : "Create link"}
        </button>
        {link && (
          <div className="share-link-row">
            <div className="form-field">
              <label htmlFor="share-url">Share URL</label>
              <input id="share-url" readOnly value={link} />
            </div>
            <button
              type="button"
              className="btn btn--secondary"
              onClick={() => navigator.clipboard.writeText(link)}
            >
              Copy link
            </button>
          </div>
        )}
      </Modal>
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
