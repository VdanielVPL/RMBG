import { useEffect, useRef, useState } from "react";

type CropEditorProps = {
    rect: DOMRect | null;
};

export default function CropEditor(props: CropEditorProps) {
    const [dragging, setDragging] = useState<boolean>(false);
    const [draggingSide, setDraggingSide] = useState<string | null>(null);
    const cropEditorRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (cropEditorRef.current) {
            cropEditorRef.current.style.opacity = '1';
        }
    },[]);

    useEffect(() => {

        const handlemouseup = () => {
            if (dragging) {
                setDragging(false);
            }
        }
        
        window.addEventListener('mouseup', handlemouseup);

        const mouseMove = (e: MouseEvent) => {
            if (dragging && cropEditorRef.current) {
                const minTop = 0;
                const maxTop = props.rect!.height - 30*2-1 - parseInt(cropEditorRef.current.style.bottom || '0', 10);
            
                const minBottom = 0;
                const maxBottom = props.rect!.height - 30*2-1 - parseInt(cropEditorRef.current.style.top || '0', 10);
            
                const minLeft = 0;
                const maxLeft = props.rect!.width - 30*2-1 - parseInt(cropEditorRef.current.style.right || '0', 10);
            
                const minRight = 0;
                const maxRight = props.rect!.width - 30*2-1 - parseInt(cropEditorRef.current.style.left || '0', 10);
                if (draggingSide === "top") {
                    const currentTop = parseInt(cropEditorRef.current.style.top || '0', 10);
                    const newTop = currentTop + e.movementY;
        
                    const clampedTop = Math.max(minTop, Math.min(newTop, maxTop));
        
                    cropEditorRef.current.style.top = `${clampedTop}px`;
                } else if (draggingSide === 'bottom') {
                    const currentBottom = parseInt(cropEditorRef.current.style.bottom || '0', 10);
                    const newBottom = currentBottom - e.movementY;
        
                    const clampedBottom = Math.max(minBottom, Math.min(newBottom, maxBottom));
        
                    cropEditorRef.current.style.bottom = `${clampedBottom}px`;
                } else if (draggingSide === 'left') {
                    const currentLeft = parseInt(cropEditorRef.current.style.left || '0', 10);
                    const newLeft = currentLeft + e.movementX;
        
                    const clampedLeft = Math.max(minLeft, Math.min(newLeft, maxLeft));
        
                    cropEditorRef.current.style.left = `${clampedLeft}px`;
                } else if (draggingSide === 'right') {
                    const currentRight = parseInt(cropEditorRef.current.style.right || '0', 10);
                    const newRight = currentRight - e.movementX;
        
                    const clampedRight = Math.max(minRight, Math.min(newRight, maxRight));
        
                    cropEditorRef.current.style.right = `${clampedRight}px`;
                }
            }
        }

        window.addEventListener('mousemove', mouseMove);


        return () => {
            window.removeEventListener('mouseup', handlemouseup);
            window.removeEventListener('mousemove', mouseMove);
        }

    },[dragging]);

    function resetClick(e: React.MouseEvent<HTMLDivElement>) {
        e.preventDefault();
    }

    function mouseDown(side: string) {
        setDraggingSide(side);
        setDragging(true);
    }

    return (
        <div ref={cropEditorRef} className="cropEditor" style={{top: '0px',right: '0px', left: '0px', bottom: '0px', opacity: '0'}}>
            <div className="top" onClick={resetClick}>
                <div onMouseDown={() => mouseDown("top")} style={{height: '30px', width: '60px'}}></div>
            </div>
            <div className="left" onClick={resetClick}>
                <div onMouseDown={() => mouseDown("left")} style={{height: '60px', width: '30px'}}></div>
            </div>
            <div className="bottom" onClick={resetClick}>
                <div onMouseDown={() => mouseDown("bottom")} style={{height: '30px', width: '60px'}}></div>
            </div>
            <div className="right" onClick={resetClick}>
                <div onMouseDown={() => mouseDown("right")} style={{height: '60px', width: '30px'}}></div>
            </div>
        </div>
    );
}