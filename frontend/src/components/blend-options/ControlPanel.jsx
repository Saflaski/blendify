import React from "react";

function ControlPanelTileButton({ children, label, onClick }) {
  return (
    <button
      onClick={onClick}
      className="group relative aspect-square w-20 select-none bg-white p-3 outline-2 outline-black transition-all flex flex-col items-center justify-center gap-1"
    >
      <div className="flex items-center justify-center flex-1 w-full">
        <div className="w-6 h-6 flex items-center justify-center">
          {children}
        </div>
      </div>

      {label ? (
        <span className="text-[9px] font-semibold tracking-wide text-neutral-800 leading-none">
          {label}
        </span>
      ) : null}
    </button>
  );
}

export function ControlPanel() {
  return (
    <div className="flex items-center justify-center bg-inherit outline-2 outline-black p-5">
      <div className="flex flex-row items-center gap-8">
        {/* --- DEFAULT --- */}
        <div>
          <ControlPanelTileButton label="Default">
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <circle cx="12" cy="12" r="9" />
            </svg>
          </ControlPanelTileButton>
        </div>

        <div className="flex gap-4">
          <ControlPanelTileButton label="Last 1 Month">
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <rect x="4" y="4" width="16" height="16" />
            </svg>
          </ControlPanelTileButton>
          <ControlPanelTileButton label="Last 3 Month">
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <polygon points="12,2 22,22 2,22" />
            </svg>
          </ControlPanelTileButton>
          <ControlPanelTileButton label="Last 1 Year">
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <path d="M12 2v20M2 12h20" />
            </svg>
          </ControlPanelTileButton>
        </div>

        {/* --- ARTIST / GENRE / SONG  --- */}
        <div className="flex gap-4">
          <ControlPanelTileButton label="Artists Only">
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <path d="M4 4h16v16H4z M4 4l16 16" />
            </svg>
          </ControlPanelTileButton>
          <ControlPanelTileButton label="Songs Only">
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <path d="M12 2L2 12l10 10 10-10z" />
            </svg>
          </ControlPanelTileButton>
          <ControlPanelTileButton label="Artists and Songs">
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <path d="M3 3h18v18H3z M3 9h18M9 3v18" />
            </svg>
          </ControlPanelTileButton>
        </div>
      </div>
    </div>
  );
}
