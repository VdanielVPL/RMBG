import { faCircle } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useEffect, useState } from "react";

interface SwitchProps{
    primaryColor?: string;
    secondaryColor?: string;
    /**
        * Switch height
        * @default 23
    */
    size?: number;
    defaultActive?: boolean;
    /**
        * @param active contains the value the switch is trying to change to. The change works if there is no **value** set in the Switch component. However, if value is set then you need to use **useState** to handle the Switch value change
        * @returns void
    */
    onChange?: (active: boolean) => void;
    /**
        * Determines the state of the switch
    */
    value?: boolean;
    /**
        * Width multiplier for the switch
        * @default 1.85
    */
    widthMultiplier?: number;
    visibility?: 'visible' | 'hidden';
}

export default function Switch(props: SwitchProps){

    const [active, setActive] = useState(props.value || props.defaultActive || false)

    const defaultProps = {
        primaryColor: 'grey',
        secondaryColor: 'darkslateblue',
        size: 23,
        widthMultiplier: 1.85,
        visibility: 'visible' as SwitchProps['visibility']
    };

    const finalProps = { ...defaultProps, ...props };

    function onChangeTrigger(){
        if(props.onChange!==undefined){
            props.onChange(!active)
        }
    }

    function onclick(){
        if(props.value===undefined){
            onChangeTrigger()
            setActive(prev => !prev)
        }else{
            onChangeTrigger()
        }
    }

    useEffect(()=>{
        if(props.value!==undefined){
            setActive(props.value)
        }
    },[props.value])

    return (
        <div className="switchComponentContainer" style={{
            color: finalProps.primaryColor,
            width: finalProps.size*finalProps.widthMultiplier+'px',
            height: finalProps.size+'px',
            borderRadius: '100vh',
            backgroundColor: active?finalProps.secondaryColor:finalProps.primaryColor,
            position: 'relative',
            display: 'flex',
            alignItems: 'center',
            cursor: 'pointer',
            visibility: finalProps.visibility
        }} onClick={onclick}>
            <div className="switchComponentCircle" style={{
                backgroundColor: 'white',
                borderRadius: '100%',
                height: finalProps.size-finalProps.size/5.75+'px',
                width: finalProps.size-finalProps.size/5.75+'px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                position: 'absolute',
                marginLeft: finalProps.size/8.5+'px',
                transform: active?`translateX(${(finalProps.size*finalProps.widthMultiplier)-(finalProps.size-finalProps.size/5.75)-finalProps.size/4.2}px)`:'translateX(0px)'
            }}>
                <FontAwesomeIcon className="switchComponentCircleIcon" icon={faCircle} color={active?finalProps.secondaryColor:finalProps.primaryColor} style={{fontSize: `${finalProps.size/1.5}px`}} />
            </div>
        </div>
    )
}