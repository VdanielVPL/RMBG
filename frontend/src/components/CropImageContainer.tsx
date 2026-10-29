import { useEffect, useRef, CSSProperties, useContext, useState, SyntheticEvent, useLayoutEffect, DragEvent } from "react";
import { Events } from "@wailsio/runtime";
import { HandleDrop, OpenImage } from "../../bindings/rmbg/app";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCloudArrowUp, faSpinner } from "@fortawesome/free-solid-svg-icons";
import { ImageContext, MainContext } from "./contexts";
import CropEditor from "./CropEditor";

export function InputImageContainer() {
    const { accentColor } = useContext(MainContext);
    const { cropImage, cropping, setCropImageByType } = useContext(ImageContext);
    const [rect, setRect] = useState<DOMRect | null>(null);
    const [isImageDark, setIsImageDark] = useState<boolean>(true);
    const imageContainer = useRef<HTMLDivElement>(null);
    const imageRef = useRef<HTMLImageElement>(null);
    const alphaPatternRef = useRef<HTMLDivElement>(null);

    function setLoadedImage(result: string[]) {
        if (result.length > 0) {
            const [url, fileType] = result;
            setCropImageByType(url, fileType);
        }
    }

    useEffect(() => {
        return Events.On('filedrop', (event) => {
            if (event.data?.[0] == "CROP") {
                HandleDrop("CROP", event.data[1], false).then(setLoadedImage);
            }
        });
    }, []);

    useLayoutEffect(() => {
        if (cropImage == "") {
            if(imageContainer.current){
                imageContainer.current.style.width = '450px';
                imageContainer.current.style.height = '450px';
                imageContainer.current.style.borderRadius = '10px';
                if(imageRef.current){
                    imageRef.current.style.width = 'auto';
                    imageRef.current.style.height = 'auto';
                }
            }
        }
    }, [cropImage]);

    const calculateBrightness = (image: HTMLImageElement) => {
        const canvas = document.createElement("canvas");
        const ctx = canvas.getContext("2d");
    
        canvas.width = image.width;
        canvas.height = image.height;
        if (!ctx) {
            return false;
        }
        ctx.drawImage(image, 0, 0, image.width, image.height);
        const imageData = ctx.getImageData(0, 0, image.width, image.height);
        const data = imageData.data;
    
        let totalBrightness = 0;
    
        for (let i = 0; i < data.length; i += 4) {
            const r = data[i];
            const g = data[i + 1];
            const b = data[i + 2];
            totalBrightness += (r + g + b) / 3;
        }
    
        const averageBrightness = totalBrightness / (image.width * image.height);
        return averageBrightness <= 128;
    };

    const handleURLDrop = (event: DragEvent<HTMLDivElement>) => {
        const url = event.dataTransfer.getData('URL') || event.dataTransfer.getData('text/plain');
        if (url.startsWith('http://') || url.startsWith('https://')) {
            HandleDrop("CROP", url, true).then(setLoadedImage);
        }
    };

    async function openDialog() {
        if (cropImage != "") {
            return;
        }
        setLoadedImage(await OpenImage("CROP"))
    }

    function onImageLoad(e: SyntheticEvent<HTMLImageElement>) {
        if(imageContainer.current != null){
            if (e.currentTarget.width < e.currentTarget.height) {
                imageContainer.current.style.width = "fit-content";
                imageContainer.current.style.height = '450px';
                e.currentTarget.style.height = "100%";
            } else {
                imageContainer.current.style.height = "fit-content";
                imageContainer.current.style.width = "450px";
                e.currentTarget.style.width = "100%";
            }
            imageContainer.current.style.borderRadius = '0px';
            setIsImageDark(calculateBrightness(e.currentTarget));
            setRect(imageContainer.current.getBoundingClientRect());
            if (alphaPatternRef.current) {
                alphaPatternRef.current.style.opacity = '1';
            }
        };
    }

    const patternStyle = {
        backgroundImage: `url("data:image/svg+xml;utf8,<svg width='20' height='20' xmlns='http://www.w3.org/2000/svg'><rect width='10' height='10' fill='%23808080'/><rect x='10' y='10' width='10' height='10' fill='%23808080'/></svg>")`,
        backgroundSize: '20px 20px',
        position: 'absolute',
        width: '100%',
        height: '100%',
        zIndex: 2,
        filter: 'brightness(0.5)',
        outlineColor: 'gray',
        outlineStyle: 'solid',
        outlineWidth: '1px',
        outlineOffset: '-1px',
        opacity: 0,
        transition: 'opacity 0.5s ease',
    } as CSSProperties;

    return (
        <div ref={imageContainer} id='CROP' className='imageContainer imageContainerCrop' onDrop={handleURLDrop} data-file-drop-target style={{backgroundColor: cropImage!="" && "rgba(255,255,255,0)", cursor: cropImage!="" && 'auto'} as CSSProperties} onClick={openDialog}>
            {cropImage!="" && 
                <img ref={imageRef} src={cropImage} style={{userSelect: 'none', pointerEvents: 'none', filter: 'brightness(0.5)'}} draggable={false} onLoad={onImageLoad}></img>
            }
            <div style={{position: 'absolute', height: '100%', width: '100%', color: 'black', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: -1, opacity: cropImage!=""?0:1}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%', maxHeight: '200px', maxWidth: '200px'}} />
            </div>
            {cropping && 
                <div style={{position: 'absolute', height: '100%', width: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 3}}>
                    <FontAwesomeIcon icon={faSpinner} color='lightgray' spinPulse style={{height: '50%', width: '50%', fontSize: '50%', maxWidth: '160px', maxHeight: '160px'}} />
                </div>
            }
            <div className="dropView" style={{position: 'absolute', zIndex: 4, backgroundColor: accentColor, width: '100%', height: '100%', opacity: '0', transition: 'opacity 0.4s ease', display: 'flex', justifyContent: 'center', alignItems: 'center', filter: cropImage&& 'brightness(0.5)'}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%', maxHeight: '200px', maxWidth: '200px'}} />
            </div>
            {(cropImage && !cropping) && <CropEditor rect={rect} imageRef={imageRef} isImageDark={isImageDark}/>}
            {cropImage!="" && <div ref={alphaPatternRef} style={patternStyle}></div>}
        </div>
    )
}