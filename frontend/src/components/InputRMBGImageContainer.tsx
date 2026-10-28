import { useEffect, useRef, useContext, useLayoutEffect } from "react";
import { Events } from "@wailsio/runtime";
import { HandleDrop, OpenImage } from "../../bindings/rmbg/app";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCloudArrowUp } from "@fortawesome/free-solid-svg-icons";
import { MainContext, ImageContext } from "./contexts";

export function InputImageContainer() {
    const { accentColor } = useContext(MainContext);
    const { inputRMBGImage, setInputRMBGImage, setOutputRMBGImage } = useContext(ImageContext);
    const imageContainer = useRef<HTMLDivElement>(null);
    const imgRef = useRef<HTMLImageElement>(null);

    function setLoadedImage(result: string[]) {
        if (result.length > 0) {
            const [base64, fileType] = result;
            setInputRMBGImage(`data:${fileType};base64,${base64}`);
            setOutputRMBGImage("");
        }
    }

    useEffect(() => {
        return Events.On('filedrop', (event) => {
            if (event.data?.[0] == "RMBG") {
                HandleDrop("RMBG", event.data[1], false).then(setLoadedImage);
            }
        });
    }, []);

    const handleURLDrop = (event: React.DragEvent<HTMLDivElement>) => {
        event.preventDefault();
        const items = event.dataTransfer?.items;
        if (items) {
            const item = items[0];
            if (item.kind === 'string') {
                item.getAsString((data) => {
                    if (data.startsWith('http://') || data.startsWith('https://')) {
                        // console.log('Dropped URL:', data);
                        HandleDrop("RMBG", data, true).then(setLoadedImage);
                        imageContainer.current?.classList.remove('file-drop-target-active');
                    } else if (data.startsWith('data:')) {
                        imageContainer.current?.classList.remove('file-drop-target-active');
                    }
                });
            }
        }
    };

    async function openDialog() {
        setLoadedImage(await OpenImage("RMBG"))
    }

    useLayoutEffect(() => {  
        if(imgRef.current){
            if (inputRMBGImage != "") {
                imgRef.current.style.width = '100%';
                imgRef.current.style.height = '100%';
            }else{
                imgRef.current.style.width = 'auto';
                imgRef.current.style.height = 'auto';
            }
        }
    }, [inputRMBGImage]);

    return (
        <div ref={imageContainer} id='RMBG' className='imageContainer' onDrop={handleURLDrop} data-file-drop-target onClick={openDialog}>
            {inputRMBGImage!="" && <img ref={imgRef} src={inputRMBGImage} style={{userSelect: 'none', pointerEvents: 'none'}} draggable={false}></img>}
            <div style={{position: 'absolute', height: '100%', width: '100%', color: 'black', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: -1,}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%'}} />
            </div>
            <div className="dropView" style={{position: 'absolute', zIndex: 5, backgroundColor: accentColor, width: '100%', height: '100%', opacity: '0', transition: 'opacity 0.4s ease', display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%'}} />
            </div>
        </div>
    )
}