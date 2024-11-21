import { useEffect, useRef, CSSProperties, useContext, useState, ReactEventHandler, SyntheticEvent } from "react";
import { OnFileDrop, OnFileDropOff } from "../../wailsjs/runtime/runtime";
import { HandleDrop, OpenImage } from "../../wailsjs/go/main/App";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCloudArrowUp } from "@fortawesome/free-solid-svg-icons";
import { ImageContext } from "./contexts/ImageContext";
import CropEditor from "./CropEditor";

export function InputImageContainer() {
    // const [filePath, setFilePath] = useState<string>('');
    const { cropImage, setCropImage } = useContext(ImageContext);
    const imageContainer = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleFileDrop = (x: number, y: number, paths: string[]) => {
            // console.log('File dropped:', paths);
            // setFilePath(paths[0]);
            HandleDrop("CROP", paths[0], false).then(([base64, fileType]) => {
                if (base64 == "" || fileType == "") {
                    setCropImage("");
                }else{
                    setCropImage(`data:${fileType};base64,${base64}`);
                }
            });
        };

        OnFileDrop(handleFileDrop, true);

        return () => {
            OnFileDropOff();
        };
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
                        HandleDrop("CROP", data, true).then(([base64, fileType]) => {
                            if (base64 == "" || fileType == "") {
                                setCropImage("");
                            }else{
                                setCropImage(`data:${fileType};base64,${base64}`);
                            }
                        });
                        imageContainer.current?.classList.remove('wails-drop-target-active');
                    }else if (data.startsWith('data:')){
                        imageContainer.current?.classList.remove('wails-drop-target-active');
                    }
                });
            }
        }
    };

    async function openDialog() {
        const result = await OpenImage("CROP")
        if (result != null) {
            const [base64, fileType] = result;
            if (base64 == "" || fileType == "") {
                setCropImage("");
            }else{
                setCropImage(`data:${fileType};base64,${base64}`);
            }
        }
    }

    function onImageLoad(e: SyntheticEvent<HTMLImageElement>) {
        if(imageContainer.current != null){
            console.log(e.currentTarget.width, e.currentTarget.height);
            if (e.currentTarget.width < e.currentTarget.height) {
                imageContainer.current.style.width = "fit-content";
                imageContainer.current.style.height = "100%";
                e.currentTarget.style.height = "100%";
            } else {
                imageContainer.current.style.height = "fit-content";
                imageContainer.current.style.width = "100%";
                e.currentTarget.style.width = "100%";
            }
            imageContainer.current.style.borderRadius = '0px';
        };
    }

    return (
        <div ref={imageContainer} className='imageContainer imageContainerCrop' onDrop={handleURLDrop} style={{'--wails-drop-target': 'drop', backgroundColor: cropImage!="" && "rgba(255,255,255,0)"} as CSSProperties} onClick={openDialog}>
            <img src={cropImage} style={{userSelect: 'none', pointerEvents: 'none'}} draggable={false} onLoad={onImageLoad}></img>
            <div style={{position: 'absolute', height: '100%', width: '100%', color: 'black', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: -1,}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%', maxHeight: '200px', maxWidth: '200px'}} />
            </div>
            {cropImage && <CropEditor />}
        </div>
    )
}