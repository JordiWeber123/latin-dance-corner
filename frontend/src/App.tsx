import { Route, Routes } from 'react-router'
import Hello from './pages/Hello.tsx'
import Home from './pages/Home.tsx'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/hello" element={<Hello />} />
    </Routes>
  )
}

export default App
