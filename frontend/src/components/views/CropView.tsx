import { useContext } from "react";
import { useNavigate } from "react-router-dom";
import { InputImageContainer } from "../CropImageContainer";
import { ImageContext, MainContext } from "../contexts";
import { Button, Switch } from "../inputs";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck, faCloudArrowDown, faCopy, faTrashCan } from "@fortawesome/free-solid-svg-icons";
import { faImage } from "@fortawesome/free-regular-svg-icons";
import { CropImage, FromCroptoRMBG, CopyImage, SaveImage, ToJPG, ToPNG, ClearImageMem } from "../../../bindings/rmbg/app";

export function CropView() {
    const { strings, isDarkMode, accentColor } = useContext(MainContext);
    const { cropImage, setCropImage, JPGcropImage, setJPGCropImage, setPNGCropImage, PNGcropImage, setInputRMBGImage, setOutputRMBGImage, cropDimens, cropping, setCropping, setIsJPG, isJPG, setCropImageByType } = useContext(ImageContext);

    let navigate = useNavigate();

    function ClearView() {
        setCropImageByType("", "");
        ClearImageMem("CROP");
    }

    function ToRMBGView() {
        if (cropImage != "") {
            setInputRMBGImage(isJPG?JPGcropImage:PNGcropImage);
            setOutputRMBGImage("");
            FromCroptoRMBG(isJPG)
            navigate('/');
        }
    }

    function CropImageButton() {
        if (cropDimens.left == 0 && cropDimens.right == 0 && cropDimens.top == 0 && cropDimens.bottom == 0) return;
        CropImage(cropDimens.left, cropDimens.right, cropDimens.top, cropDimens.bottom, isJPG).then((result) => {
            if (result.length > 0) {
                const [url, fileType] = result;
                setCropImageByType(url, fileType);
                setCropping(false);
            }
        });
    }

    function OnSwitchChange(active: boolean) {
        setIsJPG(active);
        if (active) {
            //PNG to JPG
            if (JPGcropImage == "") {
                ToJPG().then((result) => {
                    if (result.length > 0) {
                        const newImage = result[0];
                        setJPGCropImage(newImage);
                        setCropImage(newImage);
                        setCropping(false);
                    }else{
                        setIsJPG(false);
                    }
                })
            }else{
                setCropImage(JPGcropImage);
            }
        }else{
            //JPG to PNG
            if (PNGcropImage == "") {
                ToPNG().then((result) => {
                    if (result.length > 0) {
                        const newImage = result[0];
                        setPNGCropImage(newImage);
                        setCropImage(newImage);
                        setCropping(false);
                    }else{
                        setIsJPG(true);
                    }
                })
            }else{
                setCropImage(PNGcropImage);
            }
        }
    }

    return (
        <div className="View" style={{display: 'flex', flexDirection: 'column'}}>
            <span style={{color: isDarkMode?'white':'black'}} className='ViewHeader'>{strings["Crop"]}</span>
            <div style={{height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
                <div style={{display: 'flex', flexDirection: 'column', gap: '10px', alignItems: 'center'}}>
                    <div style={{height: '40px', display: 'flex', alignItems: 'center', justifyContent: 'center', width: '100%', gap: '10px'}}>
                        <span style={{visibility: cropImage==""?'hidden':'visible', color: isDarkMode?'white':'black'}}>PNG</span>
                        <Switch visibility={cropImage==""?'hidden':'visible'} primaryColor={accentColor} secondaryColor={accentColor} onChange={OnSwitchChange} value={isJPG} disabled={cropping}></Switch>
                        <span style={{visibility: cropImage==""?'hidden':'visible', color: isDarkMode?'white':'black'}}>JPEG</span>
                    </div>
                    <InputImageContainer></InputImageContainer>
                    <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-between', width: '100%', gap: '10px'}}>
                        <div style={{display: 'flex', flexDirection: 'row', gap: '10px'}}>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible', opacity: cropImage==""?0:1, backgroundColor: cropping?'indigo':undefined}} onClick={ClearView} title={strings["ClearView"]} disabled={cropping}>
                                <FontAwesomeIcon icon={faTrashCan} size='xl' style={{transform: 'scale(0.8)'}} />
                            </Button>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible', opacity: cropImage==""?0:1, backgroundColor: cropping?'indigo':undefined}} onClick={ToRMBGView} title={strings["ToRMBGView"]} disabled={cropping}>
                                <FontAwesomeIcon icon={faImage} size='xl' style={{transform: 'scale(0.8)'}}/>
                            </Button>
                        </div>
                        <div>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible', opacity: cropImage==""?0:1, backgroundColor: cropping?'indigo':undefined}} onClick={CropImageButton} title={strings["CropImage"]} disabled={cropping}>
                                <FontAwesomeIcon icon={faCheck} size='xl' style={{transform: 'scale(0.8)'}} />
                            </Button>
                        </div>
                        <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', gap: '10px'}}>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible'}} onClick={() => CopyImage('CROP', isJPG)} title={strings["CopyImage"]}>
                                <FontAwesomeIcon icon={faCopy} size='xl'style={{transform: 'scale(0.8)'}} />
                            </Button>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible'}} onClick={() => SaveImage('CROP', isJPG)} title={strings["SaveImage"]}>
                                <FontAwesomeIcon icon={faCloudArrowDown} size='xl' style={{transform: 'scale(0.8)'}}/>
                            </Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}