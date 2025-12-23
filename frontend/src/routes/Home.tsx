import { BlendsButton } from "../components/BlendsButton";
import { useState, useEffect } from "react";
import tick from "/src/assets/images/tick.svg";
import cross from "/src/assets/images/cross.svg";
import { useNavigate } from "react-router-dom";
import React from "react";

type Blend = {
  blendid: string;
  value: number;
  user: string[];
  timestamp: string;
};

var sampleJson = `{
	"blends": [
		{
			"blendid": "7673f65c-ab37-4fec-a698-5a0528b9af4d",
			"value": 55,
			"user": [
				"test2002",
				"saflas"
			],
			"timestamp": "2025-12-17T01:54:11+05:30"
		},
    {
			"blendid": "7673f65c-ab37-4fec-a698-5a0528b9af4d",
			"value": 67,
			"user": [
				"test2002",
				"ethan"
			],
			"timestamp": "2025-13-14T01:54:11+05:30"
		},
    {
			"blendid": "7673f65c-ab37-4fec-a698-5a0528b9af4d",
			"value": 98,
			"user": [
				"test2002",
				"kia"
			],
			"timestamp": "2025-02-12T01:54:11+05:30"
		}
	]
}`;

export function Home() {
  const navigate = useNavigate();
  async function AddBlend(givenURL: URL) {
    let url: URL | RequestInfo;
    // let value: any;
    const parsedGivenURL = new URL(givenURL);
    const invite = parsedGivenURL.searchParams.get("invite");
    let value: number;

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
      value = data["blendId"];
    } catch (err) {
      console.error("API error:", err);
      return;
    }

    console.log("Adding new blend from Blend Add URL Value:", value);
    navigate("/blend", {
      state: {
        id: "blendid",
        value: value,
      },
    });
  }

  const [blends, setBlends] = useState<Blend[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchBlends() {
      try {
        const url = "http://localhost:3000/v1/blend/userblends";
        const res = await fetch(url, {
          method: "GET",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
          credentials: "include",
        });
        const json = await res.json();
        setBlends(json.blends);
      } catch (err) {
        console.error("Error fetching blends:", err);
      } finally {
        setLoading(false);
      }
    }

    fetchBlends();
  }, []);

  function navToBlendPage(blendid: string) {
    // const navigate = useNavigate();
    navigate("/blend", {
      state: {
        id: "blendid",
        value: blendid,
      },
    });
  }

  return (
    <div className="min-h-screen w-full flex items-start justify-center py-5 font-[Roboto_Mono]">
      <div className="w-full max-w-xl bg-white border border-slate-300 px-5 py-6 flex flex-col gap-y-4 text-slate-900">
        <header className="w-full flex flex-col gap-1">
          <h1 className="text-xl font-semibold tracking-tight">Your blends</h1>

          <p className="text-sm text-slate-500">
            Paste a Blendify link from someone to start a blend
          </p>
          <section className="w-full">
            <AddNewBlendBar AddBlend={AddBlend} />
          </section>

          <div className="w-1/2 border-t my-4 mx-auto justify-center"></div>
          <p className="text-sm text-slate-500">
            Generate a Blendify link and send it to someone
          </p>
          <section>
            <GenerateLink />
          </section>
        </header>

        <div className="text-xs text-slate-500">
          {blends.length} {blends.length == 1 ? "blend" : "blends"}
        </div>

        <section className="w-full flex flex-col gap-3">
          <RecentOrTop />
          <ListOfBlends
            setEachBlend={setBlends}
            funcNav={navToBlendPage}
            blends={blends}
            loading={loading}
          />
        </section>
      </div>
    </div>
  );
}

function BlendSkeleton() {
  return (
    <div className="flex items-center justify-between border border-slate-200 px-3 py-2 animate-pulse">
      <div className="h-3 w-40 bg-slate-200 rounded" />
      <div className="h-2 w-16 bg-slate-200 rounded" />
    </div>
  );
}

type ListOfBlendsProps = {
  setEachBlend: React.Dispatch<React.SetStateAction<Blend[]>>;
  funcNav: (blendid: string) => void;
  blends: Blend[];
  loading: boolean;
};

