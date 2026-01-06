// import { DropDownMenu } from "../components/blend-options/dropdownmenu";
import { ControlPanel } from "../components/blend-options/ControlPanel";
import { useLocation, useNavigate } from "react-router-dom";
import React, { useRef, useState, useEffect } from "react";
import CardBackground from "@/assets/images/topography.svg";
import CopyIcon from "@/assets/images/copy.svg";
import LastfmIcon from "@/assets/images/lastfm.svg";
import BackArrow from "@/assets/images/arrow_back.svg";
import FrontArrow from "@/assets/images/arrow_front.svg";
import "@/assets/styles/index.css";
import { toBlob } from "html-to-image";

import { chromium, firefox, webkit, BrowserType } from "playwright";
import {
  ControlPanelProps,
  CardApiResponse,
  CardApiResponseSchema,
  CatalogueBlendResponse,
  CatalogueBlendSchema,
} from "../components/prop-types";
import { set, z } from "zod";
import {
  SplitRatioBar,
  SplitRatioBarSkeleton,
} from "../components/SplitRatioBar";
import { API_BASE_URL } from "../constants";

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

function useLocalStorageState<T>(key: string, initialValue: T) {
  const [state, setState] = useState<T>(() => {
    try {
      const stored = localStorage.getItem(key);
      if (!stored) return initialValue;

      const parsed: unknown = JSON.parse(stored);
      return parsed as T;
    } catch {
      return initialValue;
    }
  });

  useEffect(() => {
    localStorage.setItem(key, JSON.stringify(state));
  }, [key, state]);

  return [state, setState] as const;
}

const ARTIST_3_MONTH_KEY = "ARTIST_3_MONTH_KEY";
const TRACK_3_MONTH_KEY = "TRACK_3_MONTH_KEY";
const ARTIST_12_MONTH_KEY = "ARTIST_12_MONTH_KEY";
const TRACK_12_MONTH_KEY = "TRACK_12_MONTH_KEY";
const ARTIST_1_MONTH_KEY = "ARTIST_1_MONTH_KEY";
const TRACK_1_MONTH_KEY = "TRACK_1_MONTH_KEY";
const BLEND_ID_KEY = "blend_id";

