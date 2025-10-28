import { useEffect } from "react";

export function useTimelineKeyboard(
  activeIndex: number,
  count: number,
  onChange: (index: number) => void,
) {
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === "j" || e.key === "ArrowDown") {
        onChange(Math.min(count - 1, activeIndex + 1));
      }
      if (e.key === "k" || e.key === "ArrowUp") {
        onChange(Math.max(0, activeIndex - 1));
      }
      if (e.key === "h") onChange(Math.max(0, activeIndex - 10));
      if (e.key === "l") onChange(Math.min(count - 1, activeIndex + 10));
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [activeIndex, count, onChange]);
}
