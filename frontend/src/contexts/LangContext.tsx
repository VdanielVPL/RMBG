import { createContext, useEffect, useState, ReactNode } from 'react';
import { GetLangStrings } from '../../wailsjs/go/main/App';

type LangContextType = { [key: string]: string };

export const LangContext = createContext<LangContextType>({});

export function LangProvider({ children }: { children: ReactNode }) {
    const [langStrings, setLangStrings] = useState<LangContextType>({});

    useEffect(() => {
        GetLangStrings().then(strings => {
            setLangStrings(strings);
        });
    }, []);

    return (
        <LangContext.Provider value={langStrings}>
            {children}
        </LangContext.Provider>
    );
}