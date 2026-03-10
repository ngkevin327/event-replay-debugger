import { useRef } from "react";

let lastScrub = 0;

export function throttleScrub(fn: (v: number) => void, ms = 16) {
  return (v: number) => {
    const now = performance.now();
    if (now - lastScrub < ms) return;
    lastScrub = now;
    fn(v);
  };
}

export function TimelineScrubber({
  max,
  value,
  onChange,
}: {
  max: number;
  value: number;
  onChange: (index: number) => void;
}) {
  const throttled = useRef(throttleScrub(onChange)).current;
  return (
    <div className="timeline-scrubber-wrap">
      <label htmlFor="timeline-scrub">Scrub timeline</label>
      <input
        id="timeline-scrub"
        type="range"
        min={0}
        max={Math.max(0, max)}
        value={value}
        aria-valuemin={0}
        aria-valuemax={max}
        aria-valuenow={value}
        onChange={(e) => throttled(Number(e.target.value))}
      />
      <span className="timeline-row__meta">
        {value + 1} / {max + 1}
      </span>
    </div>
  );
}
