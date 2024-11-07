import React from 'react'
import {createRoot} from 'react-dom/client'
import './styles/style.css'
import App from './App'
import { MainProvider } from './components/contexts/MainContext'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <MainProvider>
            <App/>
        </MainProvider>
    </React.StrictMode>
)
