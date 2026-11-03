import { useEffect, useContext } from "react";
import { Events } from "@wailsio/runtime";
import { HandleDrop, OpenImage } from "../../bindings/rmbg/app";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCloudArrowUp } from "@fortawesome/free-solid-svg-icons";
import { MainContext, ImageContext } from "./contexts";

export function InputImageContainer() {
    const { accentColor } = useContext(MainContext);
    const { inputRMBGImage, setInputRMBGImage, setOutputRMBGImage } = useContext(ImageContext);

    function setLoadedImage(result: string[]) {
        if (result.length > 0) {
            setInputRMBGImage(result[0]);
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
        const url = event.dataTransfer.getData('URL') || event.dataTransfer.getData('text/plain');
        if (url.startsWith('http://') || url.startsWith('https://')) {
            HandleDrop("RMBG", url, true).then(setLoadedImage);
        }
    };

    async function openDialog() {
        setLoadedImage(await OpenImage("RMBG"))
    }

    return (
        <div id='RMBG' className='imageContainer' onDrop={handleURLDrop} data-file-drop-target onClick={openDialog}>
            {inputRMBGImage!="" && <img src={inputRMBGImage} style={{userSelect: 'none', pointerEvents: 'none'}} draggable={false}></img>}
            <div style={{position: 'absolute', height: '100%', width: '100%', color: 'black', display: 'flex', justifyContent: 'center', alignItems: 'center', zIndex: -1,}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%'}} />
            </div>
            <div className="dropView" style={{position: 'absolute', zIndex: 5, backgroundColor: accentColor, width: '100%', height: '100%', opacity: '0', transition: 'opacity 0.4s ease', display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                <FontAwesomeIcon icon={faCloudArrowUp} color="lightgray" style={{width: '60%', height: '60%', fontSize: '60%'}} />
            </div>
        </div>
    )
}