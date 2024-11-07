import { useEffect, useRef, useState } from "react";
import { OnFileDrop, OnFileDropOff } from "../../wailsjs/runtime/runtime";
import { HandleDrop } from "../../wailsjs/go/main/App";

export function ImageContainer() {
    const [filePath, setFilePath] = useState<string>('');
    const imageContainer = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleFileDrop = (x: number, y: number, paths: string[]) => {
            console.log('File dropped:', paths);
            setFilePath(paths[0]);
            HandleDrop(paths[0], false);
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
                        HandleDrop(data, true)
                        imageContainer.current?.classList.remove('wails-drop-target-active');
                    }
                });
            }
        }
    };

    return (
        <div ref={imageContainer} className='imageContainer' onDrop={handleURLDrop} style={{'--wails-drop-target': 'drop'} as React.CSSProperties}>
            {filePath}
        </div>
    )
}