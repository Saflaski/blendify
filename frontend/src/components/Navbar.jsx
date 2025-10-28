import {Link} from 'react-router-dom'
import '/src/assets/styles/navbar.css'

export function Navbar() {
    return (
        
        <div className='grid-container'>
            <div className='justify-self-start mx-4'>
                <Link to="/">
                        <button className=" button-blendify">Blendify</button>
                    </Link>
            </div>
                <div className=' flex space-x-4 nav-right'>
                    <Link to="/">
                        <button className="mx-2">Home</button>
                    </Link>
                </div>
            

            <div className="ml-auto flex space-x-4 px-3 nav-right">
                <Link to="/about">
                    <button className='mx-2'>About</button>
                </Link>
                <Link to="/privacy">
                    <button className='mx-2'>Privacy Policy</button>
                </Link>
            </div>
        </div>












        // <div className="w-full flex justify-center items-center px-2 navbar">
        // <Link to="/">
        //     <button className="mx-2">Home</button>
        // </Link>

        // <div className="ml-auto flex space-x-4">
        //     <Link to="/about">
        //     <button>About</button>
        //     </Link>
        //     <Link to="/privacy">
        //     <button>Privacy Policy</button>
        //     </Link>
        // </div>
        // </div>

    )
}