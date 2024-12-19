import { useEffect, useRef, CSSProperties, useContext, useLayoutEffect } from "react";
import { OnFileDrop, OnFileDropOff } from "../../wailsjs/runtime/runtime";
import { HandleDrop, OpenImage } from "../../wailsjs/go/main/App";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCloudArrowUp } from "@fortawesome/free-solid-svg-icons";
import { ImageContext } from "./contexts/ImageContext";
import { MainContext } from "./contexts/MainContext";

export function InputImageContainer() {
    // const [filePath, setFilePath] = useState<string>('');
    const { accentColor } = useContext(MainContext);
    const { inputRMBGImage, setInputRMBGImage } = useContext(ImageContext);
    const imageContainer = useRef<HTMLDivElement>(null);
    const imgRef = useRef<HTMLImageElement>(null);

    useEffect(() => {
        const handleFileDrop = (x: number, y: number, paths: string[]) => {
            // console.log('File dropped:', paths);
            // setFilePath(paths[0]);
            HandleDrop("RMBG", paths[0], false).then(([base64, fileType]) => {
                if (base64 == "" || fileType == "") {
                    setInputRMBGImage("");
                }else{
                    setInputRMBGImage(`data:${fileType};base64,${base64}`);
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
                        HandleDrop("RMBG", data, true).then(([base64, fileType]) => {
                            if (base64 == "" || fileType == "") {
                                setInputRMBGImage("");
                            }else{
                                setInputRMBGImage(`data:${fileType};base64,${base64}`);
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
        const result = await OpenImage("RMBG")
        if (result != null) {
            const [base64, fileType] = result;
            if (base64 == "" || fileType == "") {
                setInputRMBGImage("");
            }else{
                setInputRMBGImage(`data:${fileType};base64,${base64}`);
            }
        }
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
        <div ref={imageContainer} className='imageContainer' onDrop={handleURLDrop} style={{'--wails-drop-target': 'drop'} as CSSProperties} onClick={openDialog}>
            <img ref={imgRef} src={inputRMBGImage} style={{userSelect: 'none', pointerEvents: 'none'}} draggable={false}></img>
            <div style={{position: 'absolute', height: '100%', width: '100%', color: 'black', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: -1,}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%'}} />
            </div>
            <div className="dropView" style={{position: 'absolute', zIndex: 5, backgroundColor: accentColor, width: '100%', height: '100%', opacity: '0', transition: 'opacity 0.4s ease', display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%'}} />
            </div>
        </div>
    )
}