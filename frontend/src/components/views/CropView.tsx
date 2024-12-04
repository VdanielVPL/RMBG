import { useContext } from "react";
import { InputImageContainer } from "../CropImageContainer";
import { MainContext } from "../contexts/MainContext";
import { Button } from "../inputs/Button";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck, faCloudArrowDown, faCopy, faTrashCan } from "@fortawesome/free-solid-svg-icons";
import { ImageContext } from "../contexts/ImageContext";
import { useNavigate } from "react-router-dom";
import { faImage } from "@fortawesome/free-regular-svg-icons";
import { CropImage, FromCroptoRMBG, CopyImage, SaveImage } from "../../../wailsjs/go/main/App";
import { EventsEmit } from "../../../wailsjs/runtime/runtime";

export function CropView() {
    const { strings, isDarkMode } = useContext(MainContext);
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
        if (cropDimens.left == 0 && cropDimens.right == 0 && cropDimens.top == 0 && cropDimens.bottom == 0) {
            return;
        }
        CropImage(cropDimens.left, cropDimens.right, cropDimens.top, cropDimens.bottom).then((result) => {
            if (result != null) {
                const [base64, fileType] = result;
                if (base64 == "" || fileType == "") {
                    setCropImage("");
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