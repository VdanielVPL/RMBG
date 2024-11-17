import { FontawesomeObject } from "@fortawesome/fontawesome-svg-core";
import { ReactElement } from "react";
import { Link } from "react-router-dom";

type ButtonNavProps = {
    text: string;
    link?: string;
    icon?: ReactElement<FontawesomeObject>;
}

export function NavButton(props: ButtonNavProps) {

    return (

        <Link className="NavButton" to={props.link || '/'} draggable={false}>
            {props.icon}
            <span>{props.text}</span>
        </Link>

    )
}