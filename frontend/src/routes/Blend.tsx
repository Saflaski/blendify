// import { DropDownMenu } from "../components/blend-options/dropdownmenu";
import { ControlPanel } from "../components/blend-options/ControlPanel";
import { useLocation, useNavigate } from "react-router-dom";
import React, { useRef, useState, useEffect } from "react";
import { toBlob } from "html-to-image";
import {
  ControlPanelProps,
  BlendApiResponse,
  BlendApiResponseSchema,
} from "../components/prop-types";
import { set, z } from "zod";
import SplitRatioBar from "../components/SplitRatioBar";

// type ControlPanelProps = {
//   setBlendPercent: (num: number) => void;
//   blendApiResponse: BlendApiResponse;
// };

// type BlendApiResponse = {
//   usernames: string[];
//   overallBlendNum: number;
//   ArtistBlend: TypeBlend;
//   AlbumBlend: TypeBlend;
//   TrackBlend: TypeBlend;
// };
// type MetricKey = keyof BlendApiResponse;

// type TypeBlend = {
//   OneMonth: number;
//   ThreeMonth: number;
//   OneYear: number;
// };

export function Blend() {
  // ------ If user is from invite link and not Add button -------
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  const [blendId, setBlendId] = useState<string | null>(null);
  const [navLinkId, setNavLinkId] = useState<string | null>(null);
  const [userBlendData, setUserBlendData] = useState<BlendApiResponse>(
    {} as BlendApiResponse,
  );
  type LocationState = {
    id?: string;
    value?: string;
  };

  console.log(location.state);

  useEffect(() => {
    const state = location.state as LocationState | null;
    if (state?.id == "blendid") {
      state?.value != null
        ? setBlendId(state?.value)
        : console.log("undefined blendid given");
    } else if (state?.id == "linkid") {
      state?.value != null
        ? setNavLinkId(state?.value)
        : console.log("undefined navlinkid given");
    }
  }, [location.state]);

  console.log("NavLinkId state: ", navLinkId);
  const getBlendIdFromInviteLink = async () => {
    //From URL Paste
    const params = new URLSearchParams(location.search);
    const urlInvite = params.get("invite");

    //From Add button
    const value = location.state;
    // const navigateInvite = value?.invite;

    const navigateInvite = navLinkId;

    //Log them
    console.log("urlInvite: ", urlInvite);
    console.log("Navigated Invite Link Data: ", navigateInvite);

    const invite = navigateInvite ?? urlInvite;
    console.log("Getting blendid from Link: ", invite);
    //Get blendid as authenticated user.
    const requestBlendId = async () => {
      try {
        const res = await fetch("http://localhost:3000/v1/blend/add", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({ value: invite }),
        });

        if (res.status == 401) {
          navigate(
            `/login?redirectTo=${encodeURIComponent(location.pathname + location.search)}`,
          );
          return;
        }

        if (!res.ok) {
          const data = await res.json().catch(() => ({}));
          setError(data.message || "Invite is invalid.");
          // setLoading(false);
          return;
        }

        const data = await res.json();
        const blendIdFromAPI = data["blendId"];
        setBlendId(blendIdFromAPI);

        // setLoading(false);
      } catch (err) {
        console.error(err);
        setError("Something went wrong. Please try again.");
        setLoading(false);
      }
    };
    requestBlendId();

    // If user clicked on existing blend from homepage
  };
  if (blendId === null) {
    console.log("Getting blendid from API");
    getBlendIdFromInviteLink();
  }

  //

  console.log("Final blendId to use: ", blendId);
  useEffect(() => {
    console.log("Getting data for blendId (1): ", blendId);
    const getBlendData = async () => {
      console.log("Getting data for blendId (2): ", blendId);

      try {
        const encodedValue = encodeURIComponent(blendId as string);
        const res = await fetch(
          `http://localhost:3000/v1/blend/data?blendId=${encodedValue}`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
            credentials: "include",
          },
        );

        if (res.status == 401) {
          navigate(
            `/login?redirectTo=${encodeURIComponent(location.pathname + location.search)}`,
          );
          return;
        }

        if (!res.ok) {
          const data = await res.json().catch(() => ({}));
          setError(data.message || "Blend ID is invalid.");
          setLoading(false);
          return;
        }

        const data = await res.json();
        console.log("Blend data received:", data);
        // const userData = JSON.parse(data) as BlendApiResponse;
        const userData = BlendApiResponseSchema.parse(data);
        console.log("Parsed blend data:", userData);
        setUserBlendData(userData);
        setLoading(false);
      } catch (err) {
        console.error(err);
        setError("Something went wrong. Please try again.");
        setLoading(false);
      }
    };

    if (blendId == null) {
      setError("Could not get blendid.");
      console.log("Blend ID is null, cannot get data?");
    } else {
      getBlendData();
      console.log("Getting blend data");
    }
  }, [blendId]);

  // if (error != null) {
  //   console.error(error);
  // }

  // ----- Copy button functionality -----
  const captureRef = useRef(null); //Div to be captured
  const [isCapturing, setIsCapturing] = useState(false); //To hide the button during screenshot

  const [copied, setCopied] = useState(false); //For tooltip
  const hideTimer = useRef<number | null>(null); //Tooltip hide timer

  useEffect(() => {
    return () => {
      if (hideTimer.current !== null) {
        clearTimeout(hideTimer.current);
      }
    }; // cleanup on unmount
  }, []);

  const handleScreenshot = async () => {
    setIsCapturing(true);
    await new Promise((r) => setTimeout(r, 50));
    if (!captureRef.current) return;

    try {
      const blob = await toBlob(captureRef.current, {
        pixelRatio: 2, // like scale
        cacheBust: true,
        backgroundColor: "#F8F3E9",
        skipFonts: true, // ← avoids parsing/embedding fonts
      });

      setIsCapturing(false);

      if (!blob) throw new Error("Failed to create screenshot");

      await navigator.clipboard.write([
        new ClipboardItem({
          "image/png": blob,
        }),
      ]);

      setCopied(true);
      if (hideTimer.current !== null) {
        clearTimeout(hideTimer.current);
      }
      // clearTimeout(hideTimer.current);
      hideTimer.current = setTimeout(() => setCopied(false), 1400);
    } catch (err) {
      console.error("Clipboard fail", err);
    }

    // const link = document.createElement("a");
    // link.download = "screenshot.png";
    // link.href = dataUrl;
    // link.click();
  };
  const [blendPercent, setBlendPercent] = useState(3);
  const [mode, setMode] = useState("Default");
  const [users, setUsers] = useState<string>("You and someone");

  const props: ControlPanelProps = {
    setMode,
    setUsers,
    setBlendPercent,
    blendApiResponse: userBlendData,
  };

  useEffect(() => {
    if (userBlendData != undefined && loading == false) {
      setBlendPercent(userBlendData.OverallBlendNum);
      setMode("Default mode");
      if (userBlendData.Usernames.length == 2)
        setUsers(
          `${userBlendData.Usernames[0]} and ${userBlendData.Usernames[1]}`,
        );
    }
  }, [userBlendData]);
  // setBlendPercent(userBlendData.OverallBlendNum);

  return (
    <div className="w-full ">
      <div className="w-full md:w-[60%] flex pt-4 flex-col md:flex-row gap-x-5 mx-auto text-center px-4 gap-y-4 md:px-0 py-0 md:py-5">
        {/* <div className="flex justify-left"></div> */}

        {/* <div className="md:flex md:flex-wrap pr-2 mt-8 lg:grid lg:grid-cols-2 "> Old*/}
        {/* LEFT CONTENT AREA */}
        <div className="  md:w-[40%] flex flex-col flex-wrap items-center justify-baseline gap-y-5">
          {/* Back button */}
          {/* <div className=" text-[14px]">
            <button
              type="button"
              onClick={() => navigate("/home")}
              className="inline-flex items-center outline-2 h-auto  font-[Roboto_Mono] font-bold border border-black/10 bg-white px-4 py-2 text-black shadow-sm hover:shadow "
            >
              &lt; Blends
            </button>
          </div> */}
          {/* Player card */}
          <div className="w-full flex justify-center  ">
            <div className="w-full flex justify-center ">
              <div
                ref={captureRef}
                className="shine-element relative ring-2 ring-black bg-neutral-200 
    w-58 md:w-58 lg:w-66 p-4 aspect-[2/3]
    bg-size-[auto_120px] bg-[url(/src/assets/images/topography.svg)]"
              >
                {!isCapturing && (
                  <button
                    onClick={handleScreenshot}
                    className="absolute top-1 right-1 outline outline-black bg-inherit p-1"
                  >
                    <img
                      src="/src/assets/images/copy.svg"
                      className="w-4 h-4"
                    />
                  </button>
                )}

                {copied && (
                  <div
                    className="absolute right-8 top-2 bg-gray-500 text-white 
        text-[10px] px-2 py-0.5 shadow animate-fade-in-out"
                  >
                    Copied!
                  </div>
                )}

                <h1
                  className="mt-0 text-6xl md:text-5xl lg:text-7xl 
    leading-none font-[Filepile] font-normal -left-2 tracking-tight text-black relative inline-block"
                >
                  {loading ? "--" : blendPercent}
                  <span className="absolute bottom-1 -right-12 text-lg font-normal ">
                    /100
                  </span>
                </h1>

                <p
                  className="mt-1 font-bold text-gray-800"
                  style={{ fontSize: "clamp(1rem, 1.5vw, 1.2rem)" }}
                >
                  {users}
                </p>

                <p className="mt-1 text-sm/snug font-bold text-gray-800">
                  {mode}
                </p>

                <div className="grid grid-row-2 gap-1 text-left text-black font-[Roboto_Mono] mt-2">
                  <ul>
                    <p className="font-semibold text-base">Top Artists</p>

                    <li className="text-[13px] font-medium leading-tight">
                      - Clairo
                    </li>
                    <li className="text-[13px] font-medium leading-tight">
                      - Men I Trust
                    </li>
                    <li className="text-[13px] font-medium leading-tight">
                      - Bring Me The Horizon
                    </li>
                  </ul>

                  <ul>
                    <p className="font-semibold text-base">Top Songs</p>
                    <li className="text-[13px] font-medium leading-tight">
                      - Bababooey 2
                    </li>
                    <li className="text-[13px] font-medium leading-tight">
                      - Come Down
                    </li>
                    <li className="text-[13px] font-medium leading-tight">
                      - Bags
                    </li>
                  </ul>
                </div>

                <div className="flex flex-col justify-between gap-3 absolute bottom-8 left-1/2 -translate-x-1/2 size-12 h-auto">
                  <img src="/src/assets/images/lastfm.svg" />
                </div>

                <p className="text-center w-full font-[Filepile] text-[#000000] bg-clip-text bg-amber-500 tracking-widest font-medium absolute bottom-2 left-1/2 -translate-x-1/2 text-xs  text-shadow-2xs">
                  blendify.fm
                </p>
              </div>
            </div>
          </div>
          {/* End of player card */}

          {/* Control panel */}
          <div className=" flex flex-wrap justify-center items-center ">
            <ControlPanel {...props} />
          </div>
          {/* End of control panel */}
        </div>

        {/* RIGHT CONTENT AREA */}
        <div className=" md:w-[60%] outline-amber-600 flex flex-col flex-wrap items-center justify-baseline gap-y-5">
          {/* Top blend artists section */}
          <section className=" w-full flex flex-col">
            <h2 className="text-xl md:text-2xl font-semibold text-black mb-4 text-center">
              Top blend artists
            </h2>
            {/* Placeholder list/cards — replace with real data */}
            <div className="flex flex-col gap-y-4 items-center gap-4 text-zinc-950">
              <SplitRatioBar
                itemName="Clairo"
                valueA={40}
                valueB={30}
                urlToNavigateA="https://www.last.fm/user/saflas"
                urlToNavigateB="https://www.last.fm/user/test2002"
              />
              <SplitRatioBar
                itemName="Bring Me The Horizon"
                valueA={59}
                valueB={20}
                urlToNavigateA="https://www.last.fm/user/saflas"
                urlToNavigateB="https://www.last.fm/user/test2002"
              />
              <SplitRatioBar
                itemName="Linkin Park"
                valueA={10}
                valueB={40}
                urlToNavigateA="https://www.last.fm/user/saflas"
                urlToNavigateB="https://www.last.fm/user/test2002"
              />
            </div>
          </section>

          {/* Top blend songs section */}
          <section className=" w-full">
            <h2 className="text-xl md:text-2xl font-semibold text-black mb-4 text-center">
              Top blend songs
            </h2>
            {/* Placeholder list/cards — replace with real data */}
            <div className="flex flex-col gap-y-4 items-center gap-4  text-zinc-950">
              <SplitRatioBar
                itemName="Charm"
                valueA={40}
                valueB={30}
                urlToNavigateA="https://www.last.fm/user/saflas"
                urlToNavigateB="https://www.last.fm/user/test2002"
              />
              <SplitRatioBar
                itemName="Sempiternal"
                valueA={59}
                valueB={20}
                urlToNavigateA="https://www.last.fm/user/saflas"
                urlToNavigateB="https://www.last.fm/user/test2002"
              />
              <SplitRatioBar
                itemName="Hybrid Theory"
                valueA={10}
                valueB={40}
                urlToNavigateA="https://www.last.fm/user/saflas"
                urlToNavigateB="https://www.last.fm/user/test2002"
              />
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}

const fetchBlendPercentage = async (label) => {
  await new Promise((r) => setTimeout(r, 500));
};
