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
    <input
      type="range"
      min={0}
      max={Math.max(0, max)}
      value={value}
      aria-label="Timeline scrubber"
      onChange={(e) => throttled(Number(e.target.value))}
    />
  );
}
