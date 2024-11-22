import { useEffect, useState } from "react";

type CropEditorProps = {
    rect: DOMRect | null;
};

export default function CropEditor(props: CropEditorProps) {
    const [dragging, setDragging] = useState<boolean>(false);

    useEffect(() => {

        const handlemouseup = (e: MouseEvent) => {
            if (dragging) {
                setDragging(false);
            }
        }
        window.addEventListener('mouseup', handlemouseup);

        return () => {
            window.removeEventListener('mouseup', handlemouseup);
        }

    },[dragging]);

    useEffect(() => {
        console.log(dragging);
    }, [dragging]);

    function resetClick(e: React.MouseEvent<HTMLDivElement>) {
        e.preventDefault();
    }

    function mouseDown(e: React.DragEvent<HTMLDivElement>) {
        // e.preventDefault();
        console.log("start")
        setDragging(true);

    }
    return (
        <div className="cropEditor" style={{width: props.rect?.width, height: props.rect?.height}}>
            <div className="top" onClick={resetClick} onMouseDown={mouseDown}>
                <div></div>
            </div>
            <div className="left" onClick={resetClick}>
                <div></div>
            </div>
            <div className="bottom" onClick={resetClick}>
                <div></div>
            </div>
            <div className="right" onClick={resetClick}>
                <div></div>
            </div>
        </div>
    );
}