import { Navbar } from "./components/Navbar";
import { Outlet } from "react-router-dom";
import "/src/assets/styles/index.css";
import { useDarkMode } from "/src/hooks/DarkMode";
export function Layout() {
  const [dark, setDark] = useDarkMode();
  return (
    <>
      <Navbar dark={dark} setDark={setDark} />
      <main>
        <div className="route-layer bg-[#f8f3e9] dark:bg-[#1a1917] dark:text-amber-50">
          <Outlet />
        </div>
      </main>
    </>
  );
}
