import { useState } from 'react';
import logo from './assets/images/logo-universal.png';
import { GetLangStrings } from '../wailsjs/go/main/App';
import './App.css';

function App() {
    const [text, setText] = useState('');
    GetLangStrings().then((result) => setText(()=>result["Test"]));

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{text}</div>
        </div>
    )
}

export default App
