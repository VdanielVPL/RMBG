import { CSSProperties, ReactNode } from "react";

type ButtonProps = {
    children: ReactNode;
    style?: CSSProperties;
    onClick?: ()=>void;
}

export function Button(props: ButtonProps) {
    return (
        <button className="Button" style={props.style} onClick={props.onClick}>
            {props.children}
        </button>
    )
}