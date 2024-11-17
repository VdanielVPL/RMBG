import { createContext, ReactNode, useState } from "react";

type ImageContextType = { 
    inputRMBGImage: string;
    setInputRMBGImage: (inputRMBGImage: string) => void;
    outputRMBGImage: string;
    setOutputRMBGImage: (outputRMBGImage: string) => void;
    removingBG: boolean;
    setRemovingBG: (removingBG: boolean) => void;

    cropImage: string;
    setCropImage: (inputCropImage: string) => void;
};

export const ImageContext = createContext<ImageContextType>({
    inputRMBGImage: "",
    setInputRMBGImage: () => {},
    outputRMBGImage: "",
    setOutputRMBGImage: () => {},
    cropImage: "",
    setCropImage: () => {},
    removingBG: false,
    setRemovingBG: () => {}
});

export function ImageProvider({ children }: { children: ReactNode }) {
    const [inputRMBGImage, setInputRMBGImage] = useState<ImageContextType["inputRMBGImage"]>("");
    const [outputRMBGImage, setOutputRMBGImage] = useState<ImageContextType["outputRMBGImage"]>("");
    const [cropImage, setCropImage] = useState<ImageContextType["cropImage"]>("");
    const [removingBG, setRemovingBG] = useState<ImageContextType["removingBG"]>(false); 

    return (
        <ImageContext.Provider value={{inputRMBGImage, setInputRMBGImage, outputRMBGImage, setOutputRMBGImage, cropImage, setCropImage, removingBG, setRemovingBG}}>
            {children}
        </ImageContext.Provider>
    );
}