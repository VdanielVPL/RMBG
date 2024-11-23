import { useEffect, useRef, useState } from "react";

type CropEditorProps = {
    rect: DOMRect | null;
};

export default function CropEditor(props: CropEditorProps) {
    const [dragging, setDragging] = useState<boolean>(false);
    const cropEditorRef = useRef<HTMLDivElement>(null);

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

    function resetClick(e: React.MouseEvent<HTMLDivElement>) {
        e.preventDefault();
    }

    function mouseDown(e: React.DragEvent<HTMLDivElement>) {
        // console.log(e.movementX)
        // console.log("start")
        setDragging(true);

    }

    function mouseMove(e: React.MouseEvent<HTMLDivElement>) {
        if (dragging && cropEditorRef.current) {
            const currentTop = parseInt(cropEditorRef.current.style.top || '0', 10);
            const newTop = currentTop + e.movementY;

            const minTop = 0;
            const maxTop = props.rect!.height;

            const clampedTop = Math.max(minTop, Math.min(newTop, maxTop));

            cropEditorRef.current.style.top = `${clampedTop}px`;
        }
    }

    return (
        <div ref={cropEditorRef} className="cropEditor" style={{top: '0px',right: '0px', left: '0px', bottom: '0px'}}>
            <div className="top" onClick={resetClick} onMouseDown={mouseDown} onMouseMove={mouseMove}>
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