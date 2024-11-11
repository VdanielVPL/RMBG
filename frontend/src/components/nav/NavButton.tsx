import { Link } from "react-router-dom";

type ButtonNavProps = {
    text: string;
    link?: string;
}

export function NavButton(props: ButtonNavProps) {

    return (

        <Link className="NavButton" to={props.link || '/'} draggable={false}>
            <span>{props.text}</span>
        </Link>

    )
}