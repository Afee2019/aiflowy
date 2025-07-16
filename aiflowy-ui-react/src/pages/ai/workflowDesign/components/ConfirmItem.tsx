import React, {Fragment, useState} from 'react'
import './confirmItem.css'
export type ConfirmItemProps = {
    ref?: any,
    selectionDataType?: string,
    selectionMode?: string,
    selectionData?: any[],
    label?: string,
    value?: any,
    onChange?: (value: any) => void
}

export const ConfirmItem: React.FC<ConfirmItemProps> = (options) => {

    /*const T = [ // dataType
                            { label: "文字", value: "text" },
                            { label: "图片", value: "image" },
                            { label: "视频", value: "video" },
                            { label: "音频", value: "audio" },
                            { label: "文件", value: "file" },
                            { label: "其他", value: "other" }
                        ], x = [ // mode
                            { label: "单选", value: "single" },
                            { label: "多选", value: "multiple" }
                        ];*/

    const {selectionDataType} = options

    const [selectedValue, setSelectedValue] = useState<any>(null);

    const changeValue = (value: any) => {
        if (value === selectedValue) {
            setSelectedValue(null)
            options.onChange?.(null);
        } else {
            setSelectedValue(value);
            options.onChange?.(value);
        }
    }

    const renderSelectionComponent = () => {
        return (
            <div className="custom-radio-group">
                {options.selectionData?.map((item,i) => {
                    return (
                        <Fragment key={item + i}>
                            {selectionDataType === 'text' && <div
                                key={item + i}
                                className={`custom-radio-option ${
                                    selectedValue === item ? 'selected' : ''
                                }`}
                                style={{width: '100%', flexShrink: 0}}
                                onClick={() => changeValue(item)}
                            >
                                {item}
                            </div>}

                            {selectionDataType === 'image' && <div
                                key={item + i}
                                className={`custom-radio-option ${
                                    selectedValue === item ? 'selected' : ''
                                }`}
                                style={{padding: '0'}}
                                onClick={() => changeValue(item)}
                            >
                                <img alt={""} src={item} style={{
                                    width: '80px', height: '80px',
                                    borderRadius: '8px'
                                }} />
                            </div>}

                            {selectionDataType === 'video' && <div
                                key={item + i}
                                className={`custom-radio-option ${
                                    selectedValue === item ? 'selected' : ''
                                }`}
                                onClick={() => changeValue(item)}
                            >
                                <video controls src={item} style={{
                                    width: '162px', height: '141px'
                                }} />
                            </div>}

                            {selectionDataType === 'audio' && <div
                                key={item + i}
                                className={`custom-radio-option ${
                                    selectedValue === item ? 'selected' : ''
                                }`}
                                style={{width: '100%', flexShrink: 0}}
                                onClick={() => changeValue(item)}
                            >
                                <audio controls src={item} style={{
                                    width: '100%', height: '44px',marginTop: '8px'
                                }} />
                            </div>}
                        </Fragment>
                    )
                })}
            </div>
        )
    }
    
    return (
        <>
            {renderSelectionComponent()}
        </>
    )
}