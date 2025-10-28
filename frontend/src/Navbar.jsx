import React, { useState } from 'react';
import './navbar.css'; // Import the CSS file

const Navbar = () => {
    const [isOpen, setIsOpen] = useState(true); // State to manage mobile menu visibility

    const toggleMenu = () => {
        setIsOpen(!isOpen);
    };

    return (
        <nav className="top-0 left-0">
        <div className="">
          <img src="./src/assets/logo.png" alt="Logo" className="logo" />
        </div>
        <div className="">
          <button className="nav-link">Home</button>
          <button  className="nav-link">About</button>
          <button className="nav-link">Contact</button>
        </div>
        
      </nav>
    );
};

export default Navbar;