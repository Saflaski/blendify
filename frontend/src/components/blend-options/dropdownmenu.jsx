import React from "react";

export function DropDownMenu() {
  return (
    <div className="flex items-center justify-center bg-inherit">
      <div className="flex gap-3 border-2 border-black bg-neutral-100 p-3">
        <TileButton>
          <img className="scale-130" src="/src/assets/images/dddefault.svg" />
        </TileButton>
        <TileButton>
          <img className="scale-110" src="/src/assets/images/ddcalendar.svg" />
        </TileButton>
        <TileButton>
          <img className="scale-120" src="/src/assets/images/ddartsong.svg" />
        </TileButton>
        <TileButton>
          <img className="scale-140" src="/src/assets/images/ddgenre.svg" />
        </TileButton>
      </div>
    </div>
  );
}

function TileButton({ children, label, onClick }) {
  return (
    <button
      type="button"
      onClick={onClick}
      className="group relative aspect-square w-14 select-none bg-neutral-200 p-3 text-neutral-900 outline-2 transition-all"
    >
      <div className="flex items-center justify-center flex-1 w-full">
        <div className="w-15 scale-110 flex items-center justify-center">
          {children}
        </div>
        {label ? (
          <span className="text-[10px] font-semibold tracking-wide text-neutral-800">
            {label}
          </span>
        ) : null}
      </div>
    </button>
  );
}
