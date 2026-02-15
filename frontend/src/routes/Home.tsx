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
  cached: boolean;
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
    const inviteValue = parsedGivenURL.searchParams.get("singleinvite");
    let inviteCode: string | undefined;
    let mode: string;
    console.log("GIVEN URL: ", givenURL);
    if (inviteValue == null) {
      //Try parsing the url for a perma link
      const segments = parsedGivenURL.pathname.split("/").filter(Boolean);
      inviteCode = segments.pop();
      console.log("Perm invite code:", inviteCode);
      mode = "permanent";
    } else {
      inviteCode = inviteValue ?? undefined;
      mode = "temporary";
    }
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
          value: inviteCode,
          type: mode,
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
        console.log("Fetched blends:", json.blends);
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

  function navToBlendPage(blendid: string, cached: boolean = false) {
    // const navigate = useNavigate();
    console.log(
      "func: navToBlendPage - Navigating to blend page with blendid:",
      blendid,
    );
    localStorage.setItem(BLEND_ID_KEY, blendid);
    navigate("/blend", {
      state: {
        cached: cached,
      },
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
    <div
      className="min-h-fit w-full flex items-start 
    
    
    justify-center py-5 font-[Roboto_Mono]"
    >
      <div
        className={`w-full max-w-xl slate-bg  border-2 border-slate-900 
          shadow-[4px_4px_0_0_#000] 
          px-5 py-6 flex flex-col gap-y-4 text-slate-900`}
      >
        <header className="w-full flex flex-col gap-1">
          {!statLoading ? (
            <section className="w-full  flex flex-col gap-6 mb-8 ">
              <h1 className="text-4xl font-black  font-[Roboto_Mono] uppercase tracking-tighter ">
                Hi, <span className="text-[#FF3E00]">{userName}</span>
              </h1>

              <div className="grid grid-cols-2 font-[Sora] sm:grid-cols-4 gap-4">
                <PlayStatElement
                  number={stats.plays}
                  label="Plays"
                  color="bg-[#A7F3D0]"
                />
                <PlayStatElement
                  number={stats.artists}
                  label="Artists"
                  color="bg-[#BAE6FD]"
                />
                <PlayStatElement
                  number={stats.tracks}
                  label="Tracks"
                  color="bg-[#DDD6FE]"
                />
                <PlayStatElement
                  number={blends.length}
                  label="Blends"
                  color="bg-[#FDE68A]"
                />
              </div>
            </section>
          ) : (
            <TopUserInfoSectionSkeleton />
          )}
          <h1 className="text-xl font-semibold tracking-tight">Your blends</h1>

          <p className="text-sm text-slate-500">
            Paste a Blendify link from someone to start a blend
          </p>
          <section className="w-full">
            <AddNewBlendBar AddBlend={AddBlend} />
          </section>

          <div className="w-1/2 border-t my-4 mx-auto justify-center"></div>
          <p className="text-sm text-slate-500">
            Generate a one-time Blendify link and send it to someone
          </p>
          <section>
            <GenerateLink />
          </section>

          <div className="w-1/6 border-t my-2 mx-auto justify-center"></div>
          <p className="text-sm text-slate-500">
            Or use the permanent link below to invite anyone anytime
          </p>
          <section>
            <GeneratePermaLink />
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

export const PlayStatElement = ({ number, label, color = "bg-white" }) => {
  return (
    <div
      className={`
      flex flex-col justify-center
      p-4 border-2 border-black 
      ${color} 
      shadow-[4px_4px_0_0_#000] 
      hover:shadow-none hover:translate-x-[2px] hover:translate-y-[2px] 
      transition-all
    `}
    >
      <span className="text-3xl font-black leading-none">
        {Number(number).toLocaleString()}
      </span>
      <span className="text-[10px] uppercase font-bold tracking-tighter mt-1 text-black/60">
        {label}
      </span>
    </div>
  );
};

type ListOfBlendsProps = {
  setEachBlend: React.Dispatch<React.SetStateAction<Blend[]>>;
  funcNav: (blendid: string, cached: boolean) => void;
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

  const renderCategory = (title, blendsArray: Blend[]) => {
    if (blendsArray.length === 0) return null;

    return (
      <div w-full>
        <h3 className=" pl-2 font-bold  text-gray-700 mb-2">{title}</h3>
        <div className="space-y-2 flex-1 ">
          {blendsArray.map((blend) => (
            <div
              key={blend.blendid}
              className="
    flex w-full group relative 
    bg-black border-[3px] border-black 
    shadow-[4px_4px_0px_0px_#000]
    active:translate-[2px]
    active:shadow-[2px_2px_0px_0px_#000]
    transition-all duration-100 mb-4
  "
            >
              <button
                className="
      flex-1 flex items-center justify-between
      bg-white p-3 text-left 
      hover:bg-yellow-50 active:bg-yellow-100
      transition-colors duration-75
      min-w-0 
    "
                onClick={() => funcNav(blend.blendid, blend.cached)}
              >
                <span className="truncate font-['Roboto_Mono'] text-sm font-black uppercase">
                  {blend.user.join(" + ")}{" "}
                  <span className="text-blue-600">//</span> {blend.value}%
                </span>

                <span className="ml-4 text-[10px] font-bold uppercase bg-black text-white px-2 py-1 shrink-0">
                  {daysAgo(blend.timestamp) === 0
                    ? "TODAY"
                    : `${daysAgo(blend.timestamp)}D`}
                </span>
              </button>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  handleDelete(blend.blendid);
                }}
                className="
      flex items-center justify-center bg-red-500
      hover:bg-red-600 active:bg-red-700
      border-l-[3px] border-black
      
      w-12 opacity-100
      
      lg:w-0 lg:opacity-0 lg:border-l-0
      lg:group-hover:w-14 lg:group-hover:opacity-100 lg:group-hover:border-l-[3px]
      
      transition-all duration-200 ease-in-out overflow-hidden
    "
              >
                <img
                  src={Delete}
                  className="md:w-7 md:h-7 invert min-w-[20px]"
                  alt="Delete"
                />
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
      <div className="h-8 w-48 bg-slate-200" />

      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="flex flex-col gap-2">
            <div className="h-7 w-20 bg-slate-200 " />
            <div className="h-3 w-14 bg-slate-200 " />
          </div>
        ))}
      </div>
    </section>
  );
}

function AddNewBlendBar({ AddBlend }) {
  const [value, setValue] = useState("");
  var prefix = `${FRONTEND_URL}/blend/`;
  var prefix2 = `${FRONTEND_URL}/invite/`;
  const isValid = (value: string) => {
    //Simple URL check for now. Change slice num and url for prod
    if (
      value.slice(0, prefix.length) == prefix ||
      value.slice(0, prefix2.length) == prefix2
    ) {
      return true;
    } else return false;
  };

  return (
    <div className="flex w-full gap-3 font-mono">
      <div
        className={`
        flex w-full border-2 border-black bg-white px-3 py-2 
        shadow-[4px_4px_0_0_#000] transition-all
        focus-within:translate-x-[-2px] focus-within:translate-y-[-2px] 
        focus-within:shadow-[6px_6px_0_0_#000]
      `}
      >
        <textarea
          name="newBlend"
          placeholder={`${FRONTEND_URL}/`}
          value={value}
          onChange={(e) => setValue(e.target.value)}
          rows={1}
          className="resize-none w-full focus:outline-none overflow-hidden text-nowrap bg-transparent text-sm font-bold placeholder:text-slate-400"
        />

        {value.length > 0 && !isValid(value) && (
          <img
            src={cross}
            alt="Invalid"
            className="w-5 h-5 self-center ml-2 invert"
          />
        )}
      </div>

      <button
        onClick={() => AddBlend(value)}
        className={`
        border-2 border-black bg-[#FFD700] px-6 py-2 
        text-sm font-black uppercase tracking-tighter
        shadow-[4px_4px_0_0_#000]
        hover:bg-[#00CED1] 
        
        active:shadow-none active:translate-x-[4px] active:translate-y-[4px]
        transition-all
      `}
      >
        Add
      </button>
    </div>
  );
}

function GeneratePermaLink() {
  const [copied, setCopied] = useState(false);
  const hideTimer = useRef<number | null>(null);
  const [link, setLink] = useState("");
  async function handleGetPermaLink() {
    const newLink = await GetPermaLink();
    setLink(newLink);
  }

  useEffect(() => {
    handleGetPermaLink();
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
    <div className="flex w-full gap-2">
      <textarea
        name="newLink"
        value={link}
        readOnly={true}
        rows={1}
        className="resize-none 
        flex w-full border-2 border-[#bbc2cb] bg-gray-100 px-3 py-2 
         
         overflow-hidden text-nowrap text-slate-600 text-sm font-bold placeholder:text-slate-400"
      ></textarea>

      <div className="relative">
        {copied && (
          <div
            className="absolute right-14.5 bg-gray-500 text-white 
        text-[10px] px-2 py-0.5 shadow animate-fade-in-out"
          >
            Copied!
          </div>
        )}
        <CopyButton onClick={handleCopy} />
      </div>
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
    <div className="flex flex-col w-full gap-2">
      <textarea
        name="newLink"
        value={link}
        readOnly={true}
        rows={1}
        className="resize-none 
        flex w-full border-2 border-[#bbc2cb] bg-gray-100 px-3 py-2 
        
         overflow-hidden text-nowrap text-slate-600 text-sm font-bold placeholder:text-slate-400"
      ></textarea>

      <div className="flex flex-row gap-2 items-center justify-end">
        <div className="relative">
          {copied && (
            <div
              className="absolute right-14.5 bg-gray-500 text-white 
        text-[10px] px-2 py-0.5 shadow animate-fade-in-out"
            >
              Copied!
            </div>
          )}
          <CopyButton onClick={handleCopy} />
        </div>
        <button
          onClick={handleGenerateLink}
          className={`border-2 border-black bg-[#FFD700] px-6 py-2 
        text-sm font-black uppercase tracking-tighter
        shadow-[4px_4px_0_0_#000]
        hover:bg-[#00CED1] 
        
        active:shadow-none active:translate-x-[4px] active:translate-y-[4px]
        transition-all`}
        >
          Refresh
        </button>
      </div>
    </div>
  );
}

export const CopyButton = ({ onClick }) => {
  return (
    <button
      onClick={onClick}
      className={`border-2 border-black bg-[#FFD700] px-4 py-2 
        text-sm font-black uppercase tracking-tighter
        shadow-[4px_4px_0_0_#000]
        hover:bg-[#00CED1] 
        
        active:shadow-none active:translate-x-[4px] active:translate-y-[4px]
        transition-all`}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        height="18px"
        viewBox="0 -960 960 960"
        width="24px"
        fill="#000"
      >
        <path d="M360-240q-33 0-56.5-23.5T280-320v-480q0-33 23.5-56.5T360-880h360q33 0 56.5 23.5T800-800v480q0 33-23.5 56.5T720-240H360Zm0-80h360v-480H360v480ZM200-80q-33 0-56.5-23.5T120-160v-560h80v560h440v80H200Zm160-240v-480 480Z" />
      </svg>
    </button>
  );
};

export async function GetPermaLink() {
  console.log("Fetching outward perma blend link");
  try {
    const baseURL = `${API_BASE_URL}/blend/getpermalink`;
    const url = new URL(baseURL);
    const response = await fetch(url, { credentials: "include" });
    if (response.status == 429) {
      console.log("Error: Rate limit exceeded");
      return "Woah calm down";
    }
    if (!response.ok) {
      throw new Error(
        `Backend request error on generating new outward link. Status: ${response.status}`,
      );
    }
    const data = await response.json();
    const newLink = data["permaLinkId"];
    console.log("Perma blend Link: ", newLink);
    return `${FRONTEND_URL}/invite/` + newLink;
  } catch (err) {
    console.error("API erorr: ", err);
    return "Error no API connection";
  }
}

async function generateNewLinkSomehow() {
  console.log("Fetching outward blend link");
  try {
    const baseURL = `${API_BASE_URL}/blend/generatelink`;
    const url = new URL(baseURL);
    const response = await fetch(url, { credentials: "include" });
    if (response.status == 429) {
      console.log("Error: Rate limit exceeded");
      // setTimeout(() => {
      //   generateNewLinkSomehow();
      // }, 1000);
      // return;
      return "Woah calm down";
    }
    if (!response.ok) {
      throw new Error(
        `Backend request error on generating new outward link. Status: ${response.status}`,
      );
    }
    const data = await response.json();
    const newLink = data["linkId"];
    console.log("API response data: ", data);
    console.log("Blend Link: ", newLink);
    return `${FRONTEND_URL}/blend/?singleinvite=` + newLink;
  } catch (err) {
    console.error("API erorr: ", err);
    return "Error no API connection";
  }
}
