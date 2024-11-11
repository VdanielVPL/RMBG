import '../../styles/RMBGView.css';
import { useContext, useRef, useState } from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowRight } from '@fortawesome/free-solid-svg-icons';
import { MainContext } from "../contexts/MainContext";
import { ImageContainer } from '../ImageContainer';
import SelectList, { Selected } from '../inputs/SelectList';
import { Button } from '../inputs/Button';
import { RemoveBackground, SetModel } from '../../../wailsjs/go/main/App';

const tiles = [
    {text: "u2net", value: "u2net"},
    {text: "isnet - general", value: "isnet-general-use"},
    {text: "birefnet - general", value: "birefnet-general"},
    {text: "birefnet - massive", value: "birefnet-massive"},
]

function OptionBar(props: {callback: (selected: Selected) => void}) {
    return (
        <div style={{display: 'flex', alignItems: 'center', justifyContent: 'space-around'}}>
            <span>AI model:</span>
            <SelectList tiles={tiles} width='170px' onSelected={props.callback}/>
        </div>
    )
}

export function RMBGView() {
    const { strings } = useContext(MainContext);
    const rembgimage = useRef<HTMLImageElement>(null);

    function handleModelChange(selected: Selected) {
        SetModel(selected.value);
    }

    function removeBackground() {
        console.log("Remove background");
        RemoveBackground().then(([base64, fileType]) => {
            if (rembgimage.current) {
                if (base64 == "" || fileType == "") {
                    rembgimage.current.src = "";
                } else {
                    rembgimage.current.src = `data:${fileType};base64,${base64}`;
                }
            }
        });
    }

    return (
        <div style={{display: 'flex', flexDirection: 'column', height: '100%', gap: '10px'}}>
            <span className='ViewHeader'>{strings["BGrem"]}</span>
            <OptionBar callback={handleModelChange}></OptionBar>
            <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', height: 'inherit'}}>
                <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-evenly', alignItems: 'center'}}>
                    <ImageContainer></ImageContainer>
                    <Button style={{borderRadius: '100%', width: '45px', height: '45px'}} onClick={removeBackground}>
                        <FontAwesomeIcon icon={faArrowRight} size='2xl'/>
                    </Button>
                    <div className="imageContainer">
                        <img ref={rembgimage} style={{userSelect: 'none', pointerEvents: 'none'}} draggable={false}></img>
                    </div>
                </div>
            </div>
        </div>
    )
}