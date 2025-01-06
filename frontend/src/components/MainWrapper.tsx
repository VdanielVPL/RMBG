import { ReactNode, useContext, useEffect, useRef, useState } from "react";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";
import { ImageContext } from "./contexts/ImageContext";
import { MainContext } from "./contexts/MainContext";

export function MainWrapper({children}: {children: ReactNode}) {
    const { accentColor, strings } = useContext(MainContext);
    const { setRemovingBG, setCropping } = useContext(ImageContext);
    const [ triggerAlert, setTriggerAlert ] = useState<boolean>(false);
    const [ alertString, setAlertString ] = useState<string>("");

    useEffect(() => {
        document.body.style.setProperty('--accent-color', accentColor);
        const handleRemBGstatus = (status: boolean) => {
            setRemovingBG(status);
        }

        const handleCropStatus = (status: boolean) => {
            setCropping(status);
        }

        const handleAlert = (alert: string) => {
            let alertMessage = '';
            switch (alert) {
                case 'REMBG_NOT_FOUND':
                    alertMessage = strings['ErrorREMBG_NOT_FOUND'];
                    break;
                case 'SAVED':
                    alertMessage = strings['AlertIMAGE_SAVED'];
                    break;
                case 'COPIED':
                    alertMessage = strings['AlertIMAGE_COPIED'];
                    break;
                default:
                    alertMessage = strings['ErrorUNKNOWN'];
            }
            setTriggerAlert((prev) => !prev);
            setAlertString(alertMessage);
        }

        EventsOn('removingbg', handleRemBGstatus);
        EventsOn('cropping', handleCropStatus);
        EventsOn('alert', handleAlert);

        return () => {
            EventsOff('removingbg');
            EventsOff('cropping');
            EventsOff('alert');
        }
    }, [strings, accentColor]);

    return (
        <>
            <div style={{position: 'relative', flex: '1', boxSizing: "border-box", padding: '20px', width: '100%', height: '100%'}}>
                {children}  
            </div>
            <Alert triggerAlert={triggerAlert} alertString={alertString} accentColor={accentColor}/>
        </>
    )
}

function Alert(props: {alertString?: string, accentColor: string, triggerAlert: boolean}) {
    const { isDarkMode } = useContext(MainContext);
    const alertRef = useRef<HTMLDivElement>(null);
    const timeoutRef = useRef<number | null>(null);

    useEffect(() => {
        if (props.alertString !== "") {
            if (timeoutRef.current) {
                clearTimeout(timeoutRef.current);
            }
            alertRef.current?.style.setProperty('opacity', '1');
            timeoutRef.current = window.setTimeout(() => {
                alertRef.current?.style.setProperty('opacity', '0');
                timeoutRef.current = null;
            }, 5000);
        }
    }, [props.triggerAlert, props.alertString]);

    return (
        <div ref={alertRef} style={{pointerEvents: 'none', opacity: 0, position: 'absolute', width: 'auto', height: '20px', bottom: 40, left: '50%', transform: 'translateX(-50%)', padding: '5px', borderRadius: '10px', transition: 'opacity 0.5s ease'}}>
            <div style={{color: isDarkMode?'white':'black', padding: '10px', borderRadius: '10px', backgroundColor: isDarkMode?'rgb(34, 38, 39)':'white', border: `1px solid ${props.accentColor}`, boxShadow: `0px 0px 10px ${props.accentColor}`}}>{props.alertString}</div>
        </div>
    )
}