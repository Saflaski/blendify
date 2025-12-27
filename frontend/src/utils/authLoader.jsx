import { redirect } from "react-router-dom";
import { API_BASE_URL } from "../constants";

export async function authLoader() {
  console.log("Auth Loader loaded");

  const res = await fetch(`${API_BASE_URL}/auth/validate`, {
    credentials: "include",
    headers: { Accept: "application/json" },
  });
  console.log(res);

  if (res.ok) {
    console.log("Valid cookie");

    return null;
  } else {
    console.log("Invalid cookie");
    throw redirect("/login/");
  }
}
