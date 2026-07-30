import * as React from "react";
import { cn } from "../../lib/utils";

type Tone = "amber" | "muted" | "outline";

export function Badge({
  className,
  tone = "muted",
  children,
}: {
  className?: string;
  tone?: Tone;
  children: React.ReactNode;
}) {
  const tones: Record<Tone, string> = {
    amber: "bg-amber/15 text-amber border border-amber/30",
    muted: "bg-panel-hover text-ink-muted border border-panel-border",
    outline: "text-ink-faint border border-panel-border",
  };
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-sm px-1.5 py-0.5 text-[10px] font-mono uppercase tracking-wider",
        tones[tone],
        className
      )}
    >
      {children}
    </span>
  );
}