function ListOfBlends({
  setEachBlend,
  funcNav,
  blends,
  loading,
}: ListOfBlendsProps) {
  if (loading) {
    return (
      <div className="space-y-1">
        <BlendSkeleton />
        <BlendSkeleton />
        <BlendSkeleton />
      </div>
    );
  }

  async function handleDelete(blendIdToDelete: string) {
    try {
      const blendId = blendIdToDelete;
      var url = new URL("http://localhost:3000/v1/blend/delete");
      const response = await fetch(url, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          blendId: blendId,
        }),
      });
      if (!response.ok)
        throw new Error(`Backend request error: ${response.status}`);

      const data = await response.json();
      console.log("API response data:", data);
      setEachBlend(blends.filter((blend) => blend.blendid !== blendIdToDelete));
    } catch {
      console.error("Error deleting blend:", blendIdToDelete);
    }
  }

  function daysAgo(isoDate) {
    const now = new Date();
    const then = new Date(isoDate);
    const diff = Math.floor(
      (Number(now) - Number(then)) / (1000 * 60 * 60 * 24),
    );
    return diff;
  }

  return (
    <div className="w-full space-y-1">
      {blends.map((blend) => (
        <div
          key={blend.blendid}
          className="flex overflow-hidden w-full group relative"
        >
          <button
            className=" flex flex-1 w-full text-left items-center transition-all duration-300 ease-in-out
            justify-between border border-slate-200 px-3 py-2 hover:bg-slate-50"
            onClick={() => funcNav(blend.blendid)}
          >
            <span className="truncate font-['Roboto_Mono'] text-xs">
              {blend.user.join(" + ")} // {blend.value}%
            </span>

            <span className="text-[10px] text-right text-slate-400 ml-2 shrink-0">
              {daysAgo(blend.timestamp) === 0
                ? "added today"
                : `added ${daysAgo(blend.timestamp)}d ago`}
            </span>
          </button>
          <button
            onClick={() => {
              console.log("Deleting blend:", blend.blendid);
              handleDelete(blend.blendid);
            }}
            className="
          -transition-x-4
          opacity-0 w-0
          transition-transform duration-100 ease-in
          
          group-hover:opacity-100
          group-hover:translate-x-0
          group-hover:pointer-events-auto
          group-hover:w-auto
          group-focus-within:opacity-100
          group-focus-within:translate-x-0
          group-focus-within:pointer-events-auto
          group-focus:w-auto
          pointer-events-none
          group-hover:px-1
          hover:bg-red-100
          hover:border-1
          focus:border-1
          text-xs 
         text-white
        "
          >
            <img
              src="src/assets/images/delete.svg"
              className="bg-inherit"
              alt="Go to blend"
            />
          </button>
        </div>
      ))}
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
        {/* <button className="px-3 py-2 text-slate-500 hover:text-slate-900 transition">
          Top
        </button> */}
      </div>
    </div>
  );
}

function AddNewBlendBar({ AddBlend }) {
  const [value, setValue] = useState("");
  var prefix = `http://localhost:5173/blend/`;
  const isValid = (value: string) => {
    //Simple URL check for now. Change slice num and url for prod
    if (value.slice(0, prefix.length) == prefix) {
      return true;
    } else return false;
  };

  return (
    <div className="flex w-full gap-2">
      <div
        className={`flex w-full border border-slate-600 bg-white px-3 py-2 text-xs font-['Roboto_Mono'] focus:outline-none focus:border-slate-900`}
      >
        <textarea
          name="newBlend"
          placeholder={prefix}
          value={value}
          onChange={(e) => setValue(e.target.value)}
          rows={1}
          className="resize-none w-full focus:outline-none overflow-hidden text-nowrap flex
         
          "
        ></textarea>
        {value.length > 0 && (
          <img
            src={isValid(value) ? tick : cross}
            alt={isValid(value) ? "Valid" : "Invalid"}
            className="justify-end relative w-6 h-4 pl-1 align-middle content-center"
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
