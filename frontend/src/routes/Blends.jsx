import { BlendsButton } from "../components/BlendsButton";

export function Blends() {
  return (
    <div className="min-h-screen w-full flex items-start justify-center py-5 font-[Roboto_Mono]">
      <div className="w-full max-w-xl bg-white border border-slate-300 px-5 py-6 flex flex-col gap-y-4 text-slate-900">
        <header className="w-full flex flex-col gap-1">
          <h1 className="text-xl font-semibold tracking-tight">Your blends</h1>
          <p className="text-sm text-slate-500">
            Paste a Blendify URL to make a blend with someone
          </p>
        </header>

        <section className="w-full">
          <AddNewBlendBar />
        </section>

        <div className="text-xs text-slate-500">
          23 blends — 3 added recently
        </div>

        <section className="w-full flex flex-col gap-3">
          <RecentOrTop />

          <div className="space-y-1.5 text-sm">
            <div className="flex items-center justify-between border border-slate-200 px-3 py-2 hover:bg-slate-50 transition">
              <span className="truncate font-['Roboto_Mono'] text-xs">
                https://blendify.fm/b/my-favorite-morning-mix
              </span>
              <span className="text-[10px] text-slate-400 ml-2 shrink-0">
                added 2d ago
              </span>
            </div>
            <div className="flex items-center justify-between border border-slate-200 px-3 py-2 hover:bg-slate-50 transition">
              <span className="truncate font-['Roboto_Mono'] text-xs">
                https://blendify.fm/b/focus-grooves
              </span>
              <span className="text-[10px] text-slate-400 ml-2 shrink-0">
                added 5d ago
              </span>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}

function RecentOrTop() {
  return (
    <div className="w-full pl-2">
      <div className="flex border-b border-slate-300 text-xs font-['Roboto_Mono']">
        <button className="px-3 py-2 border-b-2 border-slate-900 font-bold">
          Recent
        </button>
        <button className="px-3 py-2 text-slate-500 hover:text-slate-900 transition">
          Top
        </button>
      </div>
    </div>
  );
}

function AddNewBlendBar() {
  return (
    <div className="flex w-full gap-2">
      <textarea
        name="newBlend"
        placeholder="https://blendify.fm/new/"
        rows={1}
        className="flex-1 resize-none overflow-hidden border border-slate-400 bg-white px-3 py-2 text-xs font-['Roboto_Mono'] focus:outline-none focus:border-slate-900"
      ></textarea>

      <button className="border border-slate-900 bg-amber-400 px-4 py-2 text-xs font-['Roboto_Mono'] font-bold tracking-wide hover:bg-amber-300 focus:outline-none focus:border-black">
        ADD
      </button>
    </div>
  );
}
