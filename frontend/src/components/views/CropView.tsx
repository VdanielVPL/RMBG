import { useContext } from "react";
import { InputImageContainer } from "../CropImageContainer";
import { MainContext } from "../contexts/MainContext";

export function CropView() {
    const { strings } = useContext(MainContext);

    return (
        <div style={{display: 'flex', flexDirection: 'column', height: '100%'}}>
            <span className='ViewHeader'>{strings["Crop"]}</span>
            <div style={{height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
                <InputImageContainer></InputImageContainer>
            </div>
        </div>
    )
}