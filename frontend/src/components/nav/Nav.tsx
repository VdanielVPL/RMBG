import { useContext } from "react";
import { MainContext } from "../contexts/MainContext";
import { NavButton } from "./NavButton";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCropSimple } from "@fortawesome/free-solid-svg-icons";
import { faImage } from "@fortawesome/free-regular-svg-icons";

export function Nav() {
    const { strings } = useContext(MainContext);

    const Buttons = [
        {
            text: strings["BGrem"],
            icon: <FontAwesomeIcon icon={faImage} size="xl" style={{transform: 'scale(0.8)'}}/>,
        },
        {
            text: strings["Crop"], 
            icon: <FontAwesomeIcon icon={faCropSimple} size="xl" style={{transform: 'scale(0.8)'}}/>,
            link: '/crop'
        },
    ]

    return (
        <div style={{width: '30%', height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'center', gap: '30px', padding: '20px', boxSizing: "border-box", maxWidth: '315px', minWidth: "fit-content"}}>
            {Buttons.map((button, index) => (
                <NavButton icon={button.icon}key={index} text={button.text} link={button.link}/>
            ))}
        </div>
    )
}