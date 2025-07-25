
import './App.css'

function App() {


  return (
    <div className="">
      <nav className="nav">
          <div className="flex space-x-4">
            <button
              className="bg-white text-blue-700 px-6 py-2 rounded-full shadow-md hover:bg-blue-100 hover:shadow-lg transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-opacity-75"
              onClick={() => navigateTo('/about')}
            >
              About
            </button>
            <button
              className="bg-white text-black-700 px-6 py-2"
              onClick={() => navigateTo('/privacypolicy')}
            >
              Privacy Policy
            </button>
          </div>
      </nav>
    <div className="min-h-screen flex items-center justify-center">
      <div className="bg-white p-8  flex flex-col items-center"> {/* Added flex-col and items-center for vertical centering */}
        <h1 className="text-4xl text-gray-800 mb-6 text-center">
          Blendify
        </h1>

        {/* Button with image inside, perfectly clipped */}
        <div class="shape text log-in-to-0c94be86b263">
          <div class="text-node-html" id="html-text-node-a0c26d1e-7be0-8055-8006-0c94be86b263" data-x="412.6666666666667" data-y="352">
            <div class="root rich-text root-0" xmlns="http://www.w3.org/1999/xhtml">
              <div class="paragraph-set root-0-paragraph-set-0">
                <p class="paragraph root-0-paragraph-set-0-paragraph-0" dir="ltr"><span class="text-node root-0-paragraph-set-0-paragraph-0-text-0">Log in to start merging your music</span></p>
              </div>
            </div>
          </div>
        </div>
        
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

function loginButton(platform) {
  switch (platform) {
    case 'spotify':
      return (
        
      <button
        className=" left-[20px] top-[24px] w-[200px] h-[100px] border-4 inset-0 border-black shadow-[4px_4px_0_0_#000] hover:shadow-[2px_2px_0_0_#000] active:translate-x-[1px] active:translate-y-[1px] transition-all overflow-hidden inset bg-transparent p-0"
      >
        <img
          src=".\src\assets\spotify.png"
          alt="Spotify Logo"
          className="inset-0 w-full h-full object-contain bg-black"
        />
      </button>

      );
   
    case 'apple':
      return (
        
        <button
          className=" left-[20px] top-[24px] w-[200px] h-[100px] border-4 inset-0 border-black shadow-[4px_4px_0_0_#000] hover:shadow-[2px_2px_0_0_#000] active:translate-x-[1px] active:translate-y-[1px] transition-all overflow-hidden inset bg-transparent p-0"
        >
          <img
            src=".\src\assets\apple.png"
            alt="Apple Logo"
            className="inset-0 w-full h-full object-contain bg-black"
          />
        </button>
  
        );
    case 'tidal':
      return (
        
        <button
          className=" left-[20px] top-[24px] w-[200px] h-[100px] border-4 inset-0 border-white shadow-[4px_4px_0_0_#000] hover:shadow-[2px_2px_0_0_#000] active:translate-x-[1px] active:translate-y-[1px] transition-all overflow-hidden inset bg-transparent p-0"
        >
          <img
            src=".\src\assets\tidal.png"
            alt="Tidal Logo"
            className="inset-0 w-full h-full object-contain bg-white"
          />
        </button>
  
        );
    default:
      return <button className="btn btn-default">Login</button>;
  }
}



export default App
