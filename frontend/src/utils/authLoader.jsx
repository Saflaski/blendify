import { redirect } from "react-router-dom";

export async function authLoader() {
  console.log("Auth Loader loaded");

  const res = await fetch("http://localhost:3000/v1/auth/validate", {
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
