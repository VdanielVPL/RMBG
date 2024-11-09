import { CSSProperties, ReactNode } from "react";

type ButtonProps = {
    children: ReactNode;
    style?: CSSProperties;
}

export function Button(props: ButtonProps) {
    return (
        <button className="Button" style={props.style}>
            {props.children}
        </button>
    )
}