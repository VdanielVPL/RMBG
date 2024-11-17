import './styles/App.css';
import { HashRouter, Route, Routes } from 'react-router-dom';
import { Nav } from './components/nav/Nav';
import { MainWrapper } from './components/MainWrapper';
import { RMBGView } from './components/views/RMBGView';
import { CropView } from './components/views/CropView';
import { ImageProvider } from './components/contexts/ImageContext';

function App() {
    return (
        <HashRouter>
            <div id="App">
                <Nav/>
                <ImageProvider>
                    <MainWrapper>
                        <Routes>
                            <Route path="/" element={<RMBGView/>} />
                            <Route path="/crop" element={<CropView/>} />
                        </Routes>
                    </MainWrapper>
                </ImageProvider>
            </div>
        </HashRouter>
    )
}

export default App
