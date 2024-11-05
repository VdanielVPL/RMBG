import './styles/App.css';
import { HashRouter, Route, Routes } from 'react-router-dom';
import { Nav } from './components/Nav';
import { MainWrapper } from './components/MainWrapper';

function App() {
    return (
        <HashRouter>
            <div id="App">
                <Nav/>
                <MainWrapper>
                    <Routes>
                        <Route path="/" element={<div>Home</div>} />
                        <Route path="/crop" element={<div>Crop</div>} />
                    </Routes>
                </MainWrapper>
            </div>
        </HashRouter>
    )
}

export default App
