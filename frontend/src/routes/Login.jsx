import "/src/assets/styles/Login.css";
import { API_BASE_URL, FRONTEND_URL } from "../constants";
import LFMIcon from "@/assets/images/lastfm.svg";
import TidalIcon from "@/assets/images/tidal.svg";
import AppleIcon from "@/assets/images/apple2.svg";
import UnderConstruction from "@/assets/images/underConstruction.jpg";
import { useLocation } from "react-router-dom";
import { useState } from "react";
export function Login() {
  const { state } = useLocation();
  const [message] = useState(state?.message);
  console.log("Message: ", message);
  return (
    <div className="">
      <div className=" w-full justify-center min-h-screen flex flex-row">
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
        border border-black
        bg-white
        text-black
        font-bold
        font-md
        font-[Roboto_Mono]
        p-[3px]
        overflow-hidden
        shadow-[4px_4px_0_0_#e0ad46]
        flex items-center justify-center
        cursor-pointer
        hover:shadow-[2px_2px_0_0_#000]
        hover:translate-[2px]
        transition-all duration-100 ease-in-out
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
    <div className="font-[Roboto_Mono]">
      <h2 className="leading-none m-0 text-6xl lg:text-8xl text-black font-semibold justify-center text-center">
        Blendify
      </h2>
      <p className=" font-mono px-5 mt-4 text-black text-2xl text-center break-words">
        Make a blend, compare and share
      </p>
    </div>
  );
}
