import { redirect } from 'react-router-dom';

export async function authLoader() {

    console.log("Auth Loader loaded")
    
    const res = await fetch('http://127.0.0.1:3000/api/validate/', {
        credentials: 'include',
        headers: { Accept: 'application/json' },
    });
    console.log(res)
    if (res.ok) {
        return null;
    } else {
        throw redirect('/login/');
    }
}
