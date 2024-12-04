import './styles/App.css';
import { HashRouter, Route, Routes, useLocation } from 'react-router-dom';
import { Nav } from './components/nav/Nav';
import { MainWrapper } from './components/MainWrapper';
import { RMBGView } from './components/views/RMBGView';
import { CropView } from './components/views/CropView';
import { ImageProvider } from './components/contexts/ImageContext';
import { TransitionGroup, CSSTransition } from 'react-transition-group';

function App() {
    return (
        <HashRouter>
            <div id="App">
                <Nav/>
                <ImageProvider>
                    <MainWrapper>
                        <Views/>
                    </MainWrapper>
                </ImageProvider>
            </div>
        </HashRouter>
    )
}

function Views() {
    const location = useLocation();

    return(
        <TransitionGroup>
            <CSSTransition key={location.pathname} classNames="fade" timeout={300} unmountOnExit>
                <Routes location={location}>
                    <Route path="/" element={<RMBGView/>} />
                    <Route path="/crop" element={<CropView/>} />
                </Routes>
            </CSSTransition>
        </TransitionGroup>
    )
}

export default App