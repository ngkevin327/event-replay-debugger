export function CopyBlock({ text, label }: { text: string; label: string }) {
  async function copy() {
    await navigator.clipboard.writeText(text);
  }
  return (
    <div className="copy-block">
      <pre>{text}</pre>
      <button type="button" onClick={copy} aria-label={`Copy ${label}`}>
        Copy
      </button>
    </div>
  );
}
