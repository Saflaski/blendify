
import './App.css'
import Navbar from './Navbar';

function App() {


  return (
    <div className=''>
      {/* <Navbar /> */}
    <div className="responsive-position">
      
      <div className=" center-div p-2  w-full"> {/* Added flex-col and items-center for vertical centering */}
        

       
        <Title />
        {/* New section for login buttons */}
        <div className="flex justify-center flex-wrap gap-4"> {/* Using flex-wrap and gap for responsiveness and spacing */}
          {loginButton('spotify')}
          {loginButton('apple')}
          {loginButton('tidal')}
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

function Title() {
  return (
    <div className="w-full max-w-screen-sm px-4 mx-auto flex flex-col items-center mb-8  text-center">
      <h1 className=" text-gray-800 mb-2 break-words">
        Blendify
      </h1>

      <p className="subheading break-words">
        Log in to start merging your music
      </p>
    </div>
  );
}

function loginButton(platform) {
  switch (platform) {
    case 'spotify':
      return (

        <button
          className="login-button-box"
        >
          <img
            src=".\src\assets\spotify.svg"
            alt="Spotify Logo"
            className=" object-center w-9/10 h-9/10 image-render-pixel bg-white "
          />
        </button>

      );

    case 'apple':
      return (

        <button
          className="login-button-box">
          <img
            src=".\src\assets\apple.svg"
            alt="Apple Logo"
            className=" object-center w-9/10 h-9/10 image-render-pixel bg-white "
          />
        </button>

      );
    case 'tidal':
      return (

        <button
          className="login-button-box">
          <img
            src=".\src\assets\tidal.svg"
            alt="Tidal Logo"
            className=" object-center w-9/10 h-9/10 image-render-pixel bg-white "
          />
        </button>

      );
    default:
      return <button className="btn btn-default">Login</button>;
  }
}



export default App
