import * as React from "react";
import { cn } from "../../lib/utils";

type Variant = "primary" | "ghost" | "danger" | "outline";
type Size = "sm" | "md";

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  size?: Size;
}

const variants: Record<Variant, string> = {
  primary:
    "bg-amber text-panel-bg font-medium hover:brightness-110 disabled:bg-amber-dim disabled:text-ink-faint",
  outline:
    "border border-panel-border text-ink hover:bg-panel-hover disabled:opacity-40",
  ghost: "text-ink-muted hover:text-ink hover:bg-panel-hover",
  danger: "text-danger hover:bg-danger/10 disabled:opacity-30",
};

const sizes: Record<Size, string> = {
  sm: "h-7 px-2.5 text-xs",
  md: "h-9 px-4 text-sm",
};

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "outline", size = "md", ...props }, ref) => (
    <button
      ref={ref}
      className={cn(
        "inline-flex items-center justify-center gap-1.5 rounded-sm transition-colors duration-150 disabled:cursor-not-allowed",
        "focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-amber/60",
        variants[variant],
        sizes[size],
        className
      )}
      {...props}
    />
  )
);
Button.displayName = "Button";
