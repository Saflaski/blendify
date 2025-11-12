import "/src/assets/styles/App.css";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { Login } from "./routes/Login";
import { Home } from "./routes/home";
import { NotFound } from "./routes/NotFound";
import { About } from "./routes/About";
import { Privacy } from "./routes/Privacy";
import { Layout } from "./Layout";
import { authLoader } from "./utils/authLoader";

const router = createBrowserRouter([
  {
    element: <Layout />,
    children: [
      { path: "/login/", element: <Login /> },
      { path: "/about/", element: <About /> },
      { path: "/privacy/", element: <Privacy /> },
      // { path: "/", element: <Home /> },
      { path: "/", element: <Home />, loader: authLoader },
      // { path: "/home", element: <Home /> },
      { path: "/home", element: <Home />, loader: authLoader },

      { path: "*", element: <NotFound /> },
    ],
  },
]);

function App() {
  return (
    <RouterProvider router={router} />

    // <Router>
    //   <Routes>
    //     <Route element={<Layout/>}>
    //       <Route path="/login/" element={<Login />} />
    //       <Route path="/about/" element={<About />} />
    //       <Route path="/privacy/" element={<Privacy />} />
    //       <Route path="/"

    //         element={<Home />} />
    //       <Route path="*" element={<NotFound />} />
    //     </Route>

    //   </Routes>
    // </Router>
  );
}

export default App;
