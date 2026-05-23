
export function NameInput({
  disabled,
  onChange,
  username,
}: {
  disabled: boolean;
  onChange: (username: string) => void;
  username: string;
}) {
  return (
    <div style={{ width: "100%" }}>
      <label
        style={{
          display: "block",
          marginBottom: "0.5rem",
          fontSize: "0.8rem",
          color: "var(--color-gold)",
          letterSpacing: "0.1em",
        }}
      >
        ▶ ENTER YOUR NAME
      </label>
      <div style={{ position: "relative" }}>
        <input
          type="text"
          value={username}
          onChange={(e) => onChange(e.target.value)}
          disabled={disabled}
          style={{
            width: "100%",
            padding: "0.8rem",
            fontSize: "1.2rem",
            fontFamily: "var(--font-dot)",
            background: "rgba(0,0,0,0.5)",
            color: "var(--color-pixel-white)",
            border: "2px solid rgba(255,255,255,0.4)",
            outline: "none",
            textAlign: "center",
          }}
          onFocus={(e) => (e.target.style.borderColor = "var(--color-gold)")}
          onBlur={(e) => (e.target.style.borderColor = "rgba(255,255,255,0.4)")}
        />
      </div>
    </div>
  );
}
