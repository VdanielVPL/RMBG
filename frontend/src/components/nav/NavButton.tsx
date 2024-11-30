import { FontawesomeObject } from "@fortawesome/fontawesome-svg-core";
import { ReactElement, useContext } from "react";
import { Link } from "react-router-dom";
import { MainContext } from "../contexts/MainContext";

type ButtonNavProps = {
    text: string;
    link?: string;
    icon?: ReactElement<FontawesomeObject>;
}

export function NavButton(props: ButtonNavProps) {
    const { accentColor } = useContext(MainContext);

    return (

        <Link className="NavButton" style={{backgroundColor: accentColor}} to={props.link || '/'} draggable={false}>
            {props.icon}
            <span>{props.text}</span>
        </Link>

    )
}