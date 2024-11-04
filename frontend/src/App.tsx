import { useContext } from 'react';
import { MainContext } from './contexts/MainContext';
import './App.css';
import { HashRouter } from 'react-router-dom';

function App() {
    const { strings } = useContext(MainContext);

    return (
        <HashRouter>
            <div id="App">
                <div id="result" className="result">{strings["Hello"]}</div>
            </div>
        </HashRouter>
    )
}

export default App
