import React, {Fragment, useState} from 'react'
import './confirmItem.css'

export type ConfirmItemMultiProps = {
    ref?: any,
    selectionDataType?: string,
    selectionData?: any[],
    label?: string,
    value?: any,
    onChange?: (value: any[]) => void
}

export const ConfirmItemMulti: React.FC<ConfirmItemMultiProps> = (options) => {

    const {selectionDataType} = options

    const [selectedValue, setSelectedValue] = useState<any[]>([]);

    const changeValue = (value: any) => {

        const isSelected = selectedValue.includes(value);

        const newSelectedValue = isSelected
            ? selectedValue.filter(item => item !== value)
            : [...selectedValue, value];

        setSelectedValue(newSelectedValue);

        if (options.onChange) {
            options.onChange(newSelectedValue);
        }
    }

    const renderSelectionComponent = () => {
        return (
            <div className="custom-radio-group">
                {options.selectionData?.map((item, i) => {
                    return (
                        <Fragment key={item + i}>
                            {selectionDataType === 'text' && <div
                                key={item + i}
                                className={`custom-radio-option ${
                                    selectedValue.includes(item) ? 'selected' : ''
                                }`}
                                style={{width: '100%', flexShrink: 0}}
                                onClick={() => changeValue(item)}
                            >
                                {item}
                            </div>}

                            {selectionDataType === 'image' && <div
                                key={item + i}
                                className={`custom-radio-option ${
                                    selectedValue.includes(item) ? 'selected' : ''
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
                                    selectedValue.includes(item) ? 'selected' : ''
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
                                    selectedValue.includes(item) ? 'selected' : ''
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