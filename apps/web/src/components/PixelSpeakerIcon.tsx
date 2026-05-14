interface PixelSpeakerIconProps {
  muted: boolean;
}

export function PixelSpeakerIcon({ muted }: PixelSpeakerIconProps) {
  const pixel = {
    position: "absolute" as const,
    background: "currentColor",
    boxShadow: "1px 1px 0 #000",
  };

  return (
    <span
      aria-hidden="true"
      style={{
        position: "relative",
        display: "inline-block",
        width: "18px",
        height: "14px",
        color: "inherit",
      }}
    >
      <span style={{ ...pixel, left: 0, top: 5, width: 3, height: 4 }} />
      <span style={{ ...pixel, left: 3, top: 3, width: 4, height: 8 }} />
      <span style={{ ...pixel, left: 7, top: 1, width: 3, height: 12 }} />
      {!muted && (
        <>
          <span style={{ ...pixel, left: 12, top: 3, width: 2, height: 2 }} />
          <span style={{ ...pixel, left: 12, top: 9, width: 2, height: 2 }} />
          <span style={{ ...pixel, left: 15, top: 1, width: 2, height: 2 }} />
          <span style={{ ...pixel, left: 15, top: 11, width: 2, height: 2 }} />
          <span style={{ ...pixel, left: 16, top: 5, width: 2, height: 4 }} />
        </>
      )}
      {muted && (
        <>
          <span style={{ ...pixel, left: 13, top: 3, width: 2, height: 2, background: "#ffb0aa" }} />
          <span style={{ ...pixel, left: 15, top: 5, width: 2, height: 2, background: "#ffb0aa" }} />
          <span style={{ ...pixel, left: 13, top: 7, width: 2, height: 2, background: "#ffb0aa" }} />
          <span style={{ ...pixel, left: 11, top: 9, width: 2, height: 2, background: "#ffb0aa" }} />
        </>
      )}
    </span>
  );
}
