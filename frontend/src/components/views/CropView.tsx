import { useContext } from "react";
import { useNavigate } from "react-router-dom";
import { InputImageContainer } from "../CropImageContainer";
import { ImageContext, MainContext } from "../contexts";
import { Button, Switch } from "../inputs";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck, faCloudArrowDown, faCopy, faTrashCan } from "@fortawesome/free-solid-svg-icons";
import { faImage } from "@fortawesome/free-regular-svg-icons";
import { CropImage, FromCroptoRMBG, CopyImage, SaveImage } from "../../../wailsjs/go/main/App";
import { EventsEmit } from "../../../wailsjs/runtime/runtime";

export function CropView() {
    const { strings, isDarkMode, accentColor } = useContext(MainContext);
    const { cropImage, setCropImage, setInputRMBGImage, cropDimens, cropping } = useContext(ImageContext);
    let navigate = useNavigate();

    function clearView() {
        setCropImage("");
    }

    function ToRMBGView() {
        if (cropImage != "") {
            setInputRMBGImage(cropImage);
            FromCroptoRMBG()
            navigate('/');
        }
    }

    function CropImageButton() {
        if (cropDimens.left == 0 && cropDimens.right == 0 && cropDimens.top == 0 && cropDimens.bottom == 0) return;
        CropImage(cropDimens.left, cropDimens.right, cropDimens.top, cropDimens.bottom).then((result) => {
            if (result != null) {
                const [base64, fileType] = result;
                if (base64 == "" || fileType == "") {
                    setCropImage("");
                    EventsEmit('cropping', false);
                }else{
                    setCropImage(`data:${fileType};base64,${base64}`);
                    EventsEmit('cropping', false);
                }
            }
        });
    }

    return (
        <div className="View" style={{display: 'flex', flexDirection: 'column'}}>
            <span style={{color: isDarkMode?'white':'black'}} className='ViewHeader'>{strings["Crop"]}</span>
            <div style={{height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
                <div style={{display: 'flex', flexDirection: 'column', gap: '10px', alignItems: 'center'}}>
                    <div style={{height: '40px', display: 'flex', alignItems: 'center', justifyContent: 'center', width: '100%', gap: '10px'}}>
                        <span style={{visibility: cropImage==""?'hidden':'visible'}}>PNG</span>
                        <Switch visibility={cropImage==""?'hidden':'visible'} primaryColor={accentColor} secondaryColor={accentColor}></Switch>
                        <span style={{visibility: cropImage==""?'hidden':'visible'}}>JPEG</span>
                    </div>
                    <InputImageContainer></InputImageContainer>
                    <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-between', width: '100%', gap: '10px'}}>
                        <div style={{display: 'flex', flexDirection: 'row', gap: '10px'}}>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible', opacity: cropImage==""?0:1, backgroundColor: cropping?'indigo':undefined}} onClick={clearView} title={strings["ClearView"]} disabled={cropping}>
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
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible'}} onClick={() => CopyImage('CROP')} title={strings["CopyImage"]}>
                                <FontAwesomeIcon icon={faCopy} size='xl'style={{transform: 'scale(0.8)'}} />
                            </Button>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible'}} onClick={() => SaveImage('CROP')} title={strings["SaveImage"]}>
                                <FontAwesomeIcon icon={faCloudArrowDown} size='xl' style={{transform: 'scale(0.8)'}}/>
                            </Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}