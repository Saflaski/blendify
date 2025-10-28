import React, { useState, useEffect } from "react";

export function NotFound() {

    const [fact, setFact] = useState("");

  useEffect(() => {
    fetch("https://uselessfacts.jsph.pl/today.json?language=en")
      .then((res) => res.json())
      .then((data) => setFact(data.text))
      .catch((err) => console.error("Error fetching fact:", err));
  }, []);

    return (
        <div>
        <div className=" center-div p-2  w-full">
            <p className="text-gray-950 text-6xl font-[Transcity]">
            Oops, this page does not seem to exist! 
            </p>
            <p className="text-gray-600 text-s">
                Anyway, here's the random fact of the day: {fact || "Loading.."}
            </p>
            
            </div >
            <footer className="text-sm text-gray-500 mt-4 border-t pt-2 align-bottom">
          Facts provided by{" "}
          <a
            href="https://uselessfacts.jsph.pl/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-500 hover:underline"
          >
            uselessfacts.jsph.pl
          </a>
        </footer>
            </div>
    )
}