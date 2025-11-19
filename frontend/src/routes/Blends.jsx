import { Button } from "@/components/ui/button";

export function Blend() {
  return (
    <div>
      <BlendButton
        onClick={async () => {
          console.log("Make New Blend clicked");
          const UA = "UUID1";
          const UB = "UUID2";
          await fetch(
            `http://localhost:3000/v1/blends/new/?UA=${UA}&UB=${UB}`,
            {
              method: "GET",
            },
          );
        }}
      >
        Make New Blend
      </BlendButton>
    </div>
  );
}

function BlendButton({ children, label, onClick }) {
  return (
    <button
      onClick={onClick}
      className="group relative aspect-square w-18.75 select-none bg-white p-3 outline-2 outline-black transition-all flex flex-col items-center justify-center gap-1"
    >
      <div className="flex items-center justify-center flex-1 w-full">
        <div className="w-6 h-3 flex items-center justify-center">
          {children}
        </div>
      </div>

      {label ? (
        <span className="text-[9px] font-semibold tracking-wide text-neutral-800 leading-none">
          {label}
        </span>
      ) : null}
    </button>
  );
}
