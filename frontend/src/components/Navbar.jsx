import { Link, Navigate } from "react-router-dom";
import { useState, useEffect, useRef } from "react";

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
    <nav className="w-full bg-[#E0AD46] h-12 flex items-center z-20">
      <div className="lg:w-[full] md:w-[60%] w-[100%] mx-auto flex justify-between relative items-center px-2">
        {/* LEFT: Menu button + dropdown */}
        <div className="relative flex items-center justify-start">
          <HomeButton onClick={Navigate("/home")} />
        </div>

        {/* CENTER: Blendify brand */}
        <div className="flex items-center justify-center">
          <Link to="/home" className="no-underline">
            <button
              type="button"
              className="
                bg-transparent
                font-bold
                text-2xl
                flex items-center justify-center
                min-h-[5px]
                text-center
                text-[#F6E8CB]
                transition-all duration-75 ease-in-out
                active:translate-x-[1px] active:translate-y-[1px]
              "
              style={{ textShadow: "2px 2px 0 #000" }}
            >
              Blendify
            </button>
          </Link>
        </div>

        {/* RIGHT: Add button */}
        <div ref={menuWrapperRef} className="flex items-center justify-end">
          <button
            type="button"
            onClick={() => setOpen((prev) => !prev)}
            className="
              bg-transparent
              text-black
              w-7 h-7
              flex items-center justify-center
            "
          >
            <img
              src="/src/assets/images/menu.svg"
              alt="menu"
              className="w-full h-full"
            />
          </button>

          {open && (
            <div
              ref={dropdownRef}
              className="
                absolute
                top-[100%] right-0
                mt-5
                z-30
                inline-block
                bg-white
                px-10 py-4                
                 outline-1 outline-[#bbb]
                shadow-md
              "
            >
              <ul className="list-none m-0 p-0 space-y-1">
                {/* <DropDownItem page="/" funcName={null} text="Home" /> */}
                <DropDownItem page="/blends" funcName={null} text="Blends" />
                <DropDownItem page="/about" funcName={null} text="About" />
                <DropDownItem page="/privacy" funcName={null} text="Privacy" />
                <DropDownItem
                  page="/login"
                  funcName={handleLogOut}
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
  await fetch("http://localhost:3000/v1/auth/logout", {
    method: "POST",
    credentials: "include",
  });

  window.location.href = "/login";
}

function DropDownItem({ page, funcName, text }) {
  return (
    <li className="font-['Roboto_Mono','monospace'] font-medium text-base bg-white">
      <Link to={page || "/"} className="no-underline">
        <button
          type="button"
          onClick={funcName || undefined}
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
        <img
          src="/src/assets/images/home.svg"
          alt="home"
          className="w-full h-full"
        />
      </button>
    </Link>
  );
}
