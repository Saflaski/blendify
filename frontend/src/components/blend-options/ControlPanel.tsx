import React, { useState } from "react";
import type { ControlPanelProps, CardApiResponse } from "../prop-types";
import ArtistIcon from "@/assets/images/artist.svg";
import AlbumIcon from "@/assets/images/artist.svg";
import TrackIcon from "@/assets/images/track.svg";
import BlendifyWhiteIcon from "@/assets/images/blendifyIconWhite.svg";
import BlendifyIcon from "@/assets/images/blendifyIcon.svg";

function ControlPanelTileButton({ highlight, children, label, onClick }) {
  return (
    <button
      onClick={onClick}
      className={`group relative aspect-square w-14 select-none  ${highlight ? "bg-[#D84727] text-slate-100 outline-[#000000]" : "bg-white text-slate-950 outline-black "}  p-3 outline-2 transition-all flex flex-col items-center justify-center gap-1`}
    >
      <div className="flex items-center justify-center flex-1 w-full">
        <div
          className={`w-full h-full flex items-center ${highlight ? "brightness-0 invert" : ""} justify-center`}
        >
          {children}
        </div>
      </div>

      {label ? (
        <span className="text-[8px] font-semibold tracking-wide  leading-none">
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

  const [selectedGroup1, setSelectedGroup1] = useState("");
  const [selectedGroup2, setSelectedGroup2] = useState("");
  const [group3Selected, setGroup3Selected] = useState(true);

  const handleGroup1Click = (option) => {
    setSelectedGroup1(option);
    setGroup3Selected(false);
  };

  const handleGroup2Click = (option) => {
    setSelectedGroup2(option);
    // if (selectedGroup1 == "") {
    //   setSelectedGroup1("3month");
    // }
    setGroup3Selected(false);
  };

  const handleGroup3Click = () => {
    setGroup3Selected(!group3Selected);
    if (!group3Selected) {
      setSelectedGroup1("");
      setSelectedGroup2("");
    }
  };

  const handleGroup3Click_Alternate = () => {
    setGroup3Selected(true);
    setSelectedGroup1("");
    setSelectedGroup2("");
  };

  async function updateBlendFromStoredValue({ mode, timeDuration }) {
    try {
      var typeBlend: {
        OneMonth: number;
        ThreeMonth: number;
        OneYear: number;
      } | null = null; // e.g., "artist", "track", "album"

      const conditionOnlyModeSelected = mode != "" && timeDuration == "";
      const conditionOnlyDurationSelected = mode == "" && timeDuration != "";
      console.log("Condition Only Mode: ", conditionOnlyModeSelected);
      console.log("Condition Only Duration: ", conditionOnlyDurationSelected);

      var displayedMode = "";
      var newVal = 0;
      console.log(mode === "default" && conditionOnlyModeSelected);
      console.log(mode === "default" && conditionOnlyDurationSelected);
      console.log(
        mode === "default" &&
          !conditionOnlyDurationSelected &&
          !conditionOnlyModeSelected,
      );
      if (
        (mode === "default" && conditionOnlyModeSelected) ||
        (mode === "default" && conditionOnlyDurationSelected) ||
        (mode === "default" &&
          !conditionOnlyDurationSelected &&
          !conditionOnlyModeSelected)
      ) {
        newVal = BlendApiResponse.OverallBlendNum;
        displayedMode = "Default mode";
        handleGroup3Click_Alternate();
      } else {
        console.log("ELSE: ", mode);
        switch (mode) {
          case "artist":
            typeBlend = BlendApiResponse.ArtistBlend;
            displayedMode = "Artists";
            handleGroup2Click("artist");
            break;
          case "track":
            typeBlend = BlendApiResponse.TrackBlend;
            displayedMode = "Songs";
            handleGroup2Click("track");
            break;
          case "album":
            typeBlend = BlendApiResponse.AlbumBlend;
            displayedMode = "Albums";
            handleGroup2Click("album");
            break;
          default:
            typeBlend = BlendApiResponse.ArtistBlend;
            displayedMode = "Artists";
            console.log("Defaulting to Artists Only");
            handleGroup2Click("artist");
            break;
        }

        switch (timeDuration) {
          case "1month":
            newVal = typeBlend.OneMonth;
            displayedMode += " in last 1 month";
            handleGroup1Click("1month");
            break;
          case "3month":
            newVal = typeBlend.ThreeMonth;
            displayedMode += " in last 3 months";
            handleGroup1Click("3month");
            break;
          case "1year":
            newVal = typeBlend.OneYear;
            displayedMode += " in last 1 year";
            handleGroup1Click("1year");
            break;
          default:
            newVal = typeBlend.OneMonth;
            displayedMode += " - Last 1 Month";
            console.log("Defaulting to Last 1 Month");
            handleGroup1Click("1month");
            break;
        }
      }

      setBlendPercent(newVal);
      setMode(displayedMode);
      setUsers(BlendApiResponse.Usernames);
      console.log("Updated blend percentage:", newVal);
    } catch (err) {
      console.error("Error retrieving stored blend percentage:", err);
    }

    return;
  }
  const [curMode, setCurMode] = useState("default");
  const [curDuration, setCurDuration] = useState("");

  React.useEffect(() => {
    console.log("UPDATE: \nmode: ", curMode, " \nDuration: ", curDuration);
    updateBlendFromStoredValue({
      // user: user,
      mode: curMode,
      timeDuration: curDuration,
    });
  }, [curMode, curDuration]);

  return (
    <div className=" flex items-center justify-center  bg-inherit outline-black ">
      <div className="grid grid-row-3 items-center gap-3">
        {/* DATE RANGES */}
        <div className="outline-2 outline-black p-2 flex gap-4">
          <ControlPanelTileButton
            highlight={selectedGroup1 == "1month"}
            label=""
            onClick={() => {
              if (curMode == "default") setCurMode("artist");
              setCurDuration("1month");
              // handleGroup1Click("1month");
            }}
          >
            <p className="font-[Roboto_Mono] text-xs font-bold">1 MONTH</p>
          </ControlPanelTileButton>
          <ControlPanelTileButton
            highlight={selectedGroup1 == "3month"}
            label=""
            onClick={() => {
              if (curMode == "default") setCurMode("artist");
              setCurDuration("3month");
              // handleGroup1Click("3month");
            }}
          >
            <p className="font-[Roboto_Mono] text-xs font-bold">3 MONTH</p>
          </ControlPanelTileButton>
          <ControlPanelTileButton
            highlight={selectedGroup1 == "1year"}
            label=""
            onClick={() => {
              if (curMode == "default") setCurMode("artist");
              setCurDuration("1year");
              // handleGroup1Click("1year");
            }}
          >
            <p className="font-[Roboto_Mono] text-xs font-bold">1 YEAR</p>
          </ControlPanelTileButton>
        </div>

        {/* --- ARTIST / GENRE / SONG  --- */}
        <div className="outline-2 outline-black p-2 flex gap-4">
          <ControlPanelTileButton
            highlight={selectedGroup2 == "artist"}
            label="Artists"
            onClick={() => {
              setCurMode("artist");
              // handleGroup2Click("artist");
            }}
          >
            <img src={ArtistIcon} alt="Artist" />
          </ControlPanelTileButton>
          <ControlPanelTileButton
            highlight={selectedGroup2 == "track"}
            label="Songs"
            onClick={() => {
              setCurMode("track");
              // handleGroup2Click("track");
            }}
          >
            <img src={TrackIcon} alt="Song" />
          </ControlPanelTileButton>
          <ControlPanelTileButton
            highlight={selectedGroup2 == "album"}
            label="Albums"
            onClick={() => {
              setCurMode("album");
              // handleGroup2Click("album");
            }}
          >
            <img src={AlbumIcon} alt="Song" />
          </ControlPanelTileButton>
        </div>

        {/* --- DEFAULT --- */}
        <div className="outline-2 outline-black w-fit mx-auto p-2">
          <ControlPanelTileButton
            highlight={group3Selected}
            onClick={() => {
              setCurMode("default");
              // handleGroup3Click();
            }}
            label="Default"
          >
            <button className="group">
              <img src={BlendifyIcon} />
            </button>
          </ControlPanelTileButton>
        </div>
      </div>
    </div>
  );
}
