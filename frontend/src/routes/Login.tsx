import { API_BASE_URL, FRONTEND_URL } from "../constants";
import LFMIcon from "@/assets/images/lastfm.svg";
import TidalIcon from "@/assets/images/tidal.svg";
import AppleIcon from "@/assets/images/apple2.svg";
import UnderConstruction from "@/assets/images/underConstruction.jpg";
import mobilePreview from "@/assets/images/mobilePreview.png";
import desktopPreview from "@/assets/images/desktopPreview.png";
import React from "react";
import { useLocation } from "react-router-dom";
import { useState } from "react";
export function Login() {
  const { state } = useLocation();
  const [message] = useState(state?.message);
  console.log("Message: ", message);
  return (
    <div className=" min-h-fit w-full  items-center justify-center">
      <div className=" min-h-fit w-full flex flex-col items-center">
        <div className=" flex flex-col p-2 w-full items-center pt-10">
          {" "}
          <Title />
          <div className="mt-1 flex flex-wrap justify-center w-1/2 gap-4 ">
            {" "}
            {loginButton("lastfm")}
            {loginButton("construction")}
            {/* {loginButton("tidal")} */}
          </div>
          {message && (
            <div className="font-[Roboto_Mono] text-red-500">{message}</div>
          )}
        </div>
        <div
          className="mx-auto p-2 mt-5  mb-4 flex flex-wrap
      shadow-[4px_4px_0_0_#000]
      hover:shadow-[2px_2px_0_0_#000]
      hover:translate-[2px]
      transition-all duration-100
       justify-center text-center border-2 border-black
        text-black font-[Roboto_Mono] bg-[#00CED1] 
        pointer-events-auto"
        >
          <a href="/privacy">
            Read our <b>privacy policy</b>
          </a>
        </div>
      </div>
      {/* Mobile and desktop preview */}
      <div className="flex flex-row flex-wrap items-end justify-center gap-10 my-10">
        <div className="flex flex-col items-center">
          <h3 className="font-[Roboto_Mono] font-medium text-black bg-[#FF8C00] ring-3 px-2 ring-black text-md mb-2">
            Mobile Preview
          </h3>
          <img
            src={mobilePreview}
            alt="Mobile Preview"
            className="w-[250px] md:w-[200px] h-auto border-4 border-black shadow-[6px_6px_0_0_#FF8C00]"
          />
        </div>
        <div className="flex flex-col items-center mb-10 ">
          <h3 className="font-[Roboto_Mono] font-medium text-black bg-[#FF8C00] ring-3 px-2 ring-black text-md mb-2">
            Desktop Preview
          </h3>
          <img
            src={desktopPreview}
            alt="Desktop Preview"
            className="md:w-[600px] w-[300px] h-auto border-4 border-black shadow-[6px_6px_0_0_#FF8C00]"
          />
        </div>
      </div>
    </div>
  );
}

function loginButton(platform) {
  const platforms = {
    lastfm: {
      src: LFMIcon,
      alt: "LastFM Logo",
      wrapperClass: "mx-8",
      handleFunc: handleLastFMClick,
    },
    apple: {
      src: { AppleIcon },
      alt: "Apple Logo",
      wrapperClass: "mx-3",
    },
    tidal: {
      src: { TidalIcon },
      alt: "Tidal Logo",
      wrapperClass: "mx-3",
    },
    construction: {
      alt: "Construction Logo",
      wrapperClass: "mx-3",
    },
  };

  const config = platforms[platform];

  return (
    <button
      className="
        relative
        mx-auto my-2
        md:my-4
        w-[18rem] h-[5rem]
        md:w-[15em] md:h-[4em]
        border-3 border-black
        bg-white
        text-black
        font-bold
        font-md
        font-[Roboto_Mono]
        p-[3px]
        overflow-hidden
        shadow-[6px_6px_0_0_#FF8C00]
        flex items-center justify-center
        cursor-pointer
        hover:shadow-[4px_4px_0_0_#000]
        hover:translate-[2px]
        transition-all duration-100 
        active:translate-x-[1px] active:translate-y-[1px]
      "
      onClick={config.handleFunc}
    >
      <div className={config.wrapperClass}>
        {platform !== "construction" ? (
          <img
            src={config.src ? config.src : null}
            alt={config.alt}
            className="object-center image-render-pixel bg-white my-40"
          />
        ) : (
          "More Coming Soon"
        )}
      </div>
    </button>
  );
}

function handleLastFMClick() {
  const returnTo = encodeURIComponent(`${FRONTEND_URL}}/home`);
  window.location.href = `${API_BASE_URL}/auth/login/lastfm?return_to=${returnTo}`;
}
function Title() {
  return (
    <div className="flex flex-col items-center justify-center p-8">
      <div className="bg-[#00CED1] border-4 border-black p-6 shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] mb-6">
        <h2 className="leading-none m-0 font-[Sora] text-6xl lg:text-8xl text-black font-black text-center uppercase tracking-tighter">
          Blendify
        </h2>
      </div>

      <p className="inline-block px-6 py-2 bg-[#FF8C00] border-[3px] border-black font-[Roboto_Mono] font-bold text-black text-xl lg:text-2xl text-center shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
        Make a blend, compare and share
      </p>
    </div>
  );
}
