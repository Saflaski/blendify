import '/src/assets/styles/home.css'

import React from "react";

export async function Home() {

    // const res = await fetch("127.0.0.1:3000/api/validate", {
    //     method: 'GET',
    //     credentials: 'include',
    //     headers: { 'Accept': 'application/json' }
    // });

    // if (!res.ok) {
    //     throw redirect("/login")
    // }
   
    return (
        <h1 className='text-gray-700'>Home</h1>
    )
}




