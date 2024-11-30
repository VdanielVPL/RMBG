import { CSSProperties, ReactNode, useContext } from "react";
import { MainContext } from "../contexts/MainContext";

type ButtonProps = {
    children: ReactNode;
    style?: CSSProperties;
    onClick?: ()=>void;
    disabled?: boolean;
    title?: string;
}

export function Button(props: ButtonProps) {
    const { accentColor } = useContext(MainContext);

    return (
        <button disabled={props.disabled} className="Button" style={{...props.style, ...{backgroundColor: accentColor}}} onClick={props.onClick} title={props.title}>
            {props.children}
        </button>
    )
}