import { Link, Navigate } from "react-router-dom";
import { useState, useEffect, useRef } from "react";
import { API_BASE_URL } from "../constants";
import HomeIcon from "@/assets/images/home.svg";
import MenuIcon from "@assets/images/menu.svg";
import BlendifyWhiteIcon from "@/assets/images/blendifyIconWhite.svg";
import BlendifyIcon from "@/assets/images/blendifyIcon.svg";
export function Navbar() {
  const [open, setOpen] = useState(false);
  const dropdownRef = useRef(null);
  const menuWrapperRef = useRef(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target) &&
        menuWrapperRef.current &&
        !menuWrapperRef.current.contains(event.target)
      ) {
        setOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <nav className="w-full bg-[#D84727] h-12 flex items-center z-20">
      <div className="lg:w-[full] md:w-[60%] w-[100%] mx-auto flex justify-between relative items-center px-2">
        {/* LEFT: HOME button */}
        <div className="relative flex items-center justify-start">
          <HomeButton onClick={Navigate("/home")} />
        </div>

        {/* CENTER: Blendify brand */}
        <div className="flex items-center  justify-center">
          <Link to="/home" className="no-underline focus:outline-none">
            <button
              type="button"
              className="
                [-webkit-tap-highlight-color:transparent]
                focus:outline-none
                font-extrabold
                text-2xl
                bg-[#D84727]
                flex items-center justify-center
                min-h-[5px]
                text-center
                text-[#F6E8CB]
                font-['Roboto_Mono','monospace']
                [text-shadow:2px_2px_0_#000]
                [drop-shadow:2px_2px_0_#000]
                transition-all duration-75 ease-in-out
                active:translate-x-[1px] active:translate-y-[1px]
                active:[text-shadow:1px_1px_0_#000]

              "
              // style={{ textShadow: "2px 2px 0 #000" }}
            >
              <img
                src={BlendifyWhiteIcon}
                className="
                w-7 h-auto
                aspect-square
                pr-1
                drop-shadow-[2px_2px_0_#000] 
                transition-all duration-75 ease-in-out
                active:translate-x-[1px] active:translate-y-[1px]
                active:drop-shadow-[1px_1px_0_#000]
              "
              />
              BLENDIFY
            </button>
          </Link>
        </div>

        {/* RIGHT: Add button */}
        <div ref={menuWrapperRef} className="flex items-center justify-end">
          <button
            type="button"
            onClick={() => setOpen((prev) => !prev)}
            className="
            bg-inherit
              text-black
              w-7 h-7
              flex items-center justify-center
            "
          >
            <img src={MenuIcon} alt="menu" className="w-full h-full" />
          </button>

          {open && (
            <div
              ref={dropdownRef}
              className="
                absolute
                top-[100%] right-0
                mt-2
                
                z-30
                inline-block
                bg-white
                px-10 py-4                
                ring-1
                ring-black
                shadow-md
                w-full
                h-dvh
                
                text-center
                lg:h-auto
                lg:text-left
                lg:w-40

                
                
              "
            >
              <ul className="list-none m-0 p-0 space-y-1">
                {/* <DropDownItem page="/" funcName={null} text="Home" /> */}
                <DropDownItem
                  page="/home"
                  setOpen={setOpen}
                  funcName={null}
                  text="Blends"
                />
                <DropDownItem
                  page="/account"
                  funcName={null}
                  setOpen={setOpen}
                  text="Account"
                />
                <DropDownItem
                  page="/about"
                  setOpen={setOpen}
                  funcName={null}
                  text="About"
                />

                <DropDownItem
                  page="/privacy"
                  setOpen={setOpen}
                  funcName={null}
                  text="Privacy"
                />

                <DropDownItem
                  page="/login"
                  funcName={handleLogOut}
                  setOpen={setOpen}
                  text="Log Out"
                />
              </ul>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}

// ---- helpers ----

async function handleLogOut() {
  await fetch(API_BASE_URL + "/auth/logout", {
    method: "POST",
    credentials: "include",
  });

  window.location.href = "/login";
}

function DropDownItem({ page, funcName, setOpen, text }) {
  return (
    <li className="font-['Roboto_Mono','monospace'] font-medium text-base bg-white">
      <Link to={page || "/"} className="no-underline">
        <button
          type="button"
          onClick={() => {
            funcName ? funcName() : null;
            setOpen(false);
          }}
          className={`bg-inherit ${
            text === "Log Out" ? "text-red-500" : "text-black"
          }`}
        >
          {text || "Lorem Ipsum"}
        </button>
      </Link>
    </li>
  );
}

function HomeButton({ onClick }) {
  return (
    <Link to={"/home"} className="no-underline">
      <button
        type="button"
        onClick={onClick}
        className="
        bg-transparent
        text-black
        w-7 h-7
        flex items-center justify-center
      "
      >
        <img src={HomeIcon} alt="home" className="w-full h-full" />
      </button>
    </Link>
  );
}
