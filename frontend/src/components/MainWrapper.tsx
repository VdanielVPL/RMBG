import { ReactNode, useContext, useEffect, useRef, useState } from "react";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";
import { ImageContext } from "./contexts/ImageContext";
import { MainContext } from "./contexts/MainContext";

export function MainWrapper({children}: {children: ReactNode}) {

    const { accentColor, strings } = useContext(MainContext);
    const { setRemovingBG, setCropping } = useContext(ImageContext);
    const [ triggerError, setTriggerError ] = useState<boolean>(false);
    const [ error, setError ] = useState<string>("");

    useEffect(() => {
        document.body.style.setProperty('--accent-color', accentColor);
        const handleRemBGstatus = (status: boolean) => {
            setRemovingBG(status);
        }

        const handleCropStatus = (status: boolean) => {
            setCropping(status);
        }

        const handleError = (error: string) => {
            let errorMessage = '';
            switch (error) {
                case 'REMBG_NOT_FOUND':
                    errorMessage = strings['ErrorREMBG_NOT_FOUND'];
                    break;
                default:
                    errorMessage = strings['ErrorUnknown'];
            }
            setTriggerError((prev) => !prev);
            setError(errorMessage);
        }

        EventsOn('removingbg', handleRemBGstatus);
        EventsOn('cropping', handleCropStatus);
        EventsOn('error', handleError);

        return () => {
            EventsOff('removingbg');
            EventsOff('cropping');
            EventsOff('error');
        }
    }, [strings, accentColor]);

    return (
        <>
            <div style={{position: 'relative', flex: '1', boxSizing: "border-box", padding: '20px', width: '100%', height: '100%'}}>
                {children}  
            </div>
            <ErrorAlert triggerError={triggerError} error={error} accentColor={accentColor}/>
        </>
    )
}

function ErrorAlert(props: {error?: string, accentColor: string, triggerError: boolean}) {
    const errorRef = useRef<HTMLDivElement>(null);
    const timeoutRef = useRef<number | null>(null);
    const { isDarkMode } = useContext(MainContext);

    useEffect(() => {
        if (props.error !== "") {
            if (timeoutRef.current) {
                clearTimeout(timeoutRef.current);
            }
            errorRef.current?.style.setProperty('opacity', '1');
            timeoutRef.current = window.setTimeout(() => {
                errorRef.current?.style.setProperty('opacity', '0');
                timeoutRef.current = null;
            }, 5000);
        }
    }, [props.triggerError, props.error]);

    return (
        <div ref={errorRef} style={{opacity: 0, position: 'absolute', width: 'auto', height: '20px', bottom: 40, left: '50%', transform: 'translateX(-50%)', padding: '5px', borderRadius: '10px', transition: 'opacity 0.5s ease'}}>
            <div style={{color: isDarkMode?'white':'black', padding: '10px', borderRadius: '10px', backgroundColor: isDarkMode?'rgb(34, 38, 39)':'white', border: `1px solid ${props.accentColor}`, boxShadow: `0px 0px 10px ${props.accentColor}`}}>{props.error}</div>
        </div>
    )
}