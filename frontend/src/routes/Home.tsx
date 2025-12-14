import { BlendsButton } from "../components/BlendsButton";
import { useState, useEffect } from "react";
import tick from "/src/assets/images/tick.svg";
import cross from "/src/assets/images/cross.svg";
import { useNavigate } from "react-router-dom";
import React from "react";

export function Home() {
  const navigate = useNavigate();
  let value: any;
  async function AddBlend(givenURL: URL) {
    let url: URL | RequestInfo;
    // let value: any;
    const parsedGivenURL = new URL(givenURL);
    const invite = parsedGivenURL.searchParams.get("invite");

    try {
      // const baseURL = "http://localhost:3000/v1/blend/add";
      // const params = new URLSearchParams({ value: invite ?? "" });

      url = new URL("http://localhost:3000/v1/blend/add");
      // url.search = params.toString();

      const response = await fetch(url, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          value: invite,
        }),
      });
      if (!response.ok)
        throw new Error(`Backend request error: ${response.status}`);

      const data = await response.json();
      console.log("API response data:", data);
      value = data["linkId"];
    } catch (err) {
      console.error("API error:", err);
      return;
    }

    // const parsed = new URL(givenURL);
    // const invite = parsed.searchParams.get("invite");

    console.log("Adding new blend from Blend Add URL Value:", value);
    navigate("/blend", { state: value });
  }

  return (
    <div className="min-h-screen w-full flex items-start justify-center py-5 font-[Roboto_Mono]">
      <div className="w-full max-w-xl bg-white border border-slate-300 px-5 py-6 flex flex-col gap-y-4 text-slate-900">
        <header className="w-full flex flex-col gap-1">
          <h1 className="text-xl font-semibold tracking-tight">Your blends</h1>
          <p className="text-sm text-slate-500">
            Generate a Blendify link and send it to someone
          </p>
          <section>
            <GenerateLink />
          </section>
          <div className="w-1/2 border-t my-4 mx-auto justify-center"></div>
          <p className="text-sm text-slate-500">
            Paste a Blendify link from someone to start a blend
          </p>
          <section className="w-full">
            <AddNewBlendBar AddBlend={AddBlend} />
          </section>
        </header>

        <div className="text-xs text-slate-500">
          23 blends — 3 added recently
        </div>

        <section className="w-full flex flex-col gap-3">
          <RecentOrTop />

          <div className="space-y-1.5 text-sm">
            <div className="flex items-center justify-between border border-slate-200 px-3 py-2 hover:bg-slate-50 transition">
              <span className="truncate font-['Roboto_Mono'] text-xs">
                Ethan + Saf // 50%
              </span>
              <span className="text-[10px] text-slate-400 ml-2 shrink-0">
                added 2d ago
              </span>
            </div>
            <div className="flex items-center justify-between border border-slate-200 px-3 py-2 hover:bg-slate-50 transition">
              <span className="truncate font-['Roboto_Mono'] text-xs">
                Laurence + Saf // 80%
              </span>
              <span className="text-[10px] text-slate-400 ml-2 shrink-0">
                added 5d ago
              </span>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}

function RecentOrTop() {
  return (
    <div className="w-full pl-2">
      <div className="flex border-b border-slate-300 text-xs font-['Roboto_Mono']">
        <button className="px-3 py-2 border-b-2 border-slate-900 font-bold">
          Recent
        </button>
        <button className="px-3 py-2 text-slate-500 hover:text-slate-900 transition">
          Top
        </button>
      </div>
    </div>
  );
}

function AddNewBlendBar({ AddBlend }) {
  const [value, setValue] = useState("");
  var prefix = `http://localhost:5173/blend/`;
  const isValid = (value) => {
    //Simple URL check for now. Change slice num and url for prod
    if (value.slice(0, prefix.length) == prefix) {
      return true;
    } else return false;
  };

  return (
    <div className="flex w-full gap-2">
      <div className="flex w-full border border-slate-600 bg-white px-3 py-2 text-xs font-['Roboto_Mono'] focus:outline-none focus:border-slate-900">
        <textarea
          name="newBlend"
          placeholder={prefix}
          value={value}
          onChange={(e) => setValue(e.target.value)}
          rows={1}
          className="resize-none w-full focus:outline-none overflow-hidden flex"
        ></textarea>
        {value.length > 0 && (
          <img
            src={isValid(value) ? tick : cross}
            alt={isValid(value) ? "Valid" : "Invalid"}
            className="justify-end relative w-4 h-4 align-middle content-center"
          />
        )}
      </div>

      <button
        onClick={() => AddBlend(value)}
        className="border border-slate-900 bg-amber-400 px-4 py-2 text-xs font-['Roboto_Mono'] font-bold tracking-wide hover:bg-amber-300 focus:outline-none focus:border-black"
      >
        Add
      </button>
    </div>
  );
}

function GenerateLink() {
  const [link, setLink] = useState("");

  async function handleGenerateLink() {
    const newLink = await generateNewLinkSomehow(); // your async fn
    setLink(newLink);
  }

  useEffect(() => {
    handleGenerateLink();
  }, []);

  const handleCopy = async () => {
    if (!link) return;
    await navigator.clipboard.writeText(link); // full URL
  };

  return (
    <div className="flex w-full gap-2">
      <textarea
        name="newLink"
        value={link}
        readOnly={true}
        rows={1}
        className="flex-1 text-[11px] sm:text-xs resize-none overflow-hidden text-nowrap  border opacity-90 border-slate-300 bg-slate-50 focus:outline-none focus:ring-0 focus:border-slate-300 px-3 py-2 text-xs font-['Roboto_Mono'] cursor-default"
      ></textarea>
      <button
        onClick={handleCopy}
        className="flex items-center justify-center border border-slate-900 bg-amber-400 px-4 py-2 text-xs font-['Roboto_Mono'] font-bold tracking-wide hover:bg-amber-300 focus:outline-none focus:border-black"
      >
        <img
          className="size-4"
          src="src/assets/images/copy.svg"
          alt="Copy URL"
        />
      </button>
      <button
        onClick={handleGenerateLink}
        className="border border-slate-900 bg-amber-400 px-4 py-2 text-xs font-['Roboto_Mono'] font-bold tracking-wide hover:bg-amber-300 focus:outline-none focus:border-black"
      >
        Refresh
      </button>
    </div>
  );
}

async function generateNewLinkSomehow() {
  console.log("Fetching outward blend link");
  try {
    const baseURL = "http://localhost:3000/v1/blend/generatelink";
    const url = new URL(baseURL);
    const response = await fetch(url, { credentials: "include" });
    if (!response.ok) {
      throw new Error(
        `Backend request error on generating new outward link. Status: ${response.status}`,
      );
    }
    const data = await response.json();
    const newLink = data["linkId"];
    console.log("API response data: ", data);
    console.log("Blend Link: ", newLink);
    return "http://localhost:5173/blend/?invite=" + newLink;
  } catch (err) {
    console.error("API erorr: ", err);
    return (
      "http://localhost:5173/blend/?invite=" + Math.floor(Math.random() * 1000)
    );
  }
}
