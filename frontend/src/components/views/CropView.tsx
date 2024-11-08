import { useContext } from "react";
import { ImageContainer } from "../ImageContainer";
import { MainContext } from "../contexts/MainContext";

export function CropView() {
    const { strings } = useContext(MainContext);

    return (
        <div style={{display: 'flex', flexDirection: 'column'}}>
            <span className='ViewHeader'>{strings["Crop"]}</span>
            <ImageContainer></ImageContainer>
        </div>
    )
}