import { ReactNode } from "react";

export function MainWrapper({children}: {children: ReactNode}) {
    return (
        <div style={{height: "100%", flex: '1', boxSizing: "border-box", padding: '20px'}}>
            {children}
        </div>
    )
}