export function Blend() {
  // ------ If user is from invite link and not Add button -------
  const [error, setError] = useState<string | null>(null);
  const [cardLoading, setCardLoading] = useState(true);
  const [catalogueLoading, setCatalogueLoading] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  const locationState = location.state as LocationState | null;

  const [blendId, setBlendId] = useState<string | null>(() =>
    getInitialBlendId(locationState),
  );

  const [navLinkId, setNavLinkId] = useState<string | null>(null);
  const [userCardData, setUserCardData] = useState<CardApiResponse>(
    {} as CardApiResponse,
  );

  const [userCatalogueArtist3MonthData, setUserCatalogueArtist3MonthData] =
    useLocalStorageState<CatalogueBlendResponse[]>(ARTIST_3_MONTH_KEY, []);
  const [userCatalogueArtist1MonthData, setUserCatalogueArtist1MonthData] =
    useLocalStorageState<CatalogueBlendResponse[]>(ARTIST_1_MONTH_KEY, []);
  const [userCatalogueArtist1YearData, setUserCatalogueArtist1YearData] =
    useLocalStorageState<CatalogueBlendResponse[]>(ARTIST_12_MONTH_KEY, []);
  const [userCatalogueTrack1YearData, setUserCatalogueTrack1YearData] =
    useLocalStorageState<CatalogueBlendResponse[]>(TRACK_12_MONTH_KEY, []);
  const [userCatalogueTrack3MonthData, setUserCatalogueTrack3MonthData] =
    useLocalStorageState<CatalogueBlendResponse[]>(TRACK_3_MONTH_KEY, []);
  const [userCatalogueTrack1MonthData, setUserCatalogueTrack1MonthData] =
    useLocalStorageState<CatalogueBlendResponse[]>(TRACK_1_MONTH_KEY, []);

  const [catArt1Year, setCatArt1Year] = useState(true);
  const [catArt3Month, setCatArt3Month] = useState(true);
  const [catArt1Month, setCatArt1Month] = useState(true);
  const [catTrack1Year, setCatTrack1Year] = useState(true);
  const [catTrack3Month, setCatTrack3Month] = useState(true);
  const [catTrack1Month, setCatTrack1Month] = useState(true);

  const currentTime = new Date().getTime();
  type LocationState = {
    id?: string;
    value?: string;
  };

  function getInitialBlendId(
    locationState: LocationState | null,
  ): string | null {
    if (locationState?.id === "blendid" && locationState.value) {
      return locationState.value;
    }

    return localStorage.getItem(BLEND_ID_KEY);
  }

  console.log(location.state);
  useEffect(() => {
    const state = location.state as LocationState | null;

    if (state?.id === "blendid" && state.value) {
      setBlendId(state.value);
      navigate(location.pathname, { replace: true });
      return;
    }

    if (state?.id === "linkid" && state.value) {
      setNavLinkId(state.value);
      setBlendId(null);
      navigate(location.pathname, { replace: true });
    }
  }, [location.state, navigate, location.pathname]);

  // useEffect(() => {
  //   const state = location.state as LocationState | null;

  //   // 1. Consume navigation state
  //   if (state?.id === "blendid" && state.value) {
  //     const newBlendId = state.value;

  //     console.log("Setting blendId from location state:", newBlendId);
  //     setBlendId(newBlendId);

  //     // 🔥 IMPORTANT: clear location.state so it doesn't re-run
  //     navigate(location.pathname, { replace: true });

  //     return;
  //   }

  //   if (state?.id === "linkid" && state.value) {
  //     const newLinkId = state.value;

  //     console.log("Setting navLinkId from location state:", newLinkId);
  //     setNavLinkId(newLinkId);
  //     setBlendId(null);

  //     // 🔥 clear navigation state here too
  //     navigate(location.pathname, { replace: true });

  //     return;
  //   }

  //   // 2. Fallback to localStorage
  //   const storedBlendId = localStorage.getItem(BLEND_ID_KEY);
  //   if (storedBlendId) {
  //     console.log("Setting blendId from localStorage:", storedBlendId);
  //     setBlendId(storedBlendId);
  //     return;
  //   }

  //   // 3. Final fallback
  //   console.log("No blendId found, setting null");
  //   setBlendId(null);
  // }, [location.state, navigate, location.pathname]);

  console.log("NavLinkId state: ", navLinkId);
  useEffect(() => {
    console.log("BlendId after checking 3 places: ", blendId);
  }, [blendId]);
  useEffect(() => {
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
          const res = await fetch(API_BASE_URL + "/blend/add", {
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
          localStorage.setItem(BLEND_ID_KEY, blendIdFromAPI);

          // setLoading(false);
        } catch (err) {
          console.error(err);
          setError("Something went wrong. Please try again.");
          // setCardLoading(false);
        }
      };
      requestBlendId();

      // If user clicked on existing blend from homepage
    };
    if (blendId == null) {
      console.log("Getting blendid from API");
      getBlendIdFromInviteLink();
    }
  }, []);

  useEffect(() => {
    if (!blendId) return;

    const loadAllCatalogueData = async () => {
      try {
        setCardLoading(true);
        setCatalogueLoading(true);

        await Promise.all([
          getCatalogueBlendData(
            "3month",
            "artist",
            blendId,
            userCatalogueArtist3MonthData,
            setUserCatalogueArtist3MonthData,
            setCatArt3Month,
            setError,
          ),
          getCatalogueBlendData(
            "3month",
            "track",
            blendId,
            userCatalogueTrack3MonthData,
            setUserCatalogueTrack3MonthData,
            setCatTrack3Month,
            setError,
          ),
          getCatalogueBlendData(
            "12month",
            "artist",
            blendId,
            userCatalogueArtist1YearData,
            setUserCatalogueArtist1YearData,
            setCatArt1Year,
            setError,
          ),
          getCatalogueBlendData(
            "12month",
            "track",
            blendId,
            userCatalogueTrack1YearData,
            setUserCatalogueTrack1YearData,
            setCatTrack1Year,
            setError,
          ),
          getCatalogueBlendData(
            "1month",
            "track",
            blendId,
            userCatalogueTrack1MonthData,
            setUserCatalogueTrack1MonthData,
            setCatArt1Month,
            setError,
          ),
          getCatalogueBlendData(
            "1month",
            "artist",
            blendId,
            userCatalogueArtist1MonthData,
            setUserCatalogueArtist1MonthData,
            setCatArt1Month,
            setError,
          ),
        ]);

        await getCardBlendData(); // runs AFTER all catalogue calls
        setCatArt1Year(false);
        setCatArt3Month(false);
        setCatTrack1Year(false);
        setCatTrack3Month(false);
        setCatTrack1Month(false);
        setCatArt1Month(false);
        setCatalogueLoading(false);
      } finally {
        setCardLoading(false);
      }
    };

    loadAllCatalogueData();
  }, [blendId]);

  // console.log("Getting data for blendId (1): ", blendId);
  const getCardBlendData = async () => {
    console.log("Getting data for blendId (2): ", blendId);

    try {
      const encodedValue = encodeURIComponent(blendId as string);
      const res = await fetch(
        `${API_BASE_URL}/blend/carddata?blendId=${encodedValue}`,
        {
          method: "GET",
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
        setCardLoading(false);
        return;
      }

      const data = await res.json();
      console.log("Blend data received:", data);
      // const userData = JSON.parse(data) as BlendApiResponse;
      const userData = CardApiResponseSchema.parse(data);
      console.log("Parsed blend data:", userData);
      setUserCardData(userData);
      setCardLoading(false);
    } catch (err) {
      console.error(err);
      setError("Something went wrong. Please try again.");
      setCardLoading(false);
    }
  };

  // if (blendId == null) {
  //   setError("Could not get blendid.");
  //   console.log("Blend ID is null, cannot get data?");
  // } else {
  //   if (catalogueLoading == 4) {
  //     getCardBlendData();
  //   }
  //   console.log("Getting card blend data");
  // }

  async function downloadCatalogueData(duration: string, category: string) {
    const params = {
      blendId: blendId as string,
      duration: duration,
      category: category,
    };

    const queryString = new URLSearchParams(params).toString();
    const res = await fetch(
      `${API_BASE_URL}/blend/cataloguedata?${queryString}`,
      {
        method: "GET",
        credentials: "include",
      },
    );

    if (res.status == 401) {
      navigate(
        `/login?redirectTo=${encodeURIComponent(location.pathname + location.search)}`,
      );
      return null;
    }

    if (!res.ok) {
      const data = await res.json().catch(() => ({}));
      setError(data.message || "Blend ID is invalid.");
      // setCatalogueLoading(false);
      return null;
    }
    return res;
  }

  const getCatalogueBlendData = async (
    duration: string,
    category: string,
    blendId: string,
    data: any[],
    setData: (data: any[]) => void,
    setCatalogueLoading: (loading: boolean) => void,
    // setLoading: (loading: boolean) => void,
    setError: (msg: string) => void,
  ) => {
    console.log("Getting data for blendId:", blendId);

    if (data.length > 0) {
      setCatalogueLoading(false);
      return data;
    } else {
      try {
        const res = await downloadCatalogueData(duration, category);

        if (!res) {
          throw new Error("Catalogue data fetch returned null");
        }

        const data = await res.json();
        console.log("Catalogue blend data received:", data);

        const parsedData = data.map((item: any) =>
          CatalogueBlendSchema.parse(item),
        );

        setData(parsedData);
        return parsedData;
      } catch (err) {
        console.error(err);
        setError("Something went wrong. Please try again.");
      } finally {
        // setCatalogueLoading(catalogueLoading + 1);
        setCatalogueLoading(false);
        // console.log("+1");
      }
    }
  };

  // useEffect(() => {
  //   console.log("Loading user catalogue artist blend data:");

  //   if (!blendId) {
  //     setError("Could not get blendid.");
  //     console.log("Blend ID is null, cannot get data?");
  //     return;
  //   }

  //   cardLoading
  //     ? getCatalogueBlendData(
  //         "3month",
  //         "artist",
  //         blendId,
  //         setUserCatalogueArtist3MonthData,
  //         setCatArt3Month,
  //         setError,
  //       )
  //     : null;
  // }, [blendId]);

  // useEffect(() => {
  //   console.log("Loading user catalogue artist blend data:");

  //   if (!blendId) {
  //     setError("Could not get blendid.");
  //     console.log("Blend ID is null, cannot get data?");
  //     return;
  //   }

  //   cardLoading
  //     ? getCatalogueBlendData(
  //         "3month",
  //         "track",
  //         blendId,
  //         setUserCatalogueTrack3MonthData,
  //         setCatTrack3Month,
  //         setError,
  //       )
  //     : null;
  // }, [blendId]);

  // useEffect(() => {
  //   console.log("Loading user catalogue artist blend data:");

  //   if (!blendId) {
  //     setError("Could not get blendid.");
  //     console.log("Blend ID is null, cannot get data?");
  //     return;
  //   }

  //   cardLoading
  //     ? getCatalogueBlendData(
  //         "12month",
  //         "artist",
  //         blendId,
  //         setUserCatalogueArtist1YearData,
  //         setCatArt1Year,
  //         setError,
  //       )
  //     : null;
  // }, [blendId]);

  // useEffect(() => {
  //   console.log("Loading user catalogue artist blend data:");

  //   if (!blendId) {
  //     setError("Could not get blendid.");
  //     console.log("Blend ID is null, cannot get data?");
  //     return;
  //   }

  //   cardLoading
  //     ? getCatalogueBlendData(
  //         "12month",
  //         "track",
  //         blendId,
  //         setUserCatalogueTrack1YearData,
  //         setCatTrack1Year,
  //         setError,
  //       )
  //     : null;
  // }, [blendId]);

  // useEffect(() => {
  //   if catalogueLoading
  // }, [catalogueLoading])

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
  const [users, setUsers] = useState<string[]>(["", ""]);

  const props: ControlPanelProps = {
    setMode,
    setUsers,
    setBlendPercent,
    blendApiResponse: userCardData,
  };

  useEffect(() => {
    if (userCardData != undefined && cardLoading == false) {
      setBlendPercent(userCardData.OverallBlendNum);
      setMode("Default mode");
      if (userCardData.Usernames.length == 2) setUsers(userCardData.Usernames);
    }
  }, [userCardData]);
  // setBlendPercent(userBlendData.OverallBlendNum);
  const [showHint, setShowHint] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => setShowHint(false), 4000);
    return () => clearTimeout(timer);
  }, []);

  type OpenSection = "3months" | "12months" | null;

  const [openSection, setOpenSection] = useState<OpenSection>("3months");

  const toggleSection = (section: OpenSection) => {
    setOpenSection((prev) => (prev === section ? null : section));
  };

  type ArtistRange = "1months" | "3months" | "12months";
  type TrackRange = "1months" | "3months" | "12months";
  const ranges: ArtistRange[] = ["1months", "3months", "12months"];
  const trackRanges: TrackRange[] = ["1months", "3months", "12months"];
  const [currentArtistRangeIndex, setCurrentArtistRangeIndex] = useState(0);
  const [currentTrackRangeIndex, setCurrentTrackRangeIndex] = useState(0);

  const artistRangeLabel = {
    "1months": "ARTIST - LAST 1 MONTH",
    "3months": "ARTIST - LAST 3 MONTHS",
    "12months": "ARTIST - LAST 12 MONTHS",
  };

  const trackRangeLabel = {
    "1months": "TRACK - LAST 1 MONTH",
    "3months": "TRACK - LAST 3 MONTHS",
    "12months": "TRACK - LAST 12 MONTHS",
  };

  const currentArtistRange = ranges[currentArtistRangeIndex];
  const currentTrackRange = ranges[currentTrackRangeIndex];
  const goPrev = (
    setCurrentRangeIndex: (value: React.SetStateAction<number>) => void,
  ) => {
    setCurrentRangeIndex((prev) => (prev === 0 ? ranges.length - 1 : prev - 1));
  };

  const goNext = (
    setCurrentRangeIndex: (value: React.SetStateAction<number>) => void,
  ) => {
    setCurrentRangeIndex((prev) => (prev === ranges.length - 1 ? 0 : prev + 1));
  };

  return (
    <div className="w-full ">
      <div className="w-full md:w-[60%] flex pt-4 flex-col md:flex-row gap-x-5 mx-auto text-center px-4 gap-y-4 md:px-0 py-0 md:py-5">
        {/* <div className="flex justify-left"></div> */}

        {/* <div className="md:flex md:flex-wrap pr-2 mt-8 lg:grid lg:grid-cols-2 "> Old*/}
        {/* LEFT CONTENT AREA */}
        <div className="  md:w-[40%] flex flex-col flex-wrap items-center justify-baseline gap-y-5">
          <div
            className={`text-black font-[Roboto_Mono] italic    ${!catalogueLoading && !cardLoading ? "hidden" : "lg:hidden block"} `}
          >
            <p className="text-lg font-semibold">Loading data</p>
            <p
              className={`${showHint ? "hidden" : "block"} text-sm transition ease-in`}
            >
              First blend? This might take a bit while we fetch all your music
              data (and stay nice to the Last.fm API).
            </p>
          </div>
          {/* Player card */}
          <div className="w-full flex justify-center  ">
            <div className="w-full flex justify-center ">
              <div
                ref={captureRef}
                className={`shine-element relative ring-2 ring-black bg-neutral-200 
    w-58 md:w-58 lg:w-66 p-4 aspect-[2/3]
    bg-size-[auto_120px] `}
                style={{
                  backgroundImage: `url(${CardBackground})`,
                }}
              >
                {!isCapturing && (
                  <button
                    onClick={handleScreenshot}
                    className="absolute top-1 right-1 outline outline-black bg-inherit p-1"
                  >
                    <img src={CopyIcon} className="w-6 h-6" />
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
    leading-none   -left-2 tracking-tight text-[#000] font-extrabold relative inline-block"
                  style={{ fontFamily: "'Roboto Mono', sans-serif" }}
                >
                  {cardLoading ? "--" : blendPercent}
                  <span className="absolute bottom-1 -right-10 text-lg font-normal ">
                    /100
                  </span>
                </h1>
                {/* <svg
                  viewBox="0 0 260 120"
                  width={220}
                  height={100}
                  className="relative top-0 left-9"
                >
                  <text
                    x="0"
                    y="96"
                    style={{ fontFamily: "Filepile" }}
                    fontSize="96"
                    dominantBaseline="alphabetic"
                  >
                    {cardLoading ? "--" : blendPercent}
                  </text>

                  <text
                    x={
                      blendPercent != undefined
                        ? blendPercent.toString().length * 56 + 35
                        : 10
                    }
                    y="101"
                    style={{ fontFamily: "Filepile" }}
                    fontSize="23"
                    dominantBaseline="alphabetic"
                  >
                    /100
                  </text>
                </svg> */}

                {/* <span className="absolute top-21.5 right-10 -translate-0 font-[Roboto_Mono] text-lg text-black font-bold ">
                    /100
                  </span> */}

                <div className="flex items-center font-[Roboto_Mono] gap-2 mt-2 justify-center text-gray-800">
                  <span
                    className="font-bold"
                    style={{ fontSize: "clamp(1rem, 1.5vw, 1.2rem)" }}
                  >
                    {users ? users[0] : "You"}
                  </span>

                  <span className="font-normal text-gray-500">and</span>

                  <span
                    className="font-bold"
                    style={{ fontSize: "clamp(1rem, 1.5vw, 1.2rem)" }}
                  >
                    {users ? users[1] : "someone"}
                  </span>
                </div>

                <p className="mt-1 text-sm/snug font-bold text-gray-800">
                  {mode}
                </p>

                <div className="grid grid-row-2  gap-1  text-left text-black font-[Roboto_Mono] mt-2">
                  <ul>
                    <p className="font-semibold text-base">Top Artists</p>

                    {userCatalogueArtist3MonthData
                      .slice(0, 3)
                      .map((item, index) => {
                        return (
                          <li
                            key={index}
                            className="text-[13px] font-medium leading-tight"
                          >
                            - {item.Name}
                          </li>
                        );
                      })}
                  </ul>

                  <ul>
                    <p className="font-semibold text-base">Top Songs</p>
                    {userCatalogueTrack3MonthData
                      .slice(0, 3)
                      .map((item, index) => {
                        return (
                          <li
                            key={index}
                            className="text-[13px] font-medium leading-tight"
                          >
                            - {item.Name}
                          </li>
                        );
                      })}
                  </ul>
                </div>

                <div className="flex flex-col justify-between gap-3 absolute bottom-8 left-1/2 -translate-x-1/2 size-12 h-auto">
                  <img src={LastfmIcon} />
                </div>

                <p className="text-center w-full font-[Quantico] font-bold text-[#404040] tracking-widest absolute bottom-2 left-1/2 -translate-x-1/2 text-[12px]  text-shadow-2xs">
                  BLENDIFY
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
          <div
            className={`text-black font-[Roboto_Mono] italic   ${!catalogueLoading && !cardLoading ? "hidden" : "hidden lg:block"} `}
          >
            <p className="text-lg font-semibold">Loading data</p>
            <p
              className={`${showHint ? "hidden" : "block"} text-sm transition ease-in`}
            >
              First blend? This might take a bit while we fetch all your music
              data (and stay nice to the Last.fm API).
            </p>
          </div>

          {/* New experimental dropdown bit */}
          <section className="w-full flex flex-col">
            <div className="flex items-center justify-center gap-4 mb-4">
              <button
                onClick={() => goPrev(setCurrentArtistRangeIndex)}
                className="text-xl font-bold text-black hover:opacity-70"
                aria-label="Previous range"
              >
                <img
                  src={BackArrow}
                  className="ring-2 pr-1.5 hover:bg-gray-200 bg-white px-1"
                  alt="Previous"
                ></img>
              </button>

              <h2 className="text-lg font-bold text-black text-center min-w-[220px]">
                {artistRangeLabel[currentArtistRange]}
              </h2>

              <button
                onClick={() => goNext(setCurrentArtistRangeIndex)}
                className="text-xl font-bold text-black hover:opacity-70"
                aria-label="Next range"
              >
                <img
                  src={FrontArrow}
                  className="ring-2 px-1 hover:bg-gray-200 bg-white pl-1.5"
                  alt="Next"
                />
              </button>
            </div>
            {currentArtistRange === "3months" && (
              <>
                {catArt3Month ? (
                  <div className="flex flex-col gap-y-4 items-center">
                    {[...Array(3)].map((_, index) => (
                      <SplitRatioBarSkeleton key={index} />
                    ))}
                  </div>
                ) : userCatalogueArtist3MonthData.length !== 0 ? (
                  <div className="w-full max-h-[280px] overflow-y-scroll">
                    <div className="flex flex-col gap-y-4 items-center text-zinc-950 px-2 pt-[2px] pb-6">
                      {userCatalogueArtist3MonthData.map((item, index) => (
                        <SplitRatioBar
                          key={index}
                          itemName={item.Name}
                          Artist={item.Artist as string}
                          valueA={item.Playcounts[0]}
                          valueB={item.Playcounts[1]}
                          ArtistUrl={item.ArtistUrl as string}
                          itemUrl={item.EntryUrl as string}
                        />
                      ))}
                    </div>
                  </div>
                ) : null}
              </>
            )}

            {currentArtistRange === "12months" && (
              <>
                {catArt1Year ? (
                  <div className="flex flex-col gap-y-4 items-center">
                    {[...Array(3)].map((_, index) => (
                      <SplitRatioBarSkeleton key={index} />
                    ))}
                  </div>
                ) : userCatalogueArtist1YearData.length !== 0 ? (
                  <div className="w-full max-h-[280px] overflow-y-scroll">
                    <div className="flex flex-col gap-y-4 items-center text-zinc-950 px-2 pt-[2px] pb-6">
                      {userCatalogueArtist1YearData.map((item, index) => (
                        <SplitRatioBar
                          key={index}
                          itemName={item.Name}
                          Artist={item.Artist as string}
                          valueA={item.Playcounts[0]}
                          valueB={item.Playcounts[1]}
                          ArtistUrl={item.ArtistUrl as string}
                          itemUrl={item.EntryUrl as string}
                        />
                      ))}
                    </div>
                  </div>
                ) : null}
              </>
            )}

            {currentArtistRange === "1months" && (
              <>
                {catArt1Month ? (
                  <div className="flex flex-col gap-y-4 items-center">
                    {[...Array(3)].map((_, index) => (
                      <SplitRatioBarSkeleton key={index} />
                    ))}
                  </div>
                ) : userCatalogueArtist1MonthData.length !== 0 ? (
                  <div className="w-full max-h-[280px] overflow-y-scroll">
                    <div className="flex flex-col gap-y-4 items-center text-zinc-950 px-2 pt-[2px] pb-6">
                      {userCatalogueArtist1MonthData.map((item, index) => (
                        <SplitRatioBar
                          key={index}
                          itemName={item.Name}
                          Artist={item.Artist as string}
                          valueA={item.Playcounts[0]}
                          valueB={item.Playcounts[1]}
                          ArtistUrl={item.ArtistUrl as string}
                          itemUrl={item.EntryUrl as string}
                        />
                      ))}
                    </div>
                  </div>
                ) : null}
              </>
            )}
          </section>

          <section className="w-full flex flex-col">
            <div className="flex items-center justify-center gap-4 mb-4">
              <button
                onClick={() => goPrev(setCurrentTrackRangeIndex)}
                className="text-xl font-bold text-black hover:opacity-70"
                aria-label="Previous range"
              >
                <img
                  src={BackArrow}
                  className="ring-2 pr-1.5 hover:bg-gray-200 bg-white px-1"
                  alt="Previous"
                ></img>
              </button>

              <h2 className="text-lg font-bold text-black text-center min-w-[220px]">
                {trackRangeLabel[currentTrackRange]}
              </h2>

              <button
                onClick={() => goNext(setCurrentTrackRangeIndex)}
                className="text-xl font-bold text-black hover:opacity-70"
                aria-label="Next range"
              >
                <img
                  src={FrontArrow}
                  className="ring-2 px-1 hover:bg-gray-200 bg-white pl-1.5"
                  alt="Next"
                ></img>
              </button>
            </div>
            {currentTrackRange === "3months" && (
              <>
                {catTrack3Month ? (
                  <div className="flex flex-col gap-y-4 items-center">
                    {[...Array(3)].map((_, index) => (
                      <SplitRatioBarSkeleton key={index} />
                    ))}
                  </div>
                ) : userCatalogueTrack3MonthData.length !== 0 ? (
                  <div className="w-full max-h-[280px] overflow-y-scroll">
                    <div className="flex flex-col gap-y-4 items-center text-zinc-950 px-2 pt-[2px] pb-6">
                      {userCatalogueTrack3MonthData.map((item, index) => (
                        <SplitRatioBar
                          key={index}
                          itemName={item.Name}
                          Artist={item.Artist as string}
                          valueA={item.Playcounts[0]}
                          valueB={item.Playcounts[1]}
                          ArtistUrl={item.ArtistUrl as string}
                          itemUrl={item.EntryUrl as string}
                        />
                      ))}
                    </div>
                  </div>
                ) : null}
              </>
            )}

            {currentTrackRange === "12months" && (
              <>
                {catTrack1Year ? (
                  <div className="flex flex-col gap-y-4 items-center">
                    {[...Array(3)].map((_, index) => (
                      <SplitRatioBarSkeleton key={index} />
                    ))}
                  </div>
                ) : userCatalogueTrack1YearData.length !== 0 ? (
                  <div className="w-full max-h-[280px] overflow-y-scroll">
                    <div className="flex flex-col gap-y-4 items-center text-zinc-950 px-2 pt-[2px] pb-6">
                      {userCatalogueTrack1YearData.map((item, index) => (
                        <SplitRatioBar
                          key={index}
                          itemName={item.Name}
                          Artist={item.Artist as string}
                          valueA={item.Playcounts[0]}
                          valueB={item.Playcounts[1]}
                          ArtistUrl={item.ArtistUrl as string}
                          itemUrl={item.EntryUrl as string}
                        />
                      ))}
                    </div>
                  </div>
                ) : null}
              </>
            )}
            {currentTrackRange === "1months" && (
              <>
                {catTrack1Month ? (
                  <div className="flex flex-col gap-y-4 items-center">
                    {[...Array(3)].map((_, index) => (
                      <SplitRatioBarSkeleton key={index} />
                    ))}
                  </div>
                ) : userCatalogueTrack1MonthData.length !== 0 ? (
                  <div className="w-full max-h-[280px] overflow-y-scroll">
                    <div className="flex flex-col gap-y-4 items-center text-zinc-950 px-2 pt-[2px] pb-6">
                      {userCatalogueTrack1MonthData.map((item, index) => (
                        <SplitRatioBar
                          key={index}
                          itemName={item.Name}
                          Artist={item.Artist as string}
                          valueA={item.Playcounts[0]}
                          valueB={item.Playcounts[1]}
                          ArtistUrl={item.ArtistUrl as string}
                          itemUrl={item.EntryUrl as string}
                        />
                      ))}
                    </div>
                  </div>
                ) : null}
              </>
            )}
          </section>
        </div>
      </div>
    </div>
  );
}

const fetchBlendPercentage = async (label) => {
  await new Promise((r) => setTimeout(r, 500));
};
