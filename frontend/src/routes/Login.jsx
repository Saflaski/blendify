import "/src/assets/styles/Login.css";

export function Login() {
  return (
    <div className="">
      {/* <Navbar /> */}
      <div className=" w-full justify-center min-h-screen flex flex-row">
        <div className=" flex flex-col p-2 w-full items-center pt-10">
          {" "}
          {/* Added flex-col and items-center for vertical centering */}
          <Title />
          {/* New section for login buttons */}
          <div className="mt-1 flex flex-wrap justify-center w-1/2 gap-4 ">
            {" "}
            {/* Using flex-wrap and gap for responsiveness and spacing */}
            {loginButton("lastfm")}
            {loginButton("apple")}
            {loginButton("tidal")}
          </div>
        </div>
      </div>
    </div>
  );
}

function loginButton(platform) {
  const platforms = {
    lastfm: {
      src: "/src//assets/images/lastfm.svg",
      alt: "LastFM Logo",
      wrapperClass: "mx-8",
      handleFunc: handleLastFMClick,
    },
    apple: {
      src: "/src//assets/images/apple2.svg",
      alt: "Apple Logo",
      wrapperClass: "mx-3",
    },
    tidal: {
      src: "/src//assets/images/tidal.svg",
      alt: "Tidal Logo",
      wrapperClass: "mx-3",
    },
  };

  const config = platforms[platform];

  return (
    <button
      className="
        relative
        mx-auto my-2
        md:my-4
        w-[18rem] h-[5rem]
        md:w-[15em] md:h-[4em]
        border border-black
        bg-white
        p-[3px]
        overflow-hidden
        shadow-[4px_4px_0_0_#e0ad46]
        transition-all duration-75 ease-in-out
        flex items-center justify-center
        cursor-pointer
        hover:shadow-[2px_2px_0_0_#000]
        active:translate-x-[1px] active:translate-y-[1px]
      "
      onClick={config.handleFunc}
    >
      <div className={config.wrapperClass}>
        <img
          src={config.src}
          alt={config.alt}
          className="object-center image-render-pixel bg-white my-40"
        />
      </div>
    </button>
  );
}

function handleLastFMClick() {
  const returnTo = encodeURIComponent("http://localhost:5173/home");
  window.location.href = `http://localhost:3000/v1/auth/login/lastfm?return_to=${returnTo}`;
}

function Title() {
  return (
    <div className="font-[Roboto_Mono]">
      <h2 className="leading-none m-0 text-6xl lg:text-8xl text-black font-semibold justify-center text-center">
        Blendify
      </h2>
      <p className=" font-mono px-5 mt-4 text-black text-2xl text-center break-words">
        Blend your music tastes and more
      </p>
    </div>
  );
}
