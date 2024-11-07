import { useContext } from "react";
import { MainContext } from "../contexts/MainContext";
import { NavButton } from "./NavButton";

export function Nav() {
    const { strings } = useContext(MainContext);

    const Buttons = [
        {
            text: strings["BGrem"]
        },
        {
            text: strings["Crop"], 
            link: '/crop'
        },
    ]

    return (
        <div style={{width: '30%', height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'center', gap: '30px', padding: '20px', boxSizing: "border-box", maxWidth: '315px', minWidth: "fit-content"}}>
            {Buttons.map((button, index) => (
                <NavButton key={index} text={button.text} link={button.link}/>
            ))}
        </div>
    )
}