import '../../styles/RMBGView.css';
import { useContext } from "react";
import { MainContext } from "../contexts/MainContext";
import { ImageContainer } from '../ImageContainer';

export function RMBGView() {
    const { strings } = useContext(MainContext);

    return (
        <div style={{display: 'flex', flexDirection: 'column'}}>
            <span className='rmbgViewHeader'>{strings["BGrem"]}</span>
            <ImageContainer></ImageContainer>
        </div>
    )
}