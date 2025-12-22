// import { DropDownMenu } from "../components/blend-options/dropdownmenu";
import { ControlPanel } from "../components/blend-options/ControlPanel";
import { useLocation, useNavigate } from "react-router-dom";
import React, { useRef, useState, useEffect } from "react";
import { toBlob } from "html-to-image";

type BlendApiResponse = {
  usernames: string[];
  overallBlendNum: number;
  ArtistBlend: TypeBlend;
  AlbumBlend: TypeBlend;
  TrackBlend: TypeBlend;
};
type MetricKey = keyof BlendApiResponse;

type TypeBlend = {
  OneMonth: number;
  ThreeMonth: number;
  OneYear: number;
};

export function Blend() {
  // ------ If user is from invite link and not Add button -------
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const location = useLocation();

  const [blendId, setBlendId] = useState<string | null>(null);
  const [navLinkId, setNavLinkId] = useState<string | null>(null);
  const [userBlendData, setUserBlendData] = useState<BlendApiResponse | null>(
    null,
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

  useEffect(() => {
    //First check if blendid provided

    const getBlendIdFromInviteLink = async () => {
      //From URL Paste
      const params = new URLSearchParams(location.search);
      const urlInvite = params.get("invite");

      //From Add button

      const value = location.state;
      // const navigateInvite = value?.invite;

      const navigateInvite = (location.state as LocationState | null)?.value;

      //Log them
      console.log("urlInvite: ", urlInvite);
      console.log("Navigated Invite Link Data: ", navigateInvite);

      // if (!invite) {
      //   setError("Missing Invite Code");
      //   setLoading(false);
      //   return;
      // }

      const invite = navigateInvite ?? urlInvite;

      //Get blendid as authenticated user.
      const checkInvite = async () => {
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
          setBlendId(data["blendId"]);

          // setLoading(false);
        } catch (err) {
          console.error(err);
          setError("Something went wrong. Please try again.");
          setLoading(false);
        }
      };

      // If user clicked on existing blend from homepage

      if (blendId == null) {
        checkInvite();
      }
    };

    //
  }, []);

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
        const userData = JSON.parse(data) as BlendApiResponse;

        console.log("Parsed blend data:", userData);
        setUserBlendData(userData);
      } catch (err) {
        console.error(err);
        setError("Something went wrong. Please try again.");
        setLoading(false);
      }

      setLoading(false);
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

  return (
    <div className="w-full flex justify-center">
      <div className="w-full md:w-[60%] mx-auto text-center px-4 md:px-0 py-0 md:py-5">
        <div className="flex justify-left">
          <button
            type="button"
            className="inline-flex items-center gap-2 outline-2 outline-black font-[Roboto_Mono] font-bold border border-black/10 bg-white px-4 py-2 text-sm text-black shadow-sm hover:shadow md:text-base"
          >
            &lt; Your blends
          </button>
        </div>

        <div className="md:flex md:flex-wrap pr-2 mt-8 lg:grid lg:grid-cols-2 ">
          {/* --- Blendify Card ---*/}
          <div className="w-full flex justify-center md:mb-10 ">
            <div //Div to be screenshotted
              ref={captureRef}
              className="shine-element relative outline-2 outline-black bg-neutral-200 lg:w-80 md:w-50 h-auto p-10 aspect-2/3
             bg-size-[auto_200px] bg-[url(/src/assets/images/topography.svg)]"
            >
              {/* Copy button */}
              {!isCapturing && (
                <button
                  onClick={handleScreenshot}
                  className="absolute outline-1 active:bg-green-600  outline-black top-2 right-2 bg-inherit text-white px-1 py-1 "
                >
                  <img src="/src/assets/images/copy.svg" />
                </button>
              )}
              {/* Tooltip */}
              {copied && (
                <div
                  className=" absolute right-15 top-3 bg-gray-500 text-white 
              text-xs px-3 py-1 shadow-lg animate-fade-in-out"
                >
                  Copied!
                </div>
              )}
              {/* Hero number */}
              <h1 className="mt-0 text-6xl leading-none font-[Roboto_Mono] tracking-tight text-black md:text-4xl lg:text-7xl">
                {loading ? "--" : blendPercent}%
              </h1>
              {/* Big important text under the 80% */}
              <p className="mt-2 text-3xl md:text-3xl lg:text-4xl font-semibold text-gray-800">
                Ethan + Saf
              </p>
              {/* Filtering Mode */}
              <p className="mt-2 text-1xl md:text-1xl lg:text-1xl font-semibold text-gray-800">
                Default Mode
              </p>
              {/* Top Songs and Artists */}
              <div className="grid grid-row-2 gap-3 text-left text-black font-[Roboto_Mono] ">
                <ul>
                  <p className="font-black">Top Artists</p>
                  <li>Clairo</li>
                  <li>Men I Trust</li>
                  <li>Bring Me The Horizon</li>
                </ul>
                <ul>
                  <p className="font-black">Top Songs</p>
                  <li>Bababooey 2</li>
                  <li>Come Down</li>
                  <li>Bags</li>
                </ul>
              </div>
              <div className="flex justify-between gap-3 absolute bottom-3 left-1/2 -translate-x-1/1 size-12 h-auto">
                <img src="/src/assets/images/lastfm.svg" />

                <img src="/src/assets/images/apple.svg" />
              </div>
            </div>
          </div>
          {/* End of player card */}

          <div className=" flex flex-wrap justify-center items-center  lg:pl-10 gap-3">
            {/* Replace this block with <DropDownMenu /> if you already have it */}
            <ControlPanel setBlendPercent={setBlendPercent} />
          </div>
        </div>

        {/* Top blend artists section */}
        <section className="mt-12 text-left">
          <h2 className="text-xl md:text-2xl font-semibold text-black mb-4 text-center md:text-left">
            Top blend artists
          </h2>
          {/* Placeholder list/cards — replace with real data */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {["Artist One", "Artist Two", "Artist Three", "Artist Four"].map(
              (name) => (
                <div
                  key={name}
                  className="rounded-2xl border border-black/10 bg-white p-4 shadow-sm hover:shadow"
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium text-black">{name}</p>
                      <p className="text-sm text-black/60">
                        Blended frequently
                      </p>
                    </div>
                    <button className="rounded-lg border border-black/10 px-3 py-1 text-sm hover:shadow">
                      View
                    </button>
                  </div>
                </div>
              ),
            )}
          </div>
        </section>

        {/* Top blend songs section */}
        <section className="mt-12 text-left">
          <h2 className="text-xl md:text-2xl font-semibold text-black mb-4 text-center md:text-left">
            Top blend songs
          </h2>
          <div className="space-y-3">
            {[
              "Song A — Artist One",
              "Song B — Artist Two",
              "Song C — Artist Three",
            ].map((title) => (
              <div
                key={title}
                className="rounded-2xl border border-black/10 bg-white p-4 shadow-sm hover:shadow"
              >
                <div className="flex items-center justify-between">
                  <p className="text-black">{title}</p>
                  <button className="rounded-lg border border-black/10 px-3 py-1 text-sm hover:shadow">
                    Play
                  </button>
                </div>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}

const fetchBlendPercentage = async (label) => {
  await new Promise((r) => setTimeout(r, 500));
};
