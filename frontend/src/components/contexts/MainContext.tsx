import { createContext, useEffect, useState, ReactNode } from 'react';
import { GetLangStrings, GetDarkMode, GetAccentColor } from '../../../wailsjs/go/main/App';

type MainContextType = { 
    strings: {[key: string]: string};
    isDarkMode: boolean;
    accentColor: string;
};

export const MainContext = createContext<MainContextType>({
    strings: {},
    isDarkMode: false,
    accentColor: ''
});

export function MainProvider({ children }: { children: ReactNode }) {
    const [contextData, setContextData] = useState<MainContextType>({
        strings: {},
        isDarkMode: false,
        accentColor: ''
    });

    useEffect(() => {
        Promise.all([
            GetLangStrings(),
            GetDarkMode(),
            GetAccentColor()
        ]).then(([strings, darkMode, accentColor]) => {
            setContextData({
                strings: strings,
                isDarkMode: darkMode,
                accentColor: accentColor
            });
        });
    }, []);

    return (
        <MainContext.Provider value={contextData}>
            {children}
        </MainContext.Provider>
    );
}