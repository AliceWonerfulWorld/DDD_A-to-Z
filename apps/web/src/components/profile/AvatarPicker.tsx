import React from "react";

interface AvatarPickerProps {
  avatarUrl: string;
  githubAvatarUrl?: string;
  disabled: boolean;
  onChange: (url: string) => void;
}

export function AvatarPicker({
  avatarUrl,
  githubAvatarUrl,
  disabled,
  onChange,
}: AvatarPickerProps) {
  const templates = [
    { id: "go", url: "/avatars/GO.png", label: "Go" },
    { id: "java", url: "/avatars/JV.png", label: "Java" },
    { id: "kotlin", url: "/avatars/KL.png", label: "Kotlin" },
    { id: "lisp", url: "/avatars/LP.png", label: "Lisp" },
    { id: "linux", url: "/avatars/LX.png", label: "Linux" },
    { id: "php", url: "/avatars/PH.png", label: "PHP" },
    { id: "python", url: "/avatars/PY.png", label: "Python" },
    { id: "rust", url: "/avatars/RS.png", label: "Rust" },
    { id: "women1", url: "/avatars/W1.png", label: "女性1" },
    { id: "women2", url: "/avatars/W2.png", label: "女性2" },
    { id: "women3", url: "/avatars/W3.png", label: "女性3" },
    { id: "women4", url: "/avatars/W4.png", label: "女性4" },
    { id: "men1", url: "/avatars/M1.png", label: "男性1" },
    { id: "men2", url: "/avatars/M2.png", label: "男性2" },
    { id: "men3", url: "/avatars/M3.png", label: "男性3" },
    { id: "men4", url: "/avatars/M4.png", label: "男性4" },
    { id: "men5", url: "/avatars/M5.png", label: "男性5" },
  ];

  if (githubAvatarUrl) {
    templates.unshift({ id: "github", url: githubAvatarUrl, label: "GitHub" });
  }

  return (
    <div style={{ width: "100%", marginTop: "-1rem" }}>
      <label
        style={{
          display: "block",
          marginBottom: "0.5rem",
          fontSize: "0.8rem",
          color: "var(--color-gold)",
          letterSpacing: "0.1em",
        }}
      >
        ▶ SELECT YOUR ICON
      </label>

      <div
        style={{
          display: "flex",
          gap: "10px",
          flexWrap: "wrap",
          justifyContent: "center",
          marginBottom: "0.8rem",
        }}
      >
        {templates.map((t) => (
          <button
            key={t.id}
            type="button"
            onClick={() => onChange(t.url)}
            disabled={disabled}
            style={{
              width: "64px",
              height: "64px",
              border:
                avatarUrl === t.url
                  ? "3px solid var(--color-gold)"
                  : "2px solid rgba(255,255,255,0.4)",
              background: "rgba(0,0,0,0.5)",
              cursor: disabled ? "default" : "pointer",
              padding: "4px",
              imageRendering: t.id === "github" ? "auto" : "pixelated",
              opacity: disabled ? 0.5 : 1,
            }}
            title={t.label}
          >
            <img
              src={t.url}
              alt={t.label}
              style={{ width: "100%", height: "100%", objectFit: "cover" }}
            />
          </button>
        ))}
      </div>
    </div>
  );
}
