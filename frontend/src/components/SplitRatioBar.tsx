import React from "react";
import { ur } from "zod/v4/locales";

type SplitRatioBarProps = {
  itemName: string;
  valueA: number;
  valueB: number;
  imageUrlA?: string;
  imageUrlB?: string;
  colorA?: string;
  colorB?: string;
  height?: string;
  urlToNavigateA: string;
  urlToNavigateB: string;
};

const SplitRatioBar: React.FC<SplitRatioBarProps> = ({
  itemName,
  valueA,
  valueB,
  imageUrlA = "https://lastfm.freetls.fastly.net/i/u/avatar170s/4711b0010c2035b2a26777f666cd3f3e.png",
  imageUrlB = "https://cdn-icons-png.flaticon.com/512/25/25231.png",
  colorA = "bg-[#987c4d]",
  colorB = "bg-[#ffedcb]",
  height = "h-10",
  urlToNavigateA,
  urlToNavigateB,
}) => {
  // var percentA: number = 0.0;
  // var percentB: number = 0.0;
  const total = valueA + valueB;
  const percentA = total === 0 ? 50 : (valueA / total) * 100;
  console.log("PERCENTS:");
  console.log(percentA);
  return (
    <div
      className={`relative w-full ${height} ring-2  ${colorB} pointer-events-none`}
    >
      <div
        className={`absolute left-0 top-0 h-full ${colorA}`}
        style={{ width: `${percentA}%` }}
      />

      {/* LEFT PIC WITH TOOLTIP */}
      <div className="absolute left-2 top-2  aspect-square group pointer-events-auto cursor-pointer">
        <button
          className="h-6 w-6 ring-2 ring-black"
          onClick={() => {
            console.log("LEFT CLICKED");
            window.open(urlToNavigateA, "_blank");
          }}
        >
          <img src={imageUrlA} alt="A" />
        </button>

        <div
          className="pointer-events-none
           absolute left-1/2 -translate-y-10
          opacity-0 group-hover:opacity-100 
          group-focus:opacity-100 transition
          bg-stone-900 text-stone-100 text-[12px] 
          font-mono px-2 py-0.5 z-50 text-nowrap"
        >
          {valueA} plays
        </div>
      </div>

      <div className="absolute right-2 top-2  aspect-square group pointer-events-auto cursor-pointer">
        <button
          className="h-6 w-6 ring-2 ring-black"
          onClick={() => {
            console.log("RIGHT CLICKED");
            window.open(urlToNavigateB, "_blank");
          }}
        >
          <img src={imageUrlB} alt="B" />
        </button>
        <div
          className="pointer-events-none
           absolute right-1/2 -translate-y-10
          opacity-0 group-hover:opacity-100 
          group-focus:opacity-100 transition
          bg-stone-900 text-stone-100 text-[12px] 
          font-mono px-2 py-0.5 z-50 text-nowrap"
        >
          {valueB} plays
        </div>
      </div>

      <div className="absolute inset-0 flex items-center justify-center">
        <span className="text-xs font-mono font-bold text-stone-900 tracking-tight">
          {itemName}
        </span>
      </div>
    </div>
  );
};

export default SplitRatioBar;
