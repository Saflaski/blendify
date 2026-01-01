import { Navigate, useLocation, useNavigate } from "react-router-dom";
import React, { useRef, useState, useEffect } from "react";
import { API_BASE_URL } from "../constants";
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

  return (
    <div className="w-full flex justify-center py-10">
      <div className="w-[50%] ring-2 ring-black px-6 py-8 space-y-10">
        <div className="flex flex-col items-center text-center space-y-4">
          <p className="text-black font-bold">Delete Account Data</p>
          <p className="text-black">
            Use the button to permanently delete everything regarding your
            account automatically.
            <br />
            This is an irreversible action.
          </p>
          <button
            className="p-2 bg-red-600 text-black font-bold font-[Roboto_Mono]"
            onClick={DeleteAccount}
          >
            Delete Account
          </button>
          <p className="text-black font-[Roboto_Mono] text-sm bg-amber-500">
            {error}
          </p>
        </div>

        <div className="border-t border-black" />

        <div className="flex flex-col items-center text-center space-y-3">
          <p className="text-black font-bold">Account Data Request</p>

          <p className="text-black">
            You may request a copy of all personal data associated with your
            account by contacting us at:
          </p>

          <p className="text-black font-bold">
            accounts@blendify.sabeehislam.com
          </p>

          <p className="text-black">
            According to GDPR regulations, this request must be completed within
            30 days.
          </p>
        </div>
      </div>
    </div>
  );
}
