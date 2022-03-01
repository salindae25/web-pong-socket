import { useState } from 'react'
import logo from './logo.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="bg-slate-50 w-screen min-h-screen flex justify-center">
      <header className="text-lg font-normal">
        <p>Web pong !</p>
       
     </header>
    </div>
  )
}

export default App
