import { Navigate, useLocation, useNavigate } from "react-router-dom";
import React, { useRef, useState, useEffect } from "react";
import { API_BASE_URL, FRONTEND_URL } from "../constants";
import { GetPermaLink } from "./Home";
export function Account() {
  const navigate = useNavigate();
  const [error, setError] = useState("");
  const DeleteAccount = async () => {
    try {
      const res = await fetch(`${API_BASE_URL}/auth/delete`, {
        method: "POST",
        credentials: "include",
      });
      if (res.status == 401) {
        navigate(
          `/login?redirectTo=${encodeURIComponent(location.pathname + location.search)}`,
        );
        return;
      }

      if (!res.ok) {
        console.log(
          "Could not complete request to delete account. Please try again or contact the email below to delete.",
        );
        setError(
          "Could not complete request to delete account. Please try again or contact the email below to delete.",
        );

        return;
      }

      navigate("/login", {
        state: {
          from: "settings",
          message: "Account deleted successfully",
        },
      });
    } catch (err) {
      console.log("Could not delete due to: ", err);
    } finally {
      console.log("Finished delete account");
    }
  };

  const [link, setLink] = useState("start");
  async function handleGetPermaLink() {
    const newLink = await GetPermaLink();
    setLink(newLink);
    console.log("Setting link as: ", newLink);
  }

  useEffect(() => {
    console.log("Getting permalink");
    handleGetPermaLink();
  }, []);

  const handleRefreshLink = async () => {
    const newLink = await RefreshPermaLink();
    setLink(newLink);
    console.log("Setting link as: ", newLink);
  };

  const RefreshPermaLink = async () => {
    console.log("Fetching NEW perma blend link");
    try {
      const baseURL = `${API_BASE_URL}/blend/refreshpermalink`;
      const url = new URL(baseURL);
      const response = await fetch(url, {
        method: "POST",
        credentials: "include",
      });
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
  };

  return (
    <div className="w-full flex justify-center py-12 px-4">
      <div className="w-full max-w-3xl ring-2 ring-black  px-8 py-10">
        <h2 className="text-center font-[Sora] text-black text-3xl font-bold">
          Account
        </h2>
        <div>
          <div className="flex flex-col items-center text-center space-y-4 py-5">
            <p className="text-black font-bold">Refresh Permanent Invite</p>
            <p className="text-black">
              Revoke and refresh your permanent invite link. Doing this will
              cause all your old invite links to stop working.
            </p>
            <div className="flex w-4/5 flex-col items-center space-x-4">
              <textarea
                name="newLink"
                value={link}
                readOnly={true}
                rows={1}
                className="resize-none 
        flex w-full border-2 border-[#bbc2cb] text-center bg-gray-100 px-3 py-2 
        
         overflow-hidden text-nowrap text-slate-600 text-sm font-bold placeholder:text-slate-400"
              />
              <button
                onClick={handleRefreshLink}
                className={`border-2 mt-2 border-black bg-[#FFD700] px-6 py-2 
        text-sm font-black uppercase tracking-tighter
        shadow-[4px_4px_0_0_#000]
        hover:bg-[#00CED1] 
        text-black
        active:shadow-none active:translate-x-[4px] active:translate-y-[4px]
        transition-all`}
              >
                Refresh Link
              </button>
            </div>
          </div>
        </div>
        <div className="border-t border-black" />
        <div className="flex flex-col items-center text-center space-y-4 pt-5">
          <p className="text-black font-bold">Delete Account Data</p>
          <p className="text-black">
            Use the button to permanently delete everything regarding your
            account automatically.
            <br />
            This is an <b>irreversible action.</b>
          </p>
          <button
            className={`border-2 mt-2 border-black bg-[#f00] px-6 py-2 
        text-sm font-black uppercase tracking-tighter
        shadow-[4px_4px_0_0_#000]
        hover:bg-[#00CED1] 
        text-black
        active:shadow-none active:translate-x-[4px] active:translate-y-[4px]
        transition-all`}
            onClick={DeleteAccount}
          >
            Delete Account
          </button>
          <p className="text-black font-[Roboto_Mono] text-sm bg-amber-500">
            {error}
          </p>
        </div>

        <div className="border-t border-black" />

        <div className="flex flex-col items-center text-center space-y-3 pt-4">
          <p className="text-black font-bold">Account Data Request</p>

          <p className="text-black">
            You may request a copy of all personal data associated with your
            account by contacting us at:
          </p>

          <p className="text-black font-bold">blendify@sabeehislam.com</p>

          <p className="text-black">
            According to GDPR regulations, this request must be completed within
            30 days.
          </p>
        </div>
      </div>
    </div>
  );
}
