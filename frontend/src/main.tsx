import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'
import { LangProvider } from './contexts/LangContext'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <LangProvider>
            <App/>
        </LangProvider>
    </React.StrictMode>
)
