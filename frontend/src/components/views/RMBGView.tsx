import '../../styles/RMBGView.css';
import { RefObject, useContext, useEffect, useRef, useState } from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowRight, faCloudArrowDown, faCopy, faImage, faSpinner } from '@fortawesome/free-solid-svg-icons';
import { MainContext } from "../contexts/MainContext";
import { InputImageContainer } from '../InputImageContainer';
import SelectList, { Selected } from '../inputs/SelectList';
import { Button } from '../inputs/Button';
import { CopyImage, RemoveBackground, SaveImage, SetModel } from '../../../wailsjs/go/main/App';
import { EventsOff, EventsOn } from '../../../wailsjs/runtime/runtime';

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

function OutputImageContainer(props: {rembgimage: RefObject<HTMLImageElement>, removingBG: boolean, trigger?: boolean}) {
    return (
        <div className="imageContainer">
            <img ref={props.rembgimage} style={{userSelect: 'none', pointerEvents: 'none'}} draggable={false}></img>
            <div style={{position: 'absolute', height: '100%', width: '100%', color: 'black', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: -1, backgroundColor: 'white'}}>
                {props.removingBG?
                    <FontAwesomeIcon icon={faSpinner} color='lightgray' spinPulse style={{height: '50%', width: '50%', fontSize: '50%'}} />
                :
                    (!props.rembgimage.current?.src || props.rembgimage.current?.src === "") && <FontAwesomeIcon icon={faImage} color='lightgray' style={{height: '50%', width: '50%', fontSize: '50%'}}/>
                }
            </div>
        </div>
    )
}

export function RMBGView() {
    const { strings } = useContext(MainContext);
    const [removingBG, setRemovingBG] = useState(false);
    const [trigger, setTrigger] = useState(false);
    const rembgimage = useRef<HTMLImageElement>(null);
    

    function handleModelChange(selected: Selected) {
        SetModel(selected.value);
    }

    function removeBackground() {
        RemoveBackground().then(([base64, fileType]) => {
            if (rembgimage.current) {
                if (base64 == "" || fileType == "") {
                    rembgimage.current.src = "";
                } else {
                    rembgimage.current.src = `data:${fileType};base64,${base64}`;
                }
            }
            setTrigger((prev) => !prev);
        });
    }

    useEffect(() => {
        const handleRemBGstatus = (status: boolean) => {
            setRemovingBG(()=>status);
        }

        EventsOn('removingbg', handleRemBGstatus);

        return () => {
            EventsOff('removingbg');
        }
    }, []);

    return (
        <div style={{display: 'flex', flexDirection: 'column', height: '100%', gap: '10px'}}>
            <span className='ViewHeader'>{strings["BGrem"]}</span>
            <OptionBar callback={handleModelChange}></OptionBar>
            <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', height: 'inherit'}}>
                <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-evenly', alignItems: 'center'}}>
                    <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', gap: '10px'}}>
                        <InputImageContainer></InputImageContainer>
                        <div style={{padding: '10px', fontSize: '1rem', fontFamily: 'Arial', visibility: 'hidden'}}>
                            s
                        </div>
                    </div>
                    <Button style={{borderRadius: '100%', width: '45px', height: '45px', backgroundColor: removingBG?'indigo':undefined}} onClick={removeBackground} disabled={removingBG}>
                        <FontAwesomeIcon icon={faArrowRight} size='2xl'/>
                    </Button>
                    <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', gap: '10px'}}>
                        <OutputImageContainer rembgimage={rembgimage} removingBG={removingBG} trigger={trigger}></OutputImageContainer>
                        <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'end', width: '100%', gap: '10px'}}>
                            <Button style={{width: '40px', height: '40px'}} onClick={SaveImage}>
                                <FontAwesomeIcon icon={faCloudArrowDown} size='xl' style={{transform: 'scale(0.8)'}}/>
                            </Button>
                            <Button style={{width: '40px', height: '40px'}} onClick={CopyImage}>
                                <FontAwesomeIcon icon={faCopy} size='xl'style={{transform: 'scale(0.8)'}} />
                            </Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}