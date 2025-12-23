import React, { useState } from "react";
import type { ControlPanelProps, BlendApiResponse } from "../prop-types";
function ControlPanelTileButton({ highlight, children, label, onClick }) {
  return (
    <button
      onClick={onClick}
      className={`group relative aspect-square w-18.75 select-none  ${highlight ? "bg-green-400" : "bg-white"}  p-3 outline-2 outline-black transition-all flex flex-col items-center justify-center gap-1`}
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

export function ControlPanel({
  setMode,
  setUsers,
  setBlendPercent,
  blendApiResponse: BlendApiResponse,
}: ControlPanelProps) {
  // async function updateBlendFromAPI({ user, mode, timeDuration }) {
  //   console.log("Updating blend from API:", { user, mode, timeDuration });
  //   try {
  //     const baseURL = "http:/localhost:3000/v1/blend/new";
  //     const params = new URLSearchParams({
  //       category: mode,
  //       user: user,
  //       timeDuration: timeDuration,
  //     });
  //     const url = new URL(baseURL);
  //     url.search = params.toString();
  //     const response = await fetch(url, { credentials: "include" });
  //     if (!response.ok) {
  //       throw new Error(`Backend request error. Status: ${response.status}`);
  //     }
  //     const data = await response.json();
  //     const newVal = data["blend_percentage"];

  //     console.log("API response data:", data);
  //     console.log("Blend percentage from API:", newVal);
  //     setBlendPercent(newVal);
  //   } catch (err) {
  //     console.error("API error:", err);
  //   }

  //   return;
  // }

  const [selectedGroup1, setSelectedGroup1] = useState(null);
  const [selectedGroup2, setSelectedGroup2] = useState(null);
  const [group3Selected, setGroup3Selected] = useState(true);

  const handleGroup1Click = (option) => {
    setSelectedGroup1(option);
    setGroup3Selected(false);
  };

  const handleGroup2Click = (option) => {
    setSelectedGroup2(option);
    setGroup3Selected(false);
  };

  const handleGroup3Click = () => {
    setGroup3Selected(!group3Selected);
    if (!group3Selected) {
      setSelectedGroup1(null);
      setSelectedGroup2(null);
    }
  };

  async function updateBlendFromStoredValue({ mode, timeDuration }) {
    console.log("Updating blend from stored value:", {
      mode,
      timeDuration,
    });
    try {
      var typeBlend: {
        OneMonth: number;
        ThreeMonth: number;
        OneYear: number;
      } | null = null; // e.g., "artist", "track", "album"
      var displayedMode = "";
      var newVal = 0;
      switch (mode) {
        case "artist":
          typeBlend = BlendApiResponse.ArtistBlend;
          displayedMode = "Artists Only";
          break;
        case "track":
          typeBlend = BlendApiResponse.TrackBlend;
          displayedMode = "Songs Only";
          break;
        case "album":
          typeBlend = BlendApiResponse.AlbumBlend;
          displayedMode = "Albums Only";
          break;
        case "default": // OverallBlendNum, not unknown case
          newVal = BlendApiResponse.OverallBlendNum;
          displayedMode = "Default";
          break;
        default:
          console.warn("Unknown mode:", mode);
          return;
      }

      if (typeBlend !== null) {
        switch (timeDuration) {
          case "1month":
            newVal = typeBlend.OneMonth;
            displayedMode += " - Last 1 Month";
            break;
          case "3month":
            newVal = typeBlend.ThreeMonth;
            displayedMode += " - Last 3 Month";
            break;
          case "1year":
            newVal = typeBlend.OneYear;
            displayedMode += " - Last 1 Year";
            break;
          default:
            console.warn("Unknown time duration:", timeDuration);
            return;
        }
      }

      setBlendPercent(newVal);
      setMode(displayedMode);
      setUsers(BlendApiResponse.Usernames.join(" + "));
      console.log("Updated blend percentage:", newVal);
    } catch (err) {
      console.error("Error retrieving stored blend percentage:", err);
    }

    return;
  }
  // const user = "test2002";
  const [curMode, setCurMode] = useState("artist");
  const [curDuration, setCurDuration] = useState("3month");
  React.useEffect(() => {
    updateBlendFromStoredValue({
      // user: user,
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
            highlight={selectedGroup1 == "1month"}
            label="Last 1 Month"
            onClick={() => {
              setCurDuration("1month");
              handleGroup1Click("1month");
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
            highlight={selectedGroup1 == "3month"}
            label="Last 3 Month"
            onClick={() => {
              setCurDuration("3month");
              handleGroup1Click("3month");
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
            highlight={selectedGroup1 == "1year"}
            label="Last 1 Year"
            onClick={() => {
              setCurDuration("1year");
              handleGroup1Click("1year");
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
            highlight={selectedGroup2 == "artist"}
            label="Artists Only"
            onClick={() => {
              setCurMode("artist");
              handleGroup2Click("artist");
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
            highlight={selectedGroup2 == "track"}
            label="Songs Only"
            onClick={() => {
              setCurMode("track");
              handleGroup2Click("track");
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
            highlight={selectedGroup2 == "album"}
            label="Albums"
            onClick={() => {
              setCurMode("album");
              handleGroup2Click("album");
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
          <ControlPanelTileButton
            highlight={group3Selected}
            onClick={() => {
              setCurMode("default");
              handleGroup3Click();
            }}
            label="Default"
          >
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
