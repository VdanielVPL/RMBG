import { createContext, ReactNode, useState } from "react";

type ImageContextType = { 
    inputRMBGImage: string;
    setInputRMBGImage: (inputRMBGImage: string) => void;
    outputRMBGImage: string;
    setOutputRMBGImage: (outputRMBGImage: string) => void;
};

export const ImageContext = createContext<ImageContextType>({
    inputRMBGImage: "",
    setInputRMBGImage: () => {},
    outputRMBGImage: "",
    setOutputRMBGImage: () => {},
});

export function ImageProvider({ children }: { children: ReactNode }) {
    const [inputRMBGImage, setInputRMBGImage] = useState<ImageContextType["inputRMBGImage"]>("");
    const [outputRMBGImage, setOutputRMBGImage] = useState<ImageContextType["outputRMBGImage"]>("");

    return (
        <ImageContext.Provider value={{inputRMBGImage, setInputRMBGImage, outputRMBGImage, setOutputRMBGImage}}>
            {children}
        </ImageContext.Provider>
    );
}