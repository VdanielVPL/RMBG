import { useContext } from 'react';
import logo from './assets/images/logo-universal.png';
import { LangContext } from './contexts/LangContext';
import './App.css';

function App() {
    const strings = useContext(LangContext);

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{strings["Hello"]}</div>
        </div>
    )
}

export default App
