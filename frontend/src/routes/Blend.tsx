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
  CatalogueTopItemsSchema,
  CatalogueTopItemsResponse,
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
const USER_A_TOP_ARTISTS_KEY = "USER_A_TOP_ARTISTS_KEY";
const USER_B_TOP_ARTISTS_KEY = "USER_B_TOP_ARTISTS_KEY";

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
  const [userATopItemsData, setUserATopItemsData] =
    useState<CatalogueTopItemsResponse>({ Items: [] });
  const [userBTopItemsData, setUserBTopItemsData] =
    useState<CatalogueTopItemsResponse>({ Items: [] });
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

  const [userATopItemsLoading, setUserATopItemsLoading] = useState(true);
  const [userBTopItemsLoading, setUserBTopItemsLoading] = useState(true);
  // const [userATopArtists, setUserATopArtists] = useLocalStorageState<string[]>(
  //   USER_A_TOP_ARTISTS_KEY,
  //   [],
  // );
  // const [userBTopArtists, setUserBTopArtists] = useLocalStorageState<string[]>(
  //   USER_B_TOP_ARTISTS_KEY,
  //   [],
  // );

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

  const downloadTopItems = async (
    blendid: string,
    category: string,
    duration: string,
    username: string,
    index: number,
    // setData: (val: CatalogueTopItemsResponse) => void,
  ): Promise<void> => {
    console.log("Downloading top items for:", username, category, duration);
    switch (index) {
      case 0:
        setUserATopItemsLoading(true);
        break;
      case 1:
        setUserBTopItemsLoading(true);
        break;
      default:
        console.error("Invalid user index for top items:", index);
    }
    try {
      const params = {
        blendId: blendid,
        duration: duration,
        category: category,
        username: username,
      };
      const queryString = new URLSearchParams(params).toString();
      const res = await fetch(
        `${API_BASE_URL}/blend/usertopitems?${queryString}`,
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
        return;
      }

      const data = await res.json();
      console.log("Top items data received:", data);

      const userData = CatalogueTopItemsSchema.parse(data);
      // console.log("Parsed blend data:", userData);
      console.log("Setting data for:", username, userData);
      switch (index) {
        case 0:
          setUserATopItemsData(userData);
          setUserATopItemsLoading(false);
          break;
        case 1:
          setUserBTopItemsData(userData);
          setUserBTopItemsLoading(false);
          break;
        default:
          console.error("Invalid user index for top items:", index);
      }
      // setData(userData);
    } catch (err) {
      console.error(err);
      setError("Something went wrong. Please try again.");
      return;
    }
  };

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
  const [cardInfo, setCardInfo] = useState([[], []]);
  const [mode, setMode] = useState("Default");
  const [users, setUsers] = useState<string[]>(["", ""]);

  const props: ControlPanelProps = {
    blendid: blendId as string,
    setMode,
    setUsers,
    // setUserATopItemsData,
    // setUserBTopItemsData,
    setBlendPercent,
    userATopItemApiResponse: userATopItemsData,
    userBTopItemApiResponse: userBTopItemsData,
    blendApiResponse: userCardData,
    downloadTopItems: downloadTopItems,
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

  const genres = [
    "rock",
    "pop",
    "edm",
    "hiphop",
    "classical",
    "jazz",
    "country",
    "metal",
    "indie",
    "punk",
    "reggae",
    "blues",
    "folk",
    "disco",
  ];

  const [genreTracks, setGenreTracks] = useState<
    CatalogueBlendResponse[] | undefined
  >([]);

  const [enabledButtons, setEnabledButtons] = useState({});
  // useEffect(() => {
  //   const genre = genres[0];
  //   setEnabledButtons({ [genre]: true });
  // }, []);
  const toggleButton = (genre) => {
    //If there is only 1 genre enabled, do not allow disabling it
    // if (enabledButtons[genre] && getEnabledGenres().length <= 1) { //UNCOMMENT IF YOU ALWAYS WANT ONE ENABLED
    //   return;
    // }
    setEnabledButtons((prev) => {
      const newState = { ...prev, [genre]: !prev[genre] };
      console.log(
        "Enabled buttons:",
        Object.keys(newState).filter((g) => newState[g]),
      );

      // Update genreTracks based on new state
      const songs = getSongsBasedOffGenres(
        Object.keys(newState).filter((g) => newState[g]),
      );
      setGenreTracks(songs);
      console.log("User 1 year track data", userCatalogueTrack1YearData);
      return newState;
    });
    console.log("Enabled buttons: ", getEnabledGenres());
    var songs = getSongsBasedOffGenres(getEnabledGenres());
    setGenreTracks(songs);
    console.log("Songs based off genres: ", songs);
    console.log(
      "BRO",
      userCatalogueTrack1MonthData.filter((item) =>
        item.Name.includes("Wildflower".toLowerCase()),
      ),
    );
  };

  useEffect(() => {
    setGenreTracks(userCatalogueTrack1YearData);
  }, [userCatalogueTrack1YearData]);
  const getEnabledGenres = () => {
    return Object.keys(enabledButtons).filter((genre) => enabledButtons[genre]);
  };

  //MOCK FUNCTION
  const getSongsBasedOffGenres = (genres: string[]) => {
    console.log("Getting songs based off genres: ", genres);
    if (genres.length === 0) {
      console.log("Not returning any songs");
      return userCatalogueTrack1YearData; //Or return all
    } else if (genres.length === 3) {
      console.log("All genres");
      return [];
    } else if (genres.includes("edm") && genres.includes("rock")) {
      return userCatalogueTrack1YearData.filter((item) =>
        item.Name.toLowerCase().includes("Drown".toLowerCase()),
      );
    } else if (genres.includes("rock")) {
      return userCatalogueTrack1YearData.filter((item) =>
        item.Artist?.toLowerCase().includes("Linkin Park".toLowerCase()),
      );
    } else if (genres.includes("pop")) {
      return userCatalogueTrack1YearData.filter((item) =>
        item.Name.toLowerCase().includes("wildflower".toLowerCase()),
      );
    }
  };

  useEffect(() => {
    console.log(
      "BROOO",
      userCatalogueTrack1YearData.filter((item) =>
        item.Name.toLowerCase().includes("Crawling".toLowerCase()),
      ),
    );
  }, [userCatalogueTrack1YearData]);
  const [genreExpanded, setGenreExpanded] = useState(false);
  return (
    <div className="w-full ">
      <div className="w-full md:w-[60%] flex pt-4 flex-col md:flex-row gap-x-5 mx-auto text-center px-4 gap-y-4 md:px-0 py-0 md:py-5">
        {/* <div className="flex justify-left"></div> */}

        {/* <div className="md:flex md:flex-wrap pr-2 mt-8 lg:grid lg:grid-cols-2 "> Old*/}
        {/* LEFT CONTENT AREA */}
        {/* This was card bit was enhanced from its previous version with AI. I am not that artistic to have come up with it myself lmao. */}
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
          <div className="w-full flex justify-center p-8">
            <div
              ref={captureRef}
              className="shine-element relative border-[3px] border-black bg-neutral-200 w-72 md:w-80 aspect-[2/3] flex flex-col overflow-hidden shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]"
              style={{
                backgroundImage: `url(${CardBackground})`,
                backgroundSize: "cover",
              }}
            >
              <div className="flex border-b-[3px] border-black bg-white/70 backdrop-blur-md h-28">
                <div className="flex-[1.5] p-3 border-r-[3px] border-black flex flex-col justify-center bg-white/40">
                  <span className="text-[10px] font-black uppercase tracking-tighter mb-0.5 text-neutral-500">
                    Score
                  </span>
                  <h1
                    className="text-5xl font-black leading-none tracking-tighter text-black"
                    style={{ fontFamily: "'Roboto Mono', monospace" }}
                  >
                    {cardLoading ? "--" : blendPercent}
                    <span className="text-sm ml-0.5">%</span>
                  </h1>
                </div>

                <div className="flex-[3] flex flex-col">
                  <div className="flex-1 flex">
                    <div className="flex-1 flex border-b-[3px]  border-black">
                      <div className="flex-[2] px-3 flex flex-col justify-center items-center ">
                        <p className="text-[11px] font-black tracking-[0.15em] justify-center font-[Quantico] text-black leading-none mb-1.5">
                          BLENDIFY
                        </p>
                        <div className="flex items-center">
                          <img
                            src={LastfmIcon}
                            className="h-2.5 w-auto opacity-70"
                            alt="Lastfm"
                          />
                        </div>
                      </div>

                      {!isCapturing && (
                        <button
                          onClick={handleScreenshot}
                          className="flex-1 border-l-[3px] border-black bg-yellow-400 hover:bg-black group transition-colors flex items-center justify-center min-w-[50px]"
                          title="Copy to Clipboard"
                        >
                          <img
                            src={CopyIcon}
                            className="pointer-events-none w-5 h-5 group-hover:invert transition-all"
                            alt="Copy"
                          />
                        </button>
                      )}
                    </div>
                  </div>
                  <div className="flex-1 bg-black  text-white px-3 flex flex-col justify-center  border-black overflow-hidden">
                    <span className="text-[7px] uppercase font-bold text-neutral-400 leading-none mb-1">
                      Current Mode
                    </span>
                    <span className="text-[11px] font-bold uppercase truncate tracking-tight">
                      {mode}
                    </span>
                  </div>
                </div>
              </div>

              <div className="bg-white text-black border-b-[3px] border-black px-4 py-2 flex justify-between items-center font-[Roboto_Mono] text-[11px] font-black uppercase tracking-tighter">
                <span className="truncate max-w-[100px]">
                  {users ? users[0] : "You"}
                </span>

                <span className="bg-black text-white text-[8px] px-1.5 py-0.5 mx-2">
                  VS
                </span>

                <span className="truncate max-w-[100px] text-right">
                  {users ? users[1] : "Someone"}
                </span>
              </div>

              {mode == "Default mode" ? (
                <div className="flex-grow p-4 flex flex-col gap-4">
                  <div className="border-[2px] border-black bg-white px-2 pb-1.5 pt-0 shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]">
                    <p className="text-[10px] font-black text-black uppercase tracking-widest mb-1 mt-0 border-b-2 border-black inline-block">
                      Top Artists
                    </p>
                    <ul className="space-y-0.5">
                      {userCatalogueArtist3MonthData
                        .slice(0, 5)
                        .map((item, index) => (
                          <li
                            key={index}
                            className="text-[10px] font-bold text-black font-[Roboto_Mono] truncate flex items-center"
                          >
                            <span className="text-[10px] mr-2 text-neutral-400">
                              0{index + 1}
                            </span>
                            {item.Name}
                          </li>
                        ))}
                    </ul>
                  </div>

                  <div className="border-[2px] border-black bg-white px-2 pb-1.5 pt-0  shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]">
                    <p className="text-[10px] font-black text-black uppercase tracking-widest mb-1 mt-0 border-b-2 border-black inline-block">
                      Top Tracks
                    </p>
                    <ul className="space-y-0.5">
                      {userCatalogueTrack3MonthData
                        .slice(0, 5)
                        .map((item, index) => (
                          <li
                            key={index}
                            className="text-[10px] font-bold  text-black font-[Roboto_Mono] truncate flex items-center"
                          >
                            <span className="text-[10px] mr-2 text-neutral-400">
                              0{index + 1}
                            </span>
                            {item.Name}
                          </li>
                        ))}
                    </ul>
                  </div>
                </div>
              ) : (
                <div className="flex-grow-2 p-4 flex flex-col gap-4">
                  {/* TOP ARTISTS: SIDE BY SIDE */}
                  <div className="border-[2px] border-black bg-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] flex flex-col">
                    {/* Section Header */}
                    <div className="border-b-[2px] border-black px-2 py-1 bg-neutral-100 flex justify-between items-center">
                      <p className="text-[9px] font-black text-black uppercase tracking-widest">
                        {mode}
                      </p>
                      <span className="text-[8px] font-bold text-neutral-400 font-[Roboto_Mono]"></span>
                    </div>

                    {/* Split Content */}
                    <div className="grid grid-cols-2 divide-x-[2px] divide-black">
                      {/* User 1 Column */}

                      <div className="p-2 overflow-hidden">
                        <p className="text-[8px] font-bold text-neutral-400 mb-2 uppercase tracking-tighter truncate"></p>
                        {userATopItemsLoading ? (
                          <ul className="space-y-0.5">
                            {" "}
                            {Array.from({ length: 5 }).map((_, index) => (
                              <li
                                key={index}
                                className="h-3 mb-1 bg-gray-300 animate-pulse"
                              />
                            ))}
                          </ul>
                        ) : (
                          <ul className="space-y-1">
                            {userATopItemsData.Items?.slice(0, 10).map(
                              (item, index) => (
                                <li
                                  key={index}
                                  className="text-[10px] font-bold text-black font-[Roboto_Mono] truncate flex items-center"
                                >
                                  <span className="text-[8px] mr-1.5 opacity-30">
                                    {index + 1}
                                  </span>
                                  {item}
                                </li>
                              ),
                            )}
                          </ul>
                        )}
                      </div>

                      {/* User 2 Column */}
                      <div className="p-2 overflow-hidden">
                        <p className="text-[8px] font-bold text-neutral-400 mb-2 uppercase tracking-tighter truncate text-right"></p>
                        {userATopItemsLoading ? (
                          <ul className="space-y-0.5">
                            {" "}
                            {Array.from({ length: 5 }).map((_, index) => (
                              <li
                                key={index}
                                className="h-3 mb-1 bg-gray-300 animate-pulse"
                              />
                            ))}
                          </ul>
                        ) : (
                          <ul className="space-y-1">
                            {userBTopItemsData.Items?.slice(0, 10).map(
                              (item, index) => (
                                <li
                                  key={index}
                                  className="text-[10px] font-bold text-black font-[Roboto_Mono] truncate flex items-center"
                                >
                                  {/* <span className="text-[8px] mr-1.5 opacity-30">
                                  {index + 1}
                                </span> */}
                                  {item}
                                </li>
                              ),
                            )}
                          </ul>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              )}
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
        <div className=" md:w-[60%] outline-amber-600 flex flex-col flex-wrap items-center justify-baseline gap-y-5 mt-10">
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

          {/* New genre thingy  */}
          <section className="w-full flex flex-col mb-6 ">
            <div className="relative flex flex-col justify-center ring-2 p-2 ring-black   ">
              <div>
                <div className=" pb-10">
                  <div
                    className={`flex flex-wrap justify-center flex-row m-2 gap-3 px-[3px] ${genreExpanded ? "max-h-[260px]" : "max-h-[110px] overflow-y-scroll py-2"}`}
                  >
                    {genres.map((genre) => (
                      <button
                        key={genre}
                        onClick={() => toggleButton(genre)}
                        className={`flex flex-row group relative w-auto h-10 min-w-5 text-sm font-[Roboto_Mono] font-medium select-none ${"active:shadow-[2px_2px_0_0_black] active:translate-[2px] shadow-[4px_4px_0_0_black]"} ${enabledButtons[genre] ? "bg-[#D84727] text-slate-100 outline-[#000000]" : "bg-white text-slate-950 outline-black "}  p-3 outline-2 transition-all flex flex-col items-center justify-center gap-1`}
                      >
                        {genre}
                      </button>
                    ))}
                  </div>

                  <button
                    onClick={() => setGenreExpanded((prev) => !prev)}
                    className={`absolute right-4 w-20 h-8 flex items-center justify-center
                  bg-black  shadow-[2px_2px_0_0_black]
                  active:translate-[1px] active:shadow-[1px_1px_0_0_black]
                  transition-all text-white font-[Roboto_Mono] font-medium text-sm`}
                    aria-label="Toggle genre list height"
                  >
                    <p className="pr-1 pl-1.5">
                      {genreExpanded ? "Less" : "Expand"}{" "}
                    </p>
                    <span
                      className={`transition-transform duration-300 ${
                        genreExpanded ? "rotate-180" : "rotate-0"
                      }`}
                    >
                      ▼
                    </span>
                  </button>
                </div>

                <div className="flex flex-col max-h-[280px] pt-2 overflow-y-scroll">
                  <div className="flex flex-col gap-y-4 items-center text-zinc-950 px-2 pt-[2px] pb-6 ">
                    {genreTracks ? (
                      genreTracks.map((item, index) => (
                        <SplitRatioBar
                          key={index}
                          itemName={item.Name}
                          Artist={item.Artist as string}
                          valueA={item.Playcounts[0]}
                          valueB={item.Playcounts[1]}
                          ArtistUrl={item.ArtistUrl as string}
                          itemUrl={item.EntryUrl as string}
                        />
                      ))
                    ) : (
                      <p className="text-black font-[Roboto_Mono]">
                        No Music Found
                      </p>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </section>

          {/* New experimental dropdown bit */}
          <section className="w-full flex flex-col ring-2 p-2 ring-black">
            <div className="flex items-center justify-center gap-4 mb-4 ">
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

          <section className="w-full flex flex-col ring-2 p-2 ring-black">
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
