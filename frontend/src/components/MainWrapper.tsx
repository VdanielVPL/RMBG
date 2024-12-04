import { ReactNode, useContext, useEffect } from "react";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";
import { ImageContext } from "./contexts/ImageContext";

export function MainWrapper({children}: {children: ReactNode}) {

    const { setRemovingBG, setCropping } = useContext(ImageContext);

    useEffect(() => {
        const handleRemBGstatus = (status: boolean) => {
            setRemovingBG(status);
        }

        const handleCropStatus = (status: boolean) => {
            setCropping(status);
        }

        EventsOn('removingbg', handleRemBGstatus);
        EventsOn('cropping', handleCropStatus);

        return () => {
            EventsOff('removingbg');
            EventsOff('cropping');
        }
    }, []);

    return (
        <div style={{position: 'relative', flex: '1', boxSizing: "border-box", padding: '20px', width: '100%', height: '100%'}}>
            {children}
        </div>
    )
}