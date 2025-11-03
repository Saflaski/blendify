import { Navbar } from "./components/Navbar"
import { Outlet } from "react-router-dom"
import '/src/assets/styles/index.css'
export function Layout() {
    return (
        <>
            <Navbar/>
            <main>
                <div className="route-layer">
                <Outlet />
                </div>
            </main>
        </>
    )

}