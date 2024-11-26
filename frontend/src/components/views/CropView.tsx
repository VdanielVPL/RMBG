import { useContext } from "react";
import { InputImageContainer } from "../CropImageContainer";
import { MainContext } from "../contexts/MainContext";
import { Button } from "../inputs/Button";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck, faTrashCan } from "@fortawesome/free-solid-svg-icons";
import { ImageContext } from "../contexts/ImageContext";
import { useNavigate } from "react-router-dom";
import { faImage } from "@fortawesome/free-regular-svg-icons";
import { CropImage, FromCroptoRMBG } from "../../../wailsjs/go/main/App";

export function CropView() {
    const { strings } = useContext(MainContext);
    const { cropImage, setCropImage, setInputRMBGImage, cropDimens } = useContext(ImageContext);
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
                }
            }
        });
    }

    return (
        <div style={{display: 'flex', flexDirection: 'column', height: '100%'}}>
            <span className='ViewHeader'>{strings["Crop"]}</span>
            <div style={{height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
                <div style={{display: 'flex', flexDirection: 'column', gap: '10px'}}>
                    <InputImageContainer></InputImageContainer>
                    <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-between', width: '100%', gap: '10px'}}>
                        <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'start', width: '100%', gap: '10px'}}>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible', opacity: cropImage==""?0:1}} onClick={clearView} title={strings["ClearView"]}>
                                <FontAwesomeIcon icon={faTrashCan} size='xl' style={{transform: 'scale(0.8)'}} />
                            </Button>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible', opacity: cropImage==""?0:1}} onClick={ToRMBGView} title={strings["ToRMBGView"]}>
                                <FontAwesomeIcon icon={faImage} size='xl' style={{transform: 'scale(0.8)'}}/>
                            </Button>
                        </div>
                        <div>
                            <Button style={{width: '40px', height: '40px', visibility: cropImage==""?'hidden':'visible', opacity: cropImage==""?0:1}} onClick={CropImageButton} title={strings["CropImage"]}>
                                <FontAwesomeIcon icon={faCheck} size='xl' style={{transform: 'scale(0.8)'}} />
                            </Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}