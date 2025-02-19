import '../../styles/RMBGView.css';
import { useContext, useLayoutEffect, useRef } from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowRight, faCloudArrowDown, faCopy, faCropSimple, faImage, faSpinner, faTrashCan } from '@fortawesome/free-solid-svg-icons';
import { MainContext, ImageContext } from "../contexts";
import { InputImageContainer } from '../InputRMBGImageContainer';
import { Button, SelectList, Selected } from '../inputs';
import { CopyImage, RemoveBackground, SaveImage, SetModel, ClearImageMem, FromRMBGtoCrop } from '../../../wailsjs/go/main/App';
import { EventsEmit } from '../../../wailsjs/runtime/runtime';
import { useNavigate } from 'react-router-dom';

const tiles: Selected[] = [
    {text: "U<sup class='supu2net'>2</sup><span style='letter-spacing: 0.1em'>&nbsp;</span>-Net", value: "u2net"},
    {text: "IS-Net - general", value: "isnet-general-use"},
    {text: "BiRefNet - general", value: "birefnet-general"},
    {text: "BiRefNet - massive", value: "birefnet-massive"},
    {text: "Bria RMBG 2.0", value: "bria-rmbg"}
]

function OptionBar(props: {callback: (selected: Selected) => void}) {

    const { model } = useContext(ImageContext);
    const { strings, isDarkMode, accentColor } = useContext(MainContext);

    return (
        <div style={{display: 'flex', alignItems: 'center', justifyContent: 'space-around'}}>
            <span style={{color: !isDarkMode?'black':undefined}}>{strings["AImodel"]}</span>
            <SelectList tiles={tiles} width='170px' defaultSelected={model} onSelected={props.callback} accentColor={accentColor} isDarkMode={isDarkMode}/>
        </div>
    )
}

function OutputImageContainer() {
    const { outputRMBGImage, removingBG } = useContext(ImageContext);
    const imgRef = useRef<HTMLImageElement>(null);

    useLayoutEffect(() => {  
        if(imgRef.current){
            if (outputRMBGImage != "") {
                imgRef.current.style.width = '100%';
                imgRef.current.style.height = '100%';
            }else{
                imgRef.current.style.width = 'auto';
                imgRef.current.style.height = 'auto';
            }
        }
    }, [outputRMBGImage]);
    
    return (
        <div className="imageContainer" style={{cursor: 'auto'}}>
            <img ref={imgRef} src={outputRMBGImage} style={{userSelect: 'none', pointerEvents: 'none'}} draggable={false}></img>
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
    const { strings, isDarkMode } = useContext(MainContext);
    const { setOutputRMBGImage, setInputRMBGImage, setCropImage, setPNGCropImage, setJPGCropImage, inputRMBGImage, outputRMBGImage, removingBG, setModel, setIsJPG } = useContext(ImageContext);
    let navigate = useNavigate();

    function handleModelChange(selected: Selected) {
        setModel(selected.value);
        SetModel(selected.value);
    }

    function removeBackground() {
        setOutputRMBGImage("");
        if (inputRMBGImage != "") {
            RemoveBackground().then((result) => {
                if (result != null) {
                    const [base64, fileType] = result;
                    if (base64 == "" || fileType == "") {
                        setOutputRMBGImage("");
                    } else {
                        setOutputRMBGImage(`data:${fileType};base64,${base64}`);
                        EventsEmit('removingbg', false);
                    }
                }
            });
        }
    }

    function ClearView() {
        setOutputRMBGImage("");
        setInputRMBGImage("");
        ClearImageMem("RMBG");
    }

    function ToCropView(inOrOut: string) {
        if (inOrOut == "in") {
            if (inputRMBGImage != "") {
                const newImage = inputRMBGImage;
                setCropImage(newImage);
                setPNGCropImage(newImage);
                setJPGCropImage("");
                FromRMBGtoCrop(0);
                navigate('/crop');
            }
        }else if (inOrOut == "out") {
            if (outputRMBGImage != "") {
                setIsJPG(false);
                const newImage = outputRMBGImage;
                setPNGCropImage(newImage);
                setCropImage(newImage);
                setJPGCropImage("");
                FromRMBGtoCrop(1);
                navigate('/crop');
            }
        }
    }

    return (
        <div className='View' style={{display: 'flex', flexDirection: 'column', gap: '10px'}}>
            <span style={{color: isDarkMode?'white':'black'}} className='ViewHeader'>{strings["BGrem"]}</span>
            <OptionBar callback={handleModelChange}></OptionBar>
            <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', height: '100%'}}>
                <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-evenly', alignItems: 'center'}}>
                    <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', gap: '10px'}}>
                        <InputImageContainer></InputImageContainer>
                        <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'start', width: '100%', gap: '10px'}}>
                            <Button style={{width: '40px', height: '40px'}} onClick={ClearView} title={strings["ClearView"]}>
                                <FontAwesomeIcon icon={faTrashCan} size='xl' style={{transform: 'scale(0.8)'}} />
                            </Button>
                            <Button style={{width: '40px', height: '40px'}} onClick={() => ToCropView('in')} title={strings["ToCropView"]}>
                                <FontAwesomeIcon icon={faCropSimple} size='xl' style={{transform: 'scale(0.8)'}}/>
                            </Button>
                        </div>
                    </div>
                    <div style={{display: 'flex', flexDirection: 'column', height: '100%', gap: '10px'}}>
                        <div style={{flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center'}}>
                            <Button style={{borderRadius: '100%', width: '45px', height: '45px', filter: removingBG?'brightness(60%) contrast(100%)':undefined}} onClick={removeBackground} disabled={removingBG} title={strings["RemoveBG"]}>
                                <FontAwesomeIcon icon={faArrowRight} size='2xl'/>
                            </Button>
                        </div>
                        <div style={{height: '40px'}}></div>
                    </div>
                    <div style={{display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', gap: '10px'}}>
                        <OutputImageContainer></OutputImageContainer>
                        <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-between', width: '100%', gap: '10px'}}>
                            <Button style={{width: '40px', height: '40px'}} onClick={() => ToCropView('out')} title={strings["ToCropView"]}>
                                <FontAwesomeIcon icon={faCropSimple} size='xl' style={{transform: 'scale(0.8)'}}/>
                            </Button>
                            <div style={{display: 'flex', flexDirection: 'row', gap: '10px'}}>
                                <Button style={{width: '40px', height: '40px'}} onClick={() => CopyImage('RMBG')} title={strings["CopyImage"]}>
                                    <FontAwesomeIcon icon={faCopy} size='xl'style={{transform: 'scale(0.8)'}} />
                                </Button>
                                <Button style={{width: '40px', height: '40px'}} onClick={() => SaveImage('RMBG')} title={strings["SaveImage"]}>
                                    <FontAwesomeIcon icon={faCloudArrowDown} size='xl' style={{transform: 'scale(0.8)'}}/>
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}