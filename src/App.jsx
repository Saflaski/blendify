
import './App.css'
import Navbar from './Navbar';

function App() {

  // const handleLastFMClick = useLoginClick("lastFM");

  return (
    <div className=''>
      {/* <Navbar /> */}
    <div className="responsive-position">
      
      <div className=" center-div p-2  w-full"> {/* Added flex-col and items-center for vertical centering */}
        

       
        <Title />
        {/* New section for login buttons */}
        <div className="flex justify-center flex-wrap gap-4"> {/* Using flex-wrap and gap for responsiveness and spacing */}
          {loginButton2('spotify')}
          {loginButton2('apple')}
          {loginButton2('tidal')}
        </div>
      </div>
      </div>
      </div>
  )
}


function navigateTo(path) {
  // This function would handle navigation logic, e.g., using React Router
  console.log(`Navigating to ${path}`);
}


function loginButton2(platform) {
  const platforms = {
     spotify: {
      src: './src/assets/lastfm.svg',
      alt: 'LastFM Logo',
      wrapperClass: 'mx-5 my-4',
    },
    apple: {
      src: './src/assets/apple2.svg',
      alt: 'Apple Logo',
      wrapperClass: 'mx-1',
    },
    tidal: {
      src: './src/assets/tidal.svg',
      alt: 'Tidal Logo',
      wrapperClass: 'mx-1',
    },
  }

  const config = platforms[platform];

  return (
      <button className="login-button-box">
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

function Title() {
  return (
    <div className="w-full max-w-screen-sm px-4 mx-auto flex flex-col items-center mb-8  text-center">
      <h1 className=" text-gray-800 mb-2 break-words">
        Blendify
      </h1>

      <p className="subheading break-words">
        Blend your music the way you know it. But cross-platform and better.
      </p>
    </div>
  );
}

function loginButton(platform) {
  switch (platform) {
    case 'spotify':
      return (

        <button
          className="login-button-box "
        >
          <div className="mx-5 my-4">
          <img
            src=".\src\assets\lastfm.svg"
            alt="LastFM Logo"
            className=" object-center my-40 image-render-pixel bg-white "
          />
        </div>
        </button>

      );

    case 'apple':
      return (

        <button
          className="login-button-box">
          <div className="mx-1">
          <img
            src=".\src\assets\apple2.svg"
            alt="LastFM Logo"
            className=" object-center my-40 image-render-pixel bg-white "
          />
        </div>

        </button>

      );
    case 'tidal':
      return (

        <button
          className="login-button-box">
          <div className="mx-1">
          <img
            src=".\src\assets\tidal.svg"
            alt="Tidal Logo"
            className=" object-center image-render-pixel bg-white "
          /></div>
        </button>

      );
    default:
      return <button className="btn btn-default">Login</button>;
  }
}



export default App
