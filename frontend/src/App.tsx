import { useState } from 'react'
import './App.css'
import Sidebar from './components/Sidebar/Sidebar'
import Dashboard from './components/Dashboard/Dashboard'


let sidebar_classes={}
function App() {

  return (
    <div>
      <Sidebar classes={sidebar_classes}/>
      <Dashboard />
    </div>
  )
}

export default App
