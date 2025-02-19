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

    PNGcropImage: string;
    setPNGCropImage: (inputCropImage: string) => void;

    JPGcropImage: string;
    setJPGCropImage: (inputCropImage: string) => void;

    isJPG: boolean;
    setIsJPG: (isJPG: boolean) => void;

    cropDimens: {left: number, right: number, top: number, bottom: number};
    setCropDimens: (cropDimens: {left: number, right: number, top: number, bottom: number}) => void;

    cropping: boolean;
    setCropping: (cropping: boolean) => void;
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
    setCropImage: () => {},

    JPGcropImage: "",
    setJPGCropImage: () => {},

    PNGcropImage: "",
    setPNGCropImage: () => {},

    isJPG: false,
    setIsJPG: () => {},

    cropDimens: {left: 0, right: 0, top: 0, bottom: 0},
    setCropDimens: () => {},

    cropping: false,
    setCropping: () => {}
});

export function ImageProvider({ children }: { children: ReactNode }) {
    const [inputRMBGImage, setInputRMBGImage] = useState<ImageContextType["inputRMBGImage"]>("");
    const [outputRMBGImage, setOutputRMBGImage] = useState<ImageContextType["outputRMBGImage"]>("");
    const [cropImage, setCropImage] = useState<ImageContextType["cropImage"]>("");
    const [JPGcropImage, setJPGCropImage] = useState<ImageContextType["JPGcropImage"]>("");
    const [PNGcropImage, setPNGCropImage] = useState<ImageContextType["PNGcropImage"]>("");
    const [isJPG, setIsJPG] = useState<ImageContextType["isJPG"]>(false);
    const [removingBG, setRemovingBG] = useState<ImageContextType["removingBG"]>(false);
    const [model, setModel] = useState<ImageContextType["model"]>("u2net");
    const [cropDimens, setCropDimens] = useState<ImageContextType["cropDimens"]>({left: 0, right: 0, top: 0, bottom: 0});
    const [cropping, setCropping] = useState<ImageContextType["cropping"]>(false);

    return (
        <ImageContext.Provider value={{
            inputRMBGImage, setInputRMBGImage, 
            outputRMBGImage, setOutputRMBGImage, 
            cropImage, setCropImage,
            JPGcropImage, setJPGCropImage,
            PNGcropImage, setPNGCropImage,
            isJPG, setIsJPG,
            removingBG, setRemovingBG, 
            setModel, model,
            cropDimens, setCropDimens,
            cropping, setCropping
        }}>
            {children}
        </ImageContext.Provider>
    );
}