export function BlendsButton({ onClick, mode, text }) {
  buttonInterior = getButtonInterior(mode, text);
  return (
    <button
      onClick={onClick}
      className="group relative aspect-square w-18.75 select-none bg-white p-3 outline-2 outline-black transition-all flex flex-col items-center justify-center gap-1"
    >
      <div className="flex items-center justify-center flex-1 w-full">
        <div className="w-3 h-3 flex items-center justify-center">
          {buttonInterior}
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

function getButtonInterior({ mode, text }) {
  if (mode == "new") {
    return <div>Add a new Blend</div>;
  } else if (mode == "existing") {
    return <div>{text}</div>;
  }
}
