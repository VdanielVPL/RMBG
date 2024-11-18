import { createContext, ReactNode, useState } from "react";

type ImageContextType = { 
    inputRMBGImage: string;
    setInputRMBGImage: (inputRMBGImage: string) => void;
    outputRMBGImage: string;
    setOutputRMBGImage: (outputRMBGImage: string) => void;

    removingBG: boolean;
    setRemovingBG: (removingBG: boolean) => void;

    model: string;
    setModel: (model: string) => void;
    
    cropImage: string;
    setCropImage: (inputCropImage: string) => void;
};

export const ImageContext = createContext<ImageContextType>({
    inputRMBGImage: "",
    setInputRMBGImage: () => {},
    outputRMBGImage: "",
    setOutputRMBGImage: () => {},

    removingBG: false,
    setRemovingBG: () => {},

    model: "",
    setModel: () => {},

    cropImage: "",
    setCropImage: () => {}
});

export function ImageProvider({ children }: { children: ReactNode }) {
    const [inputRMBGImage, setInputRMBGImage] = useState<ImageContextType["inputRMBGImage"]>("");
    const [outputRMBGImage, setOutputRMBGImage] = useState<ImageContextType["outputRMBGImage"]>("");
    const [cropImage, setCropImage] = useState<ImageContextType["cropImage"]>("");
    const [removingBG, setRemovingBG] = useState<ImageContextType["removingBG"]>(false);
    const [model, setModel] = useState<ImageContextType["model"]>("u2net");

    return (
        <ImageContext.Provider value={{
            inputRMBGImage, setInputRMBGImage, 
            outputRMBGImage, setOutputRMBGImage, 
            cropImage, setCropImage, 
            removingBG, setRemovingBG, 
            setModel, model
        }}>
            {children}
        </ImageContext.Provider>
    );
}