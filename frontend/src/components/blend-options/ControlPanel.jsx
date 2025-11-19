import React, { useState } from "react";

function ControlPanelTileButton({ children, label, onClick }) {
  return (
    <button
      onClick={onClick}
      className="group relative aspect-square w-18.75 select-none bg-white p-3 outline-2 outline-black transition-all flex flex-col items-center justify-center gap-1"
    >
      <div className="flex items-center justify-center flex-1 w-full">
        <div className="w-3 h-3 flex items-center justify-center">
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

export function ControlPanel({ setBlendPercent }) {
  async function updateBlendFromAPI({ user, mode, timeDuration }) {
    console.log("Updating blend from API:", { user, mode, timeDuration });
    try {
      const baseURL = "http:/localhost:3000/v1/blends/new";
      const params = new URLSearchParams({
        category: mode,
        user: user,
        timeDuration: timeDuration,
      });
      const url = new URL(baseURL);
      url.search = params.toString();
      const response = await fetch(url, { credentials: "include" });
      if (!response.ok) {
        throw new Error(`Backend request error. Status: ${response.status}`);
      }
      const data = await response.json();
      const newVal = data["blend_percentage"];

      console.log("API response data:", data);
      console.log("Blend percentage from API:", newVal);
      setBlendPercent(newVal);
    } catch (err) {
      console.error("API error:", err);
    }

    return;
  }
  const user = "test2002";
  const [curMode, setCurMode] = useState("artist");
  const [curDuration, setCurDuration] = useState("1month");
  React.useEffect(() => {
    updateBlendFromAPI({
      user: user,
      mode: curMode,
      timeDuration: curDuration,
    });
  }, [curMode, curDuration]);

  return (
    <div className=" flex items-center justify-center md:pt-5 bg-inherit outline-black p-5">
      <div className="grid grid-row-3 items-center gap-8">
        {/* DATE RANGES */}
        <div className="outline-2 outline-black p-2 flex gap-4">
          <ControlPanelTileButton
            label="Last 1 Month"
            onClick={() => {
              setCurDuration("1month");
            }}
          >
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <rect x="4" y="4" width="16" height="16" />
            </svg>
          </ControlPanelTileButton>
          <ControlPanelTileButton
            label="Last 3 Month"
            onClick={() => {
              setCurDuration("3month");
            }}
          >
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <polygon points="12,2 22,22 2,22" />
            </svg>
          </ControlPanelTileButton>
          <ControlPanelTileButton
            label="Last 1 Year"
            onClick={() => {
              setCurDuration("1year");
            }}
          >
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
        <div className="outline-2 outline-black p-2 flex gap-4">
          <ControlPanelTileButton
            label="Artists Only"
            onClick={() => {
              setCurMode("artist");
            }}
          >
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <path d="M4 4h16v16H4z M4 4l16 16" />
            </svg>
          </ControlPanelTileButton>
          <ControlPanelTileButton
            label="Songs Only"
            onClick={() => {
              setCurMode("track");
            }}
          >
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              stroke="black"
              strokeWidth="2"
            >
              <path d="M12 2L2 12l10 10 10-10z" />
            </svg>
          </ControlPanelTileButton>
          <ControlPanelTileButton
            label="Albums"
            onClick={() => {
              setCurMode("album");
            }}
          >
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

        {/* --- DEFAULT --- */}
        <div className="outline-2 outline-black w-fit mx-auto p-2">
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
      </div>
    </div>
  );
}
