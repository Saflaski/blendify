import {Link} from 'react-router-dom'
import '/src/assets/styles/navbar.css'
import Cookies from 'js-cookie';

export function Navbar() {
    return (
        
        <div className='grid-container'>
            <div className='justify-self-start mx-4'>
                
            </div>
            <div className=" button-blendify">
                <Link to="/home">
                        <button >Blendify</button>
                    </Link>
            </div>
                
            

            <div className="ml-auto flex space-x-4 px-3 nav-right">
                <Link to="/about">
                    <button className=''>About</button>
                </Link>
                <Link to="/privacy">
                    <button className=''>Privacy Policy</button>
                </Link>
                <Link to="/login">
                    <button onClick={ handleLogOut } className=''>Log Out</button>
                </Link>
            </div>
        </div>



    )
}

async function handleLogOut() {
  await fetch('http://127.0.0.1:3000/api/logout/', {
    method: 'POST',
    credentials: 'include', 
  });

  window.location.href = '/login';
}


