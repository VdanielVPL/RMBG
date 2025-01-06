import { faAngleUp } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useEffect, useRef, useState } from "react"

interface SelectListProps{
    tiles: Selected[]
    width?: string;
    onSelected?: (selected: Selected)=>void;
    defaultSelected?: "first" | "none" | string;
    isDarkMode?: boolean;
    accentColor?: string;
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
        }else{
            setSelected(props.tiles.filter(tile=>tile.value==defaultSelected)[0])
        }
    },[props.tiles])

    return (
        <div ref={selectListRef}>
            <div style={{position: "relative",backgroundColor: props.accentColor?props.accentColor:'rgb(34, 38, 39)', borderTopLeftRadius: '10px', borderTopRightRadius: '10px', borderBottomLeftRadius: expanded?'0px':'10px',borderBottomRightRadius: expanded?'0px':'10px', cursor: 'pointer', width: props.width?props.width:'370px'}} onClick={onClickList}>
                <div style={{display: 'grid', gridTemplateColumns: '1fr auto',padding: '5px 5px 5px 10px'}}>
                    <span style={{userSelect: 'none', fontWeight: 500, display: 'inline-block', position: 'relative'}} dangerouslySetInnerHTML={{__html: selected?.text ?? ''}}></span>
                    <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center', transition: 'rotate 0.1s ease', width: '100%', height: '100%', rotate: expanded?'180deg':'0deg', userSelect: 'none'}}>
                        <FontAwesomeIcon icon={faAngleUp} />
                    </div>
                </div>
                <div style={{display:expanded?'block':'none', position: 'absolute', top: '100%', width: '100%', backgroundColor: (props.isDarkMode!==null&&props.isDarkMode)?'rgb(30, 32, 33)':'white', zIndex: '3', borderBottomLeftRadius: '10px', borderBottomRightRadius: '10px', outlineWidth: '2px', outlineColor: props.accentColor, outlineOffset: '-2px', outlineStyle: 'solid'}}>
                    <div style={{padding: '5px 5px 5px 5px', display: 'flex', flexDirection: 'column', gap: '5px'}}>
                        {props.tiles.map((item)=>(
                            <div className="selectListOption" onClick={()=>{onClickOption(item.text, item.value)}} style={{padding: '5px',fontSize: '15px', fontWeight: 500, textAlign: 'left', backgroundColor: item.text==selected?.text?(props.isDarkMode!==null&&props.isDarkMode)?'rgb(34, 38, 39)':'rgb(220, 220, 220)':'inherit', borderRadius: '10px', transition: 'background-color 0.35s ease, color 0.35s ease', color: (props.isDarkMode!==null&&props.isDarkMode)?'white':'black'}} key={item.value}><span style={{position: 'relative'}} dangerouslySetInnerHTML={{__html: item.text}}></span></div>
                        ))}
                    </div>
                </div>
            </div>
        </div>  
    )
}