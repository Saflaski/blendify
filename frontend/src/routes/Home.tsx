import { BlendsButton } from "../components/BlendsButton";
import { useState, useEffect, useMemo, useRef } from "react";
import tick from "@/assets/images/tick.svg";
import cross from "@/assets/images/cross.svg";
import { useNavigate } from "react-router-dom";
import React from "react";
import "/src/assets/styles/home.css";
import { API_BASE_URL, FRONTEND_URL } from "../constants";
import Delete from "@/assets/images/delete.svg";
import Copy from "@/assets/images/copy.svg";

const BLEND_ID_KEY = "blend_id";

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
  localStorage.clear();
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

      url = new URL(`${API_BASE_URL}/blend/add`);
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
      console.log("API Home response data:", data);
      const blendId = data["blendId"];

      console.log("Adding new blend from Blend Add URL Value:", blendId);
      localStorage.setItem(BLEND_ID_KEY, blendId);
      navigate("/blend");
      return;
    } catch (err) {
      console.error("API error:", err);
      return;
    } finally {
    }
  }

  const [blends, setBlends] = useState<Blend[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchBlends() {
      try {
        const url = `${API_BASE_URL}/blend/userblends`;
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

    async function fetchUserInfo() {
      console.log("Fetching user info");
      try {
        const url = `${API_BASE_URL}/blend/userinfo`;
        const res = await fetch(url, {
          method: "GET",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
          credentials: "include",
        });
        const json = await res.json();
        setStats({
          plays: json.playcount,
          artists: json.artist,
          tracks: json.track,
        });
        setUsername(json.username);
      } catch (err) {
        console.error("Error fetching userinfo:", err);
      } finally {
        setStatLoading(false);
      }
    }

    fetchUserInfo();
  }, []);

  function navToBlendPage(blendid: string) {
    // const navigate = useNavigate();
    console.log(
      "func: navToBlendPage - Navigating to blend page with blendid:",
      blendid,
    );
    localStorage.setItem(BLEND_ID_KEY, blendid);
    navigate("/blend", {
      // state: {
      //   id: "blendid",
      //   value: blendid,
      // },
    });
  }

  const [userName, setUsername] = useState("XYZ");
  const [statLoading, setStatLoading] = useState(true);
  const [stats, setStats] = useState({
    plays: 0,
    artists: 0,
    tracks: 0,
  });
  return (
    <div className="min-h-screen w-full flex items-start dark:text-white justify-center py-5 font-[Roboto_Mono]">
      <div
        className={`w-full max-w-xl bg-[#F3ECDC] dark:bg-[#242321] 
          dark:shadow-[3px_3px_#F6E8CB] shadow-[3px_3px_#000]  dark:border-[#EED5A0]
           border-2 border-black px-5 py-6 flex flex-col gap-y-4 
           text-slate-900 dark:text-[#dfdcd7]`}
      >
        <header className="w-full flex flex-col gap-1">
          {!statLoading ? (
            <section className="w-full flex flex-col gap-4 mb-6">
              <h1 className="text-3xl font-bold tracking-tight">
                Hi,{" "}
                <span className="text-slate-600 dark:text-slate-50">
                  {userName.toUpperCase()}
                </span>
              </h1>

              <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
                <div className="flex flex-col">
                  <span className="text-2xl font-semibold">
                    {Number(stats.plays).toLocaleString()}
                  </span>
                  <span className="text-xs uppercase tracking-wide text-slate-500">
                    Plays
                  </span>
                </div>

                <div className="flex flex-col">
                  <span className="text-2xl font-semibold">
                    {Number(stats.artists).toLocaleString()}
                  </span>
                  <span className="text-xs uppercase tracking-wide text-slate-500">
                    Artists
                  </span>
                </div>

                <div className="flex flex-col">
                  <span className="text-2xl font-semibold">
                    {Number(stats.tracks).toLocaleString()}
                  </span>
                  <span className="text-xs uppercase tracking-wide text-slate-500">
                    Tracks
                  </span>
                </div>

                <div className="flex flex-col">
                  <span className="text-2xl font-semibold">
                    {Number(blends.length).toLocaleString()}
                  </span>
                  <span className="text-xs uppercase tracking-wide text-slate-500">
                    Blends
                  </span>
                </div>
              </div>
            </section>
          ) : (
            <TopUserInfoSectionSkeleton />
          )}
          <h1 className="text-xl font-semibold tracking-tight">Your blends</h1>

          <p className="text-sm text-slate-500 dark:text-slate-300">
            Paste a Blendify link from someone to start a blend
          </p>
          <section className="w-full">
            <AddNewBlendBar AddBlend={AddBlend} />
          </section>

          <div className="w-1/2 border-t my-4 mx-auto justify-center"></div>
          <p className="text-sm text-slate-500  dark:text-slate-300">
            Generate a Blendify link and send it to someone
          </p>
          <section>
            <GenerateLink />
          </section>
        </header>

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

  const renderCategory = (title, blendsArray) => {
    if (blendsArray.length === 0) return null;

    return (
      <div w-full>
        <h3 className="text-xs pl-2 font-bold dark:text-gray-100  text-gray-700 mb-2">
          {title}
        </h3>
        <div className="space-y-1">
          {blendsArray.map((blend) => (
            <div
              key={blend.blendid}
              className="flex overflow-hidden shadow-[2px_2px_black] 
              active:shadow-[0px_0px_black] active:translate-0.5
              hover:translate-[-2px] hover:shadow-[4px_4px_black]
              dark:shadow-[2px_2px_#F6E8CB] dark:active:shadow-[0px_0px_#F6E8CB]
              dark:hover:shadow-[4px_4px_#F6E8CB]
              transition-all
              bg-[#00CED1] dark:bg-[#d84827] dark:text-black
              border-1  border-[#000] dark:border-[#FFF] group relative"
            >
              <button
                className=" flex flex-1 w-full text-left transition-all duration-300 ease-in-out
            justify-between  dark:border-slate-200 border-r border-slate-700 px-3 py-2
            items-stretch overflow-hidden
             hover:bg-slate-50 dark:hover:bg-slate-700 hover:text-black dark:hover:text-white"
                onClick={() => funcNav(blend.blendid)}
              >
                <span className="truncate font-['Roboto_Mono'] text-sm">
                  {blend.user.join(" + ")} // {blend.value}%
                </span>

                <span className="text-[12px] text-right text-slate-800 dark:text-slate-300   ml-2 shrink-0">
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

                flex items-center justify-center
      bg-red-50 dark:bg-slate-100
      text-white
      transition-all duration-200 ease-in-out
      
      w-10 opacity-100 pointer-events-auto border-y border-r border-slate-200
      
      lg:w-0 lg:opacity-0 lg:pointer-events-none lg:border-none
      lg:group-hover:w-10 lg:group-hover:opacity-100 lg:group-hover:pointer-events-auto lg:group-hover:border-y lg:group-hover:border-r lg:group-hover:border-slate-200
      lg:group-focus-within:w-10 lg:group-focus-within:opacity-100 lg:group-focus-within:pointer-events-auto
      
      hover:bg-red-100 dark:hover:bg-red-500
      
        "
              >
                <img src={Delete} className="" alt="Delete Blend" />
              </button>
            </div>
          ))}
        </div>
      </div>
    );
  };

  async function handleDelete(blendIdToDelete: string) {
    try {
      const blendId = blendIdToDelete;
      var url = new URL(`${API_BASE_URL}/blend/delete`);
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

  const categoriseBlendsByDate = (blends: Blend[]) => {
    {
      /* Example categorization logic based on date */
    }
    const categories: {
      today: Blend[];
      yesterday: Blend[];
      thisWeek: Blend[];
      older: Blend[];
    } = {
      today: [],
      yesterday: [],
      thisWeek: [],
      older: [],
    };
    const now = new Date();
    blends.forEach((blend) => {
      const blendDate = new Date(blend.timestamp);
      const diffInDays = Math.floor(
        (Number(now) - Number(blendDate)) / (1000 * 60 * 60 * 24),
      );
      if (diffInDays === 0) {
        categories.today.push(blend);
      } else if (diffInDays <= 7) {
        categories.thisWeek.push(blend);
      } else {
        categories.older.push(blend);
      }
    });
    Object.keys(categories).forEach((key) => {
      categories[key as keyof typeof categories].sort(
        (a, b) =>
          new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime(),
      );
    });
    return categories;
  };

  const categorizedBlends = useMemo(
    () => categoriseBlendsByDate(blends),
    [blends],
  );
  return (
    <div className="w-full space-y-1">
      {renderCategory("Today", categorizedBlends.today)}
      {renderCategory("This Week", categorizedBlends.thisWeek)}
      {renderCategory("Older", categorizedBlends.older)}
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

function TopUserInfoSectionSkeleton() {
  return (
    <section className="w-full flex flex-col gap-4 mb-6 animate-pulse">
      <div className="h-8 w-48 bg-slate-200 rounded-md" />

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="flex flex-col gap-2">
            <div className="h-7 w-20 bg-slate-200 rounded-md" />
            <div className="h-3 w-14 bg-slate-200 rounded-md" />
          </div>
        ))}
      </div>
    </section>
  );
}

function AddNewBlendBar({ AddBlend }) {
  const [value, setValue] = useState("");
  var prefix = `${FRONTEND_URL}/blend/`;
  const isValid = (value: string) => {
    //Simple URL check for now. Change slice num and url for prod
    if (value.slice(0, prefix.length) == prefix) {
      return true;
    } else return false;
  };

  return (
    <div className="flex  w-full gap-2">
      <div
        className={`flex w-full border border-slate-600
            dark:bg-[#1a1917] dark:border-slate-50 dark:text-white
            text-black bg-white  px-3 py-2 text-xs font-['Roboto_Mono']
            shadow-[2px_2px_#000] dark:shadow-[2px_2px_#F6E8CB] 
            focus:outline-none focus:border-slate-900`}
      >
        <textarea
          name="newBlend"
          placeholder={prefix}
          value={value}
          onChange={(e) => setValue(e.target.value)}
          rows={1}
          className="resize-none w-full dark:text-white dark:border-slate-50
          
           dark:bg-inherit bg-white text-black focus:outline-none 
           overflow-hidden text-nowrap flex
         
          "
        ></textarea>
        {value.length > 0 && !isValid(value) && (
          <img
            src={cross}
            alt={isValid(value) ? "Valid" : "Invalid"}
            className="justify-end relative w-6 h-4 pl-1 align-middle content-center"
          />
        )}
      </div>

      <button
        onClick={() => AddBlend(value)}
        className={`border 
          text-black
          transition-all duration-100
        border-slate-900 bg-[#D84727] shadow-[3px_3px_0px_0px_rgba(0,0,0,1)]
        active:shadow-none active:translate-x-[2px] active:translate-y-[2px]
        dark:border-[#e6e1d7] dark:bg-[#242321] dark:text-[#e6e1d7] dark:shadow-[3px_3px_0px_0px_#F6E8CB]
        dark:active:shadow-none
          px-4 py-2 text-xs font-['Roboto_Mono'] font-bold tracking-wide  
          focus:outline-none focus:border-black`}
      >
        ADD
      </button>
    </div>
  );
}

function GenerateLink() {
  const [link, setLink] = useState("");
  const [copied, setCopied] = useState(false);
  const hideTimer = useRef<number | null>(null);

  async function handleGenerateLink() {
    const newLink = await generateNewLinkSomehow(); // your async fn
    setLink(newLink);
  }

  useEffect(() => {
    handleGenerateLink();
  }, []);

  const handleCopy = async () => {
    setCopied(false);
    if (!link) return;
    await navigator.clipboard.writeText(link); // full URL
    setCopied(true);
    if (hideTimer.current !== null) {
      clearTimeout(hideTimer.current);
    }
    // clearTimeout(hideTimer.current);
    hideTimer.current = setTimeout(() => setCopied(false), 1400);
  };

  return (
    <div className="flex w-full gap-3 items-center">
      <textarea
        name="newLink"
        value={link}
        readOnly={true}
        rows={1}
        className="flex-1 text-[11px] sm:text-xs resize-none overflow-hidden text-nowrap 
      transition-all duration-200
      bg-[#f8f3e9] border-slate-900 dark:border-1 text-slate-900 shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]
       dark:bg-slate-700 dark:border-[#e6e1d7] dark:text-[#e6e1d7] dark:shadow-[2px_2px_0px_0px_#F6E8CB]
      focus:outline-none px-3 py-2 font-mono cursor-default"
      />

      <div className="relative flex gap-2">
        {copied && (
          <div className="absolute -top-8 left-1/2 -translate-x-1/2 bg-slate-900 dark:bg-[#F6E8CB] text-white dark:text-black text-[10px] px-2 py-1 rounded shadow-lg">
            Copied!
          </div>
        )}

        <button
          onClick={handleCopy}
          className="flex items-center justify-center border transition-all duration-100
        border-slate-900 bg-[#D84727] shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]
        active:shadow-none active:translate-x-[2px] active:translate-y-[2px]
        dark:border-[#e6e1d7] dark:bg-[#242321] dark:shadow-[2px_2px_0px_0px_#F6E8CB]
        dark:active:shadow-none
        px-4 py-2"
        >
          <svg
            className="fill-slate-900 dark:fill-[#F6E8CB]"
            xmlns="http://www.w3.org/2000/svg"
            height="18px"
            viewBox="0 -960 960 960"
            width="18px"
          >
            <path d="M360-240q-33 0-56.5-23.5T280-320v-480q0-33 23.5-56.5T360-880h360q33 0 56.5 23.5T800-800v480q0 33-23.5 56.5T720-240H360Zm0-80h360v-480H360v480ZM200-80q-33 0-56.5-23.5T120-160v-560h80v560h440v80H200Zm160-240v-480 480Z" />
          </svg>
        </button>

        <button
          onClick={handleGenerateLink}
          className="border font-mono text-xs font-bold tracking-widest transition-all duration-100
        border-slate-900 bg-[#D84727] shadow-[3px_3px_0px_0px_rgba(0,0,0,1)]
        active:shadow-none active:translate-x-[2px] active:translate-y-[2px]
        dark:border-[#e6e1d7] dark:bg-[#242321] dark:text-[#e6e1d7] text-black dark:shadow-[3px_3px_0px_0px_#F6E8CB]
        dark:active:shadow-none
        px-4 py-2"
        >
          REFRESH
        </button>
      </div>
    </div>
  );
}

async function generateNewLinkSomehow() {
  console.log("Fetching outward blend link");
  try {
    const baseURL = `${API_BASE_URL}/blend/generatelink`;
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
    return `${FRONTEND_URL}/blend/?invite=` + newLink;
  } catch (err) {
    console.error("API erorr: ", err);
    return "Error no API connection";
  }
}
