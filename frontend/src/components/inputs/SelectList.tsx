import { useEffect, useRef, useState } from "react"

interface SelectListProps{
    tiles: Selected[]
    width?: string;
    onSelected?: (selected: Selected)=>void;
    defaultSelected?: "first" | "none"

}

export interface Selected{
    text: string;
    value: string;
    isSelected?: boolean;
}

export default function SelectList(props: SelectListProps) {
    const [expanded, setExpanded] = useState<boolean>(false)
    const [selected, setSelected] = useState<Selected | null>(null)
    const selectListRef = useRef<HTMLDivElement>(null);
    let defaultSelected = props.defaultSelected?props.defaultSelected:"first"

    function onClickList(){
        setExpanded((prev)=>!prev)
    }

    function onClickOption(text: string, value: string){
        setSelected({text: text, value: value})
        if(props.onSelected){
            props.onSelected({text: text, value: value})
        }
    }
    
    const handleClickOutside = (event: Event) => {
        if (selectListRef.current && !selectListRef.current.contains(event.target as Node)) {
            setExpanded(false);
        }
    };

    useEffect(()=>{
        if(expanded){
            window.addEventListener("click", handleClickOutside)
        }else{
            window.removeEventListener("click", handleClickOutside)
        }
        return ()=>{
            window.removeEventListener("click", handleClickOutside)
        }
    },[expanded])

    useEffect(()=>{

        let tileSelected = props.tiles.filter(tile=>tile.isSelected==true)

        if(tileSelected.length==1){
            setSelected(tileSelected[0])
        }else if(defaultSelected=="first"){
            setSelected(props.tiles[0])
        }

    },[props.tiles])

    return (
        <div ref={selectListRef}>
            <div style={{position: "relative",backgroundColor: 'rgb(34, 38, 39)', borderTopLeftRadius: '10px', borderTopRightRadius: '10px', borderBottomLeftRadius: expanded?'0px':'10px',borderBottomRightRadius: expanded?'0px':'10px', cursor: 'pointer', width: props.width?props.width:'370px'}} onClick={onClickList}>
                <div style={{display: 'grid', gridTemplateColumns: '1fr auto',padding: '5px 5px 5px 10px'}}>
                    <span style={{userSelect: 'none'}}>{selected?.text}</span>
                    <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center', width: '100%', height: '100%', rotate: expanded?'180deg':'0deg', userSelect: 'none'}}>^</div>
                </div>
                <div style={{display:expanded?'block':'none', position: 'absolute', top: '100%', width: '100%', backgroundColor: 'rgb(30, 32, 33)', zIndex: '1', borderBottomLeftRadius: '10px', borderBottomRightRadius: '10px'}}>
                    <div style={{padding: '5px 5px 5px 5px', display: 'flex', flexDirection: 'column', gap: '5px'}}>
                        {props.tiles.map((item)=>(
                            <div onClick={()=>{onClickOption(item.text, item.value)}} className="selectListOption" style={{padding: '5px', textAlign: 'left', backgroundColor: item.text==selected?.text?'rgb(34, 38, 39)':'inherit', borderRadius: '10px'}} key={item.value}>{item.text}</div>
                        ))}
                    </div>
                </div>
            </div>
        </div>  
    )
}