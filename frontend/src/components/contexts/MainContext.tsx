import { createContext, useEffect, useState, ReactNode } from 'react';
import { GetLangStrings, GetDarkMode } from '../../../wailsjs/go/main/App';

type MainContextType = { 
    strings: {[key: string]: string};
    isDarkMode: boolean;
};

export const MainContext = createContext<MainContextType>({
    strings: {},
    isDarkMode: false
});

export function MainProvider({ children }: { children: ReactNode }) {
    const [contextData, setContextData] = useState<MainContextType>({
        strings: {},
        isDarkMode: false
    });

    useEffect(() => {
        Promise.all([
            GetLangStrings(),
            GetDarkMode()
        ]).then(([strings, darkMode]) => {
            setContextData({
                strings: strings,
                isDarkMode: darkMode
            });
        });
    }, []);

    return (
        <MainContext.Provider value={contextData}>
            {children}
        </MainContext.Provider>
    );
}