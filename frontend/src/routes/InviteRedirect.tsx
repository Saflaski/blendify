import { useEffect } from "react";
import { replace, useNavigate, useParams } from "react-router-dom";
import React from "react";
import { API_BASE_URL, BLEND_ID_KEY } from "../constants";
type InviteParams = {
  inviteCode: string;
};

type ApiResponse = {
  blendId: string;
};

export function InviteRedirect() {
  const { inviteCode } = useParams<InviteParams>();
  const navigate = useNavigate();

  useEffect(() => {
    if (!inviteCode) return;
    async function handleInvite() {
      try {
        const url = new URL(`${API_BASE_URL}/blend/add`);
        const res = await fetch(url, {
          method: "POST",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({
            value: inviteCode,
            type: "permanent",
          }),
        });
        if (res.status == 401) {
          navigate(
            `/login?redirectTo=${encodeURIComponent(location.pathname + location.search)}`,
          );
          return;
        }

        if (!res.ok) {
          throw new Error("Invalid invite");
        }

        const data = await res.json();
        console.log("API Home response data:", data);
        const blendId = data["blendId"];

        console.log("Adding new blend from Blend Add URL Value:", blendId);
        localStorage.setItem(BLEND_ID_KEY, blendId);
        navigate("/blend");
        // navigate("/blend", {
        //   state: { id: data.blendId },
        //   replace: true,
        // });
      } catch (err) {
        console.error(err);
        navigate("/home", { replace: true });
      }
    }
    handleInvite();
  }, [inviteCode, navigate]);

  return (
    <div>
      <p>Accepting Invite</p>
    </div>
  );
}
