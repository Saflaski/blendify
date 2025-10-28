
import '/src/assets/styles/App.css'
import {BrowserRouter as Router, Routes, Route} from 'react-router-dom'
import { Login } from './routes/Login'
import { Home } from './routes/home'
import { NotFound } from './routes/NotFound'
import { About } from './routes/About'
import {Privacy} from './routes/Privacy'
import {Layout} from "./Layout"
function App() {

  return (

    
    <Router>
      <Routes>
        <Route element={<Layout/>}>
          <Route path="/login/" element={<Login />} />
          <Route path="/about/" element={<About />} />
          <Route path="/privacy/" element={<Privacy />} />
          <Route path="/" element={<Home />} />
          <Route path="*" element={<NotFound />} />
        </Route>
        
        
        
      </Routes>
    </Router>
  )
}


export default App
