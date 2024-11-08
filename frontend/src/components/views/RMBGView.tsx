import '../../styles/RMBGView.css';
import { useContext } from "react";
import { MainContext } from "../contexts/MainContext";
import { ImageContainer } from '../ImageContainer';
import SelectList from '../inputs/SelectList';


const tiles = [
    {text: "u2net", value: "u2net"},
    {text: "isnet - general", value: "isnet-general-use"},
    {text: "birefnet - general", value: "birefnet-general"},
    {text: "birefnet - massive", value: "birefnet-massive"},
]

function OptionBar() {
    return (
        <div style={{display: 'flex', alignItems: 'center', justifyContent: 'space-around'}}>
            <span>AI model:</span>
            <SelectList tiles={tiles} width='170px'/>
        </div>
    )
}

export function RMBGView() {
    const { strings } = useContext(MainContext);

    return (
        <div style={{display: 'flex', flexDirection: 'column'}}>
            <span className='ViewHeader'>{strings["BGrem"]}</span>
            <OptionBar></OptionBar>
            <ImageContainer></ImageContainer>
        </div>
    )
}