import React, { useState, useEffect } from "react";
import '/src/assets/styles/auxPages.css'

export function NotFound() {

    const [fact, setFact] = useState("");

  useEffect(() => {
    fetch("https://uselessfacts.jsph.pl/today.json?language=en")
      .then((res) => res.json())
      .then((data) => setFact(data.text))
      .catch((err) => console.error("Error fetching fact:", err));
  }, []);

    return (
        <div className="text-center min-h-full">
            <p className="text-gray-950 text-6xl">
            Oops, this page does not seem to exist! 
            </p>
            <p className="text-gray-600 text-s mt-3">
                Anyway, here's the random fact of the day: {fact || "Loading.."}
            </p>
            <footer className="text-sm text-gray-500 border-t mt-10 position-fixed align-bottom">
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
        </div >
            
        
    )
}