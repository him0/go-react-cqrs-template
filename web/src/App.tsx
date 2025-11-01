import { useState } from 'react'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="App">
      <h1>User Management App</h1>
      <p>
        This is a sample application using Go (DDD), Vite, React, OpenAPI, and Orval.
      </p>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
      </div>
      <p className="info">
        API code will be generated using Orval from OpenAPI spec.
      </p>
    </div>
  )
}

export default App
