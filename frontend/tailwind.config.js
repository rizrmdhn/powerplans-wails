/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ["class"],
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        panel: {
          bg: "#14161A",
          surface: "#1B1E24",
          border: "#282C34",
          hover: "#20242C",
        },
        amber: {
          DEFAULT: "#E8A33D",
          dim: "#7A5A2A",
        },
        danger: "#C1443C",
        ink: {
          DEFAULT: "#EDEAE2",
          muted: "#8B8F98",
          faint: "#565A63",
        },
      },
      fontFamily: {
        mono: ["JetBrains Mono", "ui-monospace", "SFMono-Regular", "monospace"],
        sans: ["Inter", "ui-sans-serif", "system-ui", "sans-serif"],
      },
      boxShadow: {
        glow: "0 0 12px 2px rgba(232, 163, 61, 0.35)",
      },
      borderRadius: {
        sm: "4px",
        md: "6px",
      },
    },
  },
  plugins: [],
};
