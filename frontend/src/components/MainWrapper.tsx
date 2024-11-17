import { ReactNode, useContext, useEffect } from "react";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";
import { ImageContext } from "./contexts/ImageContext";

export function MainWrapper({children}: {children: ReactNode}) {

    const { setRemovingBG } = useContext(ImageContext);

    useEffect(() => {
        const handleRemBGstatus = (status: boolean) => {
            setRemovingBG(status);
        }

        EventsOn('removingbg', handleRemBGstatus);

        return () => {
            EventsOff('removingbg');
        }
    }, []);

    return (
        <div style={{height: "100%", flex: '1', boxSizing: "border-box", padding: '20px'}}>
            {children}
        </div>
    )
}