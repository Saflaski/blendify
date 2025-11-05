import { DropDownMenu } from "../components/blend-options/dropdownmenu";
import React from "react";

export function Home() {
  return (
    <div className="w-full flex justify-center">
      <div className="w-full md:w-[60%] mx-auto text-center px-4 md:px-0 py-0 md:py-5">
        <div className="flex justify-left">
          <button
            type="button"
            className="inline-flex items-center gap-2 outline-2 font-[Roboto_Mono] font-bold border border-black/10 bg-white px-4 py-2 text-sm text-black shadow-sm hover:shadow md:text-base"
          >
            &lt; Your blends
          </button>
        </div>

        {/* Hero number */}
        <h1 className="mt-8 text-7xl leading-none font-[Roboto_Mono] tracking-tight text-black md:text-8xl lg:text-9xl">
          80
        </h1>

        {/* Buttons directly below 80% */}
        <div className="mt-6 flex flex-wrap items-center justify-center gap-3">
          {/* Replace this block with <DropDownMenu /> if you already have it */}
          <button className="rounded-xl bg-black text-white px-5 py-2 text-sm md:text-base hover:opacity-90">
            Blend more
          </button>
          <button className="rounded-xl bg-white border border-black/10 px-5 py-2 text-sm md:text-base hover:shadow">
            Share
          </button>
          <button className="rounded-xl bg-white border border-black/10 px-5 py-2 text-sm md:text-base hover:shadow">
            Details
          </button>
          {/* <DropDownMenu /> */}
        </div>

        {/* Big important text under the 80% */}
        <p className="mt-6 text-3xl md:text-4xl lg:text-5xl font-semibold text-black">
          X + Y
        </p>

        {/* Top blend artists section */}
        <section className="mt-12 text-left">
          <h2 className="text-xl md:text-2xl font-semibold text-black mb-4 text-center md:text-left">
            Top blend artists
          </h2>
          {/* Placeholder list/cards — replace with real data */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {["Artist One", "Artist Two", "Artist Three", "Artist Four"].map(
              (name) => (
                <div
                  key={name}
                  className="rounded-2xl border border-black/10 bg-white p-4 shadow-sm hover:shadow"
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium text-black">{name}</p>
                      <p className="text-sm text-black/60">
                        Blended frequently
                      </p>
                    </div>
                    <button className="rounded-lg border border-black/10 px-3 py-1 text-sm hover:shadow">
                      View
                    </button>
                  </div>
                </div>
              ),
            )}
          </div>
        </section>

        {/* Top blend songs section */}
        <section className="mt-12 text-left">
          <h2 className="text-xl md:text-2xl font-semibold text-black mb-4 text-center md:text-left">
            Top blend songs
          </h2>
          <div className="space-y-3">
            {[
              "Song A — Artist One",
              "Song B — Artist Two",
              "Song C — Artist Three",
            ].map((title) => (
              <div
                key={title}
                className="rounded-2xl border border-black/10 bg-white p-4 shadow-sm hover:shadow"
              >
                <div className="flex items-center justify-between">
                  <p className="text-black">{title}</p>
                  <button className="rounded-lg border border-black/10 px-3 py-1 text-sm hover:shadow">
                    Play
                  </button>
                </div>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}
