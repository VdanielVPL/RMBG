import React from 'react'
import { createRoot } from 'react-dom/client'
import './styles/style.css'
import App from './App'
import { MainProvider } from './components/contexts'

// Wails highlights the image container only for files, links are highlighted here
let linkTarget: Element | null = null

function highlightLinkTarget(target: Element | null) {
    if (target == linkTarget) return
    linkTarget?.classList.remove('file-drop-target-active')
    target?.classList.add('file-drop-target-active')
    linkTarget = target
}

// Drags without files (links, text) only drop on the image containers, anywhere else the link would open in a new window
document.addEventListener('dragover', (event) => {
    if (event.dataTransfer && !event.dataTransfer.types.includes('Files')) {
        event.preventDefault()
        const target = (event.target as Element).closest('[data-file-drop-target]')
        event.dataTransfer.dropEffect = target ? 'copy' : 'none'
        highlightLinkTarget(target)
    }
})
document.addEventListener('dragleave', (event) => {
    if (!event.relatedTarget) highlightLinkTarget(null)
})
document.addEventListener('drop', (event) => {
    event.preventDefault()
    highlightLinkTarget(null)
})

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <MainProvider>
            <App/>
        </MainProvider>
    </React.StrictMode>
)
