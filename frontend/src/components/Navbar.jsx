import {Link} from 'react-router-dom'
import '/src/assets/styles/navbar.css'
import Cookies from 'js-cookie';
import { useState, useEffect, useRef } from 'react';

export function Navbar() {

    const [open, setOpen] = useState(false)
    const dropdownRef = useRef(null);
    const dropButtonRef = useRef(null);

    
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (
                !dropdownRef.current?.contains(event.target) && !dropButtonRef.current?.contains(event.target)
            ) {
                setOpen(false);
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside)
        };
    }, []
    );

    

    return (
        <div>
        <div className='grid-container'>
            <div className='menuOver'>
                <div className='menu-trigger' ref={dropButtonRef}>
                    <MenuButton  onClick={() => {setOpen(!open)}} />
                </div>   
            </div>
            <div className=" button-blendify">
                <Link to="/home">
                        <button >Blendify</button>
                    </Link> 
            </div>
        </div>
            <div className={`dropdown-menu ${open ? 'active' : 'inactive'}`}
                ref={dropdownRef}
            >
                <ul>
                    <DropDownItem page={"/"} funcName={null} text={"Home"} />
                    <DropDownItem page={"/about"} funcName={null} text={"About"} />
                    <DropDownItem page={"/privacy"} funcName={null} text={"Privacy"} />
                    <DropDownItem page={"/login"} funcName={handleLogOut} text={"Log Out"} />
                    </ul>
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

function DropDownItem({ page, funcName, text }) {
    return (
        <li className='dropDownItem'>
            <Link to={page || "/asd"}>
                <button
                    onClick={funcName || null}
                    className={text == "Log Out" ? 'text-red-500' : 'text-black'}
                >{text || "Lorem Ipsum"}</button>
            </Link>
        </li>
    );
}

function MenuButton({ onClick }) {
    return (
    <div>
            <button>
                <img src="/src/assets/images/MenuButton.svg"
                    onClick={onClick}
                />
        </button>
    </div>)
}

