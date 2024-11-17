import '../../styles/RMBGView.css';
import { useContext } from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowRight, faCloudArrowDown, faCopy, faImage, faSpinner, faTrashCan } from '@fortawesome/free-solid-svg-icons';
import { MainContext } from "../contexts/MainContext";
import { InputImageContainer } from '../InputRMBGImageContainer';
import SelectList, { Selected } from '../inputs/SelectList';
import { Button } from '../inputs/Button';
import { CopyImage, RemoveBackground, SaveImage, SetModel, ClearImageMem } from '../../../wailsjs/go/main/App';
import { EventsEmit } from '../../../wailsjs/runtime/runtime';
import { ImageContext } from '../contexts/ImageContext';

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

function OutputImageContainer() {
    const { outputRMBGImage, removingBG } = useContext(ImageContext);
    
    return (
        <div className="imageContainer" style={{cursor: 'auto'}}>
            <img src={outputRMBGImage} style={{userSelect: 'none', pointerEvents: 'none'}} draggable={false}></img>
            <div style={{position: 'absolute', height: '100%', width: '100%', color: 'black', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: -2, backgroundColor: 'white'}}>
                {removingBG?
                    <FontAwesomeIcon icon={faSpinner} color='lightgray' spinPulse style={{height: '50%', width: '50%', fontSize: '50%'}} />
                :
                    (!outputRMBGImage || outputRMBGImage === "") && <FontAwesomeIcon icon={faImage} color='lightgray' style={{height: '50%', width: '50%', fontSize: '50%'}}/>
                }
            </div>
        </div>
    )
}

export function RMBGView() {
    const { strings } = useContext(MainContext);
    const { setOutputRMBGImage, setInputRMBGImage, inputRMBGImage, removingBG } = useContext(ImageContext);

    function handleModelChange(selected: Selected) {
        SetModel(selected.value);
    }

    function removeBackground() {
        setOutputRMBGImage("");
        if (inputRMBGImage != "") {
            RemoveBackground().then(([base64, fileType]) => {
                if (base64 == "" || fileType == "") {
                    setOutputRMBGImage("");
                } else {
                    setOutputRMBGImage(`data:${fileType};base64,${base64}`);
                    EventsEmit('removingbg', false);
                }
            });
        }
    }

    function ClearView() {
        setOutputRMBGImage("");
        setInputRMBGImage("");
        ClearImageMem("RMBG");
    }

    return (
        <div style={{display: 'flex', flexDirection: 'column', height: '100%', gap: '10px'}}>
            <span className='ViewHeader'>{strings["BGrem"]}</span>
            <OptionBar callback={handleModelChange}></OptionBar>
            <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', height: 'inherit'}}>
                <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-evenly', alignItems: 'center'}}>
                    <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', gap: '10px'}}>
                        <InputImageContainer></InputImageContainer>
                        <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'start', width: '100%', gap: '10px'}}>
                            <Button style={{width: '40px', height: '40px'}} onClick={ClearView} title='Wyczyść widok'>
                                <FontAwesomeIcon icon={faTrashCan} size='xl' style={{transform: 'scale(0.8)'}} />
                            </Button>
                        </div>
                    </div>
                    <div style={{display: 'flex', flexDirection: 'column', height: '100%', gap: '10px'}}>
                        <div style={{flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center'}}>
                            <Button style={{borderRadius: '100%', width: '45px', height: '45px', backgroundColor: removingBG?'indigo':undefined}} onClick={removeBackground} disabled={removingBG} title='Usuń tło'>
                                <FontAwesomeIcon icon={faArrowRight} size='2xl'/>
                            </Button>
                        </div>
                        <div style={{height: '40px'}}></div>
                    </div>
                    <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', gap: '10px'}}>
                        <OutputImageContainer></OutputImageContainer>
                        <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'end', width: '100%', gap: '10px'}}>
                            <Button style={{width: '40px', height: '40px'}} onClick={SaveImage} title='Zapisz zdjęcie'>
                                <FontAwesomeIcon icon={faCloudArrowDown} size='xl' style={{transform: 'scale(0.8)'}}/>
                            </Button>
                            <Button style={{width: '40px', height: '40px'}} onClick={CopyImage} title='Kopiuj do schowka'>
                                <FontAwesomeIcon icon={faCopy} size='xl'style={{transform: 'scale(0.8)'}} />
                            </Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}