import './styles/App.css';
import { HashRouter, Route, Routes, useLocation } from 'react-router-dom';
import { Nav } from './components/nav/Nav';
import { MainWrapper } from './components/MainWrapper';
import { RMBGView, CropView } from './components/views';
import { ImageProvider } from './components/contexts';
import { AnimatePresence, motion } from 'motion/react';

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
        <AnimatePresence mode="wait">
            <motion.div
                key={location.pathname}
                initial={{ opacity: 0}}
                animate={{ opacity: 1}}
                exit={{ opacity: 0}}
                transition={{ duration: 0.19 }}
            >
                <Routes location={location}>
                    <Route path="/" element={<RMBGView/>} />
                    <Route path="/crop" element={<CropView/>} />
                </Routes>
            </motion.div>
        </AnimatePresence>
    )
}

export default App