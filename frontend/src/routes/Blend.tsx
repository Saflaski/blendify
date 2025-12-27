// import { DropDownMenu } from "../components/blend-options/dropdownmenu";
import { ControlPanel } from "../components/blend-options/ControlPanel";
import { useLocation, useNavigate } from "react-router-dom";
import React, { useRef, useState, useEffect } from "react";
import CardBackground from "@/assets/images/topography.svg";
import CopyIcon from "@/assets/images/copy.svg";
import LastfmIcon from "@/assets/images/lastfm.svg";
import "@/assets/styles/fonts.css";
import "@/assets/styles/index.css";
import { toBlob } from "html-to-image";
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

export function Blend() {
  // ------ If user is from invite link and not Add button -------
  const [error, setError] = useState<string | null>(null);
  const [cardLoading, setCardLoading] = useState(true);
  const [catalogueLoading, setCatalogueLoading] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  const [blendId, setBlendId] = useState<string | null>(null);
  const [navLinkId, setNavLinkId] = useState<string | null>(null);
  const [userCardData, setUserCardData] = useState<CardApiResponse>(
    {} as CardApiResponse,
  );
  const [userCatalogueArtist3MonthData, setUserCatalogueArtist3MonthData] =
    useState<CatalogueBlendResponse[]>([]);
  const [userCatalogueArtist1YearData, setUserCatalogueArtist1YearData] =
    useState<CatalogueBlendResponse[]>([]);
  const [userCatalogueTrack1YearData, setUserCatalogueTrack1YearData] =
    useState<CatalogueBlendResponse[]>([]);
  const [userCatalogueTrack3MonthData, setUserCatalogueTrack3MonthData] =
    useState<CatalogueBlendResponse[]>([]);

  const [catArt1Year, setCatArt1Year] = useState(true);
  const [catArt3Month, setCatArt3Month] = useState(true);
  const [catTrack1Year, setCatTrack1Year] = useState(true);
  const [catTrack3Month, setCatTrack3Month] = useState(true);

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
  if (blendId === null) {
    console.log("Getting blendid from API");
    getBlendIdFromInviteLink();
  }

  //

  console.log("Final blendId to use: ", blendId);
  useEffect(() => {
    console.log("Getting data for blendId (1): ", blendId);
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

    if (blendId == null) {
      setError("Could not get blendid.");
      console.log("Blend ID is null, cannot get data?");
    } else {
      getCardBlendData();
      console.log("Getting card blend data");
    }
  }, [blendId]);

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
    setData: (data: any[]) => void,
    setLoading: (loading: boolean) => void,
    setError: (msg: string) => void,
  ) => {
    console.log("Getting data for blendId:", blendId);

    try {
      setLoading(true);

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
      setLoading(false);
    }
  };

  useEffect(() => {
    console.log("Loading user catalogue artist blend data:");

    if (!blendId) {
      setError("Could not get blendid.");
      console.log("Blend ID is null, cannot get data?");
      return;
    }

    getCatalogueBlendData(
      "3month",
      "artist",
      blendId,
      setUserCatalogueArtist3MonthData,
      // setCatalogueLoading,
      setCatArt3Month,
      setError,
    );
  }, [blendId]);

  useEffect(() => {
    console.log("Loading user catalogue artist blend data:");

    if (!blendId) {
      setError("Could not get blendid.");
      console.log("Blend ID is null, cannot get data?");
      return;
    }

    getCatalogueBlendData(
      "3month",
      "track",
      blendId,
      setUserCatalogueTrack3MonthData,
      // setCatalogueLoading,
      setCatTrack3Month,
      setError,
    );
  }, [blendId]);

  useEffect(() => {
    console.log("Loading user catalogue artist blend data:");

    if (!blendId) {
      setError("Could not get blendid.");
      console.log("Blend ID is null, cannot get data?");
      return;
    }

    getCatalogueBlendData(
      "1year",
      "artist",
      blendId,
      setUserCatalogueTrack1YearData,
      // setCatalogueLoading,
      setCatArt1Year,
      setError,
    );
  }, [blendId]);

  useEffect(() => {
    console.log("Loading user catalogue artist blend data:");

    if (!blendId) {
      setError("Could not get blendid.");
      console.log("Blend ID is null, cannot get data?");
      return;
    }

    getCatalogueBlendData(
      "1year",
      "track",
      blendId,
      setUserCatalogueTrack1YearData,
      // setCatalogueLoading,
      setCatTrack1Year,
      setError,
    );
  }, [blendId]);

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

  return (
    <div className="w-full ">
      <div className="w-full md:w-[60%] flex pt-4 flex-col md:flex-row gap-x-5 mx-auto text-center px-4 gap-y-4 md:px-0 py-0 md:py-5">
        {/* <div className="flex justify-left"></div> */}

        {/* <div className="md:flex md:flex-wrap pr-2 mt-8 lg:grid lg:grid-cols-2 "> Old*/}
        {/* LEFT CONTENT AREA */}
        <div className="  md:w-[40%] flex flex-col flex-wrap items-center justify-baseline gap-y-5">
          {/* Player card */}
          <div className="w-full flex justify-center  ">
            <div className="w-full flex justify-center ">
              <div
                ref={captureRef}
                className={`shine-element relative ring-2 ring-black bg-neutral-200 
    w-58 md:w-58 lg:w-66 p-4 aspect-[2/3]
    bg-size-[auto_120px] bg-[url(${CardBackground})]`}
              >
                {!isCapturing && (
                  <button
                    onClick={handleScreenshot}
                    className="absolute top-1 right-1 outline outline-black bg-inherit p-1"
                  >
                    <img src={CopyIcon} className="w-4 h-4" />
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
    leading-none  font-normal -left-2 tracking-tight text-black relative inline-block"
                  style={{ fontFamily: "'Filepile', sans-serif" }}
                >
                  {cardLoading ? "--" : blendPercent}
                  <span className="absolute bottom-1 -right-12 text-lg font-normal ">
                    /100
                  </span>
                </h1>

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

                <p className="text-center w-full font-[toxigenesis] text-[#404040] tracking-widest font-medium absolute bottom-2 left-1/2 -translate-x-1/2 text-[10px]  text-shadow-2xs">
                  BLENDIFY.FM
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
          {catArt3Month ? (
            <section className="w-full flex flex-col">
              <h2 className="text-lg md:text-lg font-bold text-black mb-4 text-center">
                ARTISTS - LAST 3 MONTHS
              </h2>

              <div className="flex flex-col gap-y-4 items-center">
                {[...Array(3)].map((_, index) => (
                  <SplitRatioBarSkeleton key={index} />
                ))}
              </div>
            </section>
          ) : userCatalogueArtist3MonthData.length != 0 ? (
            <section className="w-full flex flex-col">
              <h2 className="text-lg md:text-lg font-bold text-black mb-4 text-center">
                ARTISTS - LAST 3 MONTHS
              </h2>

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
            </section>
          ) : null}

          {/* Top blend songs section */}
          {catArt1Year ? (
            <section className="w-full flex flex-col">
              <h2 className="text-lg md:text-lg font-bold text-black mb-4 text-center">
                ARTISTS - LAST 1 YEAR
              </h2>

              <div className="flex flex-col gap-y-4 items-center">
                {[...Array(3)].map((_, index) => (
                  <SplitRatioBarSkeleton key={index} />
                ))}
              </div>
            </section>
          ) : userCatalogueArtist1YearData.length != 0 ? (
            <section className="w-full flex flex-col">
              <h2 className="text-lg md:text-lg font-bold text-black mb-4 text-center">
                ARTISTS - LAST 1 YEAR
              </h2>

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
            </section>
          ) : null}

          {catTrack3Month ? (
            <section className="w-full flex flex-col">
              <h2 className="text-lg md:text-lg font-bold text-black mb-4 text-center">
                TRACKS - LAST 3 MONTHS
              </h2>

              <div className="flex flex-col gap-y-4 items-center">
                {[...Array(3)].map((_, index) => (
                  <SplitRatioBarSkeleton key={index} />
                ))}
              </div>
            </section>
          ) : userCatalogueTrack3MonthData.length != 0 ? (
            <section className="w-full flex flex-col">
              <h2 className="text-lg md:text-lg font-bold text-black mb-4 text-center">
                TRACKS - LAST 3 MONTHS
              </h2>

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
            </section>
          ) : null}

          {catTrack1Year ? (
            <section className="w-full flex flex-col">
              <h2 className="text-lg md:text-lg font-bold text-black mb-4 text-center">
                TRACKS - LAST 1 YEAR
              </h2>

              <div className="flex flex-col gap-y-4 items-center">
                {[...Array(3)].map((_, index) => (
                  <SplitRatioBarSkeleton key={index} />
                ))}
              </div>
            </section>
          ) : userCatalogueTrack1YearData.length != 0 ? (
            <section className="w-full flex flex-col">
              <h2 className="text-lg md:text-lg font-bold text-black mb-4 text-center">
                TRACKS - LAST 1 YEAR
              </h2>

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
            </section>
          ) : null}
        </div>
      </div>
    </div>
  );
}

const fetchBlendPercentage = async (label) => {
  await new Promise((r) => setTimeout(r, 500));
};
