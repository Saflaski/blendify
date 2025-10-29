import '/src/assets/styles/Login.css'

export function Login() {


    return (
    <div className=''>
      {/* <Navbar /> */}
    <div className="responsive-position">
      
      <div className=" center-div p-2  w-full"> {/* Added flex-col and items-center for vertical centering */}
        

       
        <Title />
        {/* New section for login buttons */}
        <div className="mt-1 flex justify-center flex-wrap gap-4"> {/* Using flex-wrap and gap for responsiveness and spacing */}
          {loginButton('lastfm')}
          {loginButton('apple')}
          {loginButton('tidal')}
        </div>
      </div>
      </div>
      </div>
  )
}




function loginButton(platform) {
  const platforms = {
     lastfm: {
      src: '/src//assets/images/lastfm.svg',
      alt: 'LastFM Logo',
      wrapperClass: 'mx-8',
      handleFunc: handleLastFMClick
    },
    apple: {
      src: '/src//assets/images/apple2.svg',
      alt: 'Apple Logo',
      wrapperClass: 'mx-1',
    },
    tidal: {
      src: '/src//assets/images/tidal.svg',
      alt: 'Tidal Logo',
      wrapperClass: '',
    },
  }

  const config = platforms[platform];

  return (
      <button className="login-button-box" onClick={config.handleFunc}>
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
  const returnTo = encodeURIComponent('http://127.0.0.1:5173/home');
  window.location.href = `http://127.0.0.1:3000/oauth/lastfm/login?return_to=${returnTo}`;
}

function Title() {
  return (
    <div className="">
      <h1 className= "">
        Blendify
      </h1>
      <p className='subheading text-gray-500 text-center break-words'>
        Blend your music tastes and more
      </p>

      
    </div>
    
  );
}

