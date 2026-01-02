import { href, useNavigate } from "react-router-dom";

export function About() {
  const navigate = useNavigate();
  return (
    <div className="w-full flex justify-center py-12 px-4">
      <div className="w-full max-w-3xl ring-2 ring-black  px-8 py-10">
        <h2 className="text-center font-[Roboto_Mono] text-black text-3xl font-semibold">
          About
        </h2>

        <p className="text-center text-black font-[Roboto_Mono] py-4">
          Blendify was made as an inspiration from Spotify's blend feature. But
          made to be cross platform eventually as well as present the 'blend' as
          more than just a singular number.
        </p>
        <div className="border-t border-black py-2  "></div>
        <div className="text-center text-black font-[Roboto_Mono] text-sm opacity-80">
          <p className="flex flex-col sm:flex-row items-center justify-center gap-1 sm:gap-2">
            <span>Made by</span>

            <a
              href="https://sabeehislam.com"
              target="_blank"
              rel="noopener noreferrer"
              className="font-bold underline underline-offset-4 hover:opacity-100 transition"
            >
              Sabeeh Islam
            </a>

            <span className="hidden sm:inline">·</span>

            <a
              href="https://github.com/saflaski"
              target="_blank"
              rel="noopener noreferrer"
              className="underline underline-offset-4 hover:opacity-100 transition"
            >
              GitHub
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
