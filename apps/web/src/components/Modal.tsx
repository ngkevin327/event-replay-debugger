import type { ReactNode } from "react";

export function Modal({
  open,
  onClose,
  title,
  children,
  footer,
}: {
  open: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  footer?: ReactNode;
}) {
  if (!open) return null;

  return (
    <dialog
      open
      className="modal-dialog"
      onClose={onClose}
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div
        className="modal-card"
        role="document"
        onClick={(e) => e.stopPropagation()}
        onKeyDown={() => {}}
      >
        {title && <h2>{title}</h2>}
        {children}
        {footer && <div className="modal-actions">{footer}</div>}
      </div>
    </dialog>
  );
}
