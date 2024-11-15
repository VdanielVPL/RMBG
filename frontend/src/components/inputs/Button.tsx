import { CSSProperties, ReactNode } from "react";

type ButtonProps = {
    children: ReactNode;
    style?: CSSProperties;
    onClick?: ()=>void;
    disabled?: boolean;
    title?: string;
}

export function Button(props: ButtonProps) {
    return (
        <button disabled={props.disabled} className="Button" style={props.style} onClick={props.onClick} title={props.title}>
            {props.children}
        </button>
    )
}