import { useEffect, useRef, useState, CSSProperties } from "react";
import { OnFileDrop, OnFileDropOff } from "../../wailsjs/runtime/runtime";
import { HandleDrop, OpenImage } from "../../wailsjs/go/main/App";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCloudArrowUp } from "@fortawesome/free-solid-svg-icons";

export function InputImageContainer() {
    const [filePath, setFilePath] = useState<string>('');
    const imageContainer = useRef<HTMLDivElement>(null);
    const image = useRef<HTMLImageElement>(null);

    useEffect(() => {
        const handleFileDrop = (x: number, y: number, paths: string[]) => {
            console.log('File dropped:', paths);
            setFilePath(paths[0]);
            HandleDrop(paths[0], false).then(([base64, fileType]) => {
                if (image.current) {
                    if (base64 == "" || fileType == "") {
                        image.current.src = "";

                    }else{
                        image.current.src = `data:${fileType};base64,${base64}`;
                    }
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
                        console.log('Dropped URL:', data);
                        HandleDrop(data, true).then(([base64, fileType]) => {
                            if (image.current) {
                                if (base64 == "" || fileType == "") {
                                    image.current.src = "";
                                }else{
                                    image.current.src = `data:${fileType};base64,${base64}`;
                                }
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
        const result = await OpenImage()
        if (result != null) {
            const [base64, fileType] = result;
            if (image.current) {
                if (base64 == "" || fileType == "") {
                    image.current.src = "";
                }else{
                    image.current.src = `data:${fileType};base64,${base64}`;
                }
            }
        }
    }

    return (
        <div ref={imageContainer} className='imageContainer' onDrop={handleURLDrop} style={{'--wails-drop-target': 'drop'} as CSSProperties} onClick={openDialog}>
            <img ref={image} style={{userSelect: 'none', pointerEvents: 'none'}} draggable={false}></img>
            <div style={{position: 'absolute', height: '100%', width: '100%', color: 'black', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: -1, backgroundColor: 'white'}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%'}} />
            </div>
        </div>
    )
}