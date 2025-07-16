import React, {useEffect, useRef, useState} from 'react';
import TreeClassifiedPage from "../../components/TreeClassifiedPage";
import CrudPage from "../../components/CrudPage";
import {ColumnsConfig} from "../../components/AntdCrud";
import {Form, FormInstance, Input, Select, Switch, Tag} from "antd";
import CustomLeftArrowIcon from "../../components/CustomIcon/CustomLeftArrowIcon.tsx";
import KeywordSearchForm from "../../components/AntdCrud/KeywordSearchForm.tsx";
import {RawValueType} from 'rc-select/lib/interface';


const columns: ColumnsConfig<any> = [
    {
        title: '标识',
        key: 'id',
        hidden: true,
        form: {
            type: "hidden"
        }
    },
    {
        title: '自定义输入标识',
        key: 'llmExtraConfig.isCustomInput',
        hidden: true,
        form: {
            type: "hidden"
        }
    },
    {
        title: '自定义输入标识',
        key: 'isCustomInput',
        hidden: true,
        form: {
            type: "hidden"
        }
    },
    {
        title: 'Icon',
        dataIndex: 'icon',
        key: 'icon',
        form: {
            type: "image",
        },
        render: (text: any) => <img src={text} alt="" style={{width: 30, height: 30}}/>,
    },
    {
        title: '名称',
        dataIndex: 'title',
        key: 'title',
        placeholder: "请输入名称",
        supportSearch: true,
        form: {
            type: "input",
            rules: [{required: true, message: '请输入名称'}],
        }
    },
    {
        title: '供应商',
        dataIndex: 'brand',
        key: 'brand',
        hidden: true,
        dict: '/api/v1/aiLlmBrand/list?asTree=true',
        form: {
            type: "select",
            attrs: {
                fieldNames: {
                    label: "title",
                    value: "key"
                }
            },
            rules: [{required: true, message: '请选择模型供应商'}]
        },

    },
    {
        title: '模型名称',
        dataIndex: 'llmModel',
        key: 'llmModel',
        hidden: true,
        form: {
            type: "input",
            rules: [{required: true, message: '请输入模型名称'}],
        },
    },
    {
        title: '请求地址',
        dataIndex: 'llmEndpoint',
        key: 'llmEndpoint',
        hidden: true,
        editCondition: (data) => {
            return  data?.isCustomInput
        }

    },
    {
        title: 'API Key',
        dataIndex: 'llmApiKey',
        key: 'llmApiKey',
        hidden: true,
    },
    {
        title: 'API Secret',
        dataIndex: 'llmExtraConfig.apiSecret',
        key: 'llmExtraConfig.apiSecret',
        hidden: true,
        form: {
            type: "input",
        },
        editCondition: (data: any) => {
            return data?.brand === "spark" || data?.brand === "ollama"
        }
    },
    {
        title: 'App Id',
        dataIndex: 'llmExtraConfig.appId',
        key: 'llmExtraConfig.appId',
        hidden: true,
        form: {
            type: "input",
        },
        editCondition: (data: any) => {
            return data?.brand === "spark"
        }
    },
    {
        title: 'Version',
        dataIndex: 'llmExtraConfig.version',
        key: 'llmExtraConfig.version',
        hidden: true,
        form: {
            type: "input",
        },
        editCondition: (data) => {
            return data?.brand === "spark";
        }
    },
    {
        title: '对话路径',
        dataIndex: 'llmExtraConfig.chatPath',
        key: 'llmExtraConfig.chatPath',
        hidden: true,
        form: {type: 'input'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '向量化路径',
        dataIndex: 'llmExtraConfig.embedPath',
        key: 'llmExtraConfig.embedPath',
        hidden: true,
        form: {type: 'input'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '能力',
        dataIndex: 'supportFeatures',
        key: 'supportFeatures',
        editCondition: () => false,
        render: (features: any) => {
            if (!features) {
                return ''
            }

            return features.map((feature: any) => <Tag key={feature}>{feature}</Tag>)
        },
    },
    {
        title: '描述',
        dataIndex: 'description',
        key: 'description',
        width: "20%",
        form: {
            type: "TextArea",
            attrs: {
                rows: 3
            }
        }
    },
    {
        title: '对话模型',
        dataIndex: 'supportChat',
        key: 'supportChat',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '方法调用',
        dataIndex: 'supportFunctionCalling',
        key: 'supportFunctionCalling',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (_data: any) => {
            return false;
        },
    },
    {
        title: '向量化',
        dataIndex: 'supportEmbed',
        key: 'supportEmbed',
        hidden: true,
        form: {
            type: 'switch',
        },
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '重排',
        dataIndex: 'supportReranker',
        key: 'supportReranker',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '文生图',
        dataIndex: 'supportTextToImage',
        key: 'supportTextToImage',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '图生图',
        dataIndex: 'supportImageToImage',
        key: 'supportImageToImage',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '文生音频',
        dataIndex: 'supportTextToAudio',
        key: 'supportTextToAudio',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '音频转音频',
        dataIndex: 'supportAudioToAudio',
        key: 'supportAudioToAudio',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '文生视频',
        dataIndex: 'supportTextToVideo',
        key: 'supportTextToVideo',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '图片转视频',
        dataIndex: 'supportImageToVideo',
        key: 'supportImageToVideo',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
    {
        title: '多模态',
        dataIndex: 'options.multimodal',
        key: 'options.multimodal',
        hidden: true,
        form: {type: 'switch'},
        editCondition: (data) => {
            return data?.brand && data?.isCustomInput;
        },
    },
];

const ModelNameInput: React.FC<{
    treeData: any[];
    columnConfig: any;
    form: FormInstance;
    value?: string;
    onChange?: (value: string) => void;
}> = ({
          treeData,
          columnConfig,
          form,
          value,
          onChange,
      }) => {
    const [modelOptions, setModelOptions] = useState([]);
    const [isCustomInput, setIsCustomInput] = useState<boolean | null>(null);

    // 监听brand字段变化
    const brandValue = Form.useWatch('brand', form);

    // 使用 useRef 来追踪初始化状态，避免无限循环
    const hasInitialized = useRef(false);
    const prevBrandValue = useRef(brandValue);


    // 当品牌变化时，需要判断是否为编辑状态
    useEffect(() => {
        const brandData = treeData.find(brand => brand.key === brandValue);
        const options = brandData?.options.modelList?.map((model: { title: any; llmModel: any; }) => ({
            label: model.title,
            value: model.llmModel
        })) || [];
        setModelOptions(options);

        // 只有在品牌真正发生变化时才执行清空逻辑
        const brandChanged = prevBrandValue.current !== brandValue;
        prevBrandValue.current = brandValue;

        if (!brandValue) {
            form.setFieldValue("brand", undefined);
            form.setFieldValue("llmEndpoint", undefined);
            form.setFieldValue("llmExtraConfig.chatPath", undefined);
            form.setFieldValue("llmExtraConfig.embedPath", undefined);

            // 只有在非编辑状态下且品牌发生变化时才清空模型名称和描述
            const formId = form.getFieldValue('id');
            if (!formId && brandChanged && hasInitialized.current) {
                form.setFieldValue("llmModel", undefined);
                form.setFieldValue("description", undefined);
                onChange?.('');
            }


        } else {
            form.setFieldValue("brand", brandValue);
            if (brandData) {
                const options = brandData.options;
                form.setFieldValue("llmEndpoint", options.llmEndpoint);
                form.setFieldValue("llmExtraConfig.chatPath", options.chatPath);
                form.setFieldValue("llmExtraConfig.embedPath", options.embedPath);

                // 品牌变化时，如果不是编辑状态，清空模型名称
                const formId = form.getFieldValue('id');
                if (!formId && brandChanged && hasInitialized.current) {
                    form.setFieldValue("llmModel", undefined);
                    onChange?.('');
                }
            }

            columns.forEach((column) => {
                if (column.editCondition && column.onValuesChange) {
                    column.onValuesChange("", "");
                }
            });

        }

        hasInitialized.current = true;
    }, [brandValue, treeData]);

    useEffect(() => {
        if (isCustomInput == null) {
            const dbIsCustomInput = form.getFieldValue('llmExtraConfig.isCustomInput');


            if (dbIsCustomInput !== undefined && dbIsCustomInput !== null) {
                setIsCustomInput(dbIsCustomInput);
                form.setFieldValue("isCustomInput", dbIsCustomInput);
            } else {
                // 如果数据库中没有这个值，根据当前模型值判断是否应该使用自定义输入
                if (value && brandValue && treeData.length > 0) {
                    const brandData = treeData.find(brand => brand.key === brandValue);
                    const modelList = brandData?.options?.modelList || [];
                    const isModelInList = modelList.some((model: any) => model.llmModel === value);
                    setIsCustomInput(!isModelInList);
                    form.setFieldValue("isCustomInput", !isModelInList);
                } else if (!value) {
                    setIsCustomInput(false);
                    form.setFieldValue("isCustomInput", false);
                }
            }
        }

        columns.forEach((column) => {
            if (column.editCondition && column.onValuesChange) {
                column.onValuesChange("", "");
            }
        });


    }, [value, brandValue, treeData]);




    useEffect(() => {
        const updatedFields = {
            brand: brandValue,
        };

        form.setFieldsValue(updatedFields);
        form.setFieldValue("llmExtraConfig.isCustomInput", isCustomInput);
        form.setFieldValue("isCustomInput", isCustomInput);

        const allValues = form.getFieldsValue();

        columns.forEach(column => {
            if (column.onValuesChange){
                column.onValuesChange(column, "");
            }
        })




    }, [isCustomInput]);

    const findModelListByBrandKey = (brandKey: string) => {
        const options = treeData.find(tree => tree.key === brandKey)?.options;
        let llmList = [];
        if (options) {
            llmList = options.modelList;
        }
        return llmList;
    }

    const setFormSupportValues = (form: FormInstance, model: any) => {
        form.setFieldValue("supportChat", model.supportChat);
        form.setFieldValue("supportFunctionCalling", model.supportFunctionCalling);
        form.setFieldValue("supportEmbed", model.supportEmbed);
        form.setFieldValue("supportReranker", model.supportReranker);
        form.setFieldValue("supportTextToImage", model.supportTextToImage);
        form.setFieldValue("supportImageToImage", model.supportImageToImage);
        form.setFieldValue("supportTextToAudio", model.supportTextToAudio);
        form.setFieldValue("supportAudioToAudio", model.supportAudioToAudio);
        form.setFieldValue("supportTextToVideo", model.supportTextToVideo);
        form.setFieldValue("supportImageToVideo", model.supportImageToVideo);

        if (model.multimodal) {
            form.setFieldValue("options.multimodal", model.multimodal);
        } else {
            form.setFieldValue("options.multimodal", false);
        }

        if (model.version) {
            form.setFieldValue("llmExtraConfig.version", model.version);
        }


        if (model.description) {
            form.setFieldValue("description", model.description);
        }
    }

    return (
        <div>
            {isCustomInput ? (
                <Input
                    placeholder={columnConfig.placeholder}
                    {...columnConfig.form?.attrs}
                    value={value}
                    onChange={(e) => {
                        const inputValue = e.target.value;
                        onChange?.(inputValue);
                        form.setFieldValue('llmModel', inputValue);
                    }}
                />
            ) : (
                <Select
                    allowClear
                    showSearch
                    placeholder={columnConfig.placeholder}
                    {...columnConfig.form?.attrs}
                    options={modelOptions}
                    value={value}
                    onChange={(selectValue) => {
                        onChange?.(selectValue);
                        form.setFieldValue('llmModel', selectValue);


                        // 设置其他相关字段
                        const llmList = findModelListByBrandKey(brandValue);
                        if (llmList) {
                            const model = llmList.find((llm: { llmModel: any; }) => llm.llmModel === selectValue);
                            if (model) {
                                setFormSupportValues(form, model);
                            }
                        }
                    }}
                    filterOption={(input, option) => {
                        return ((option?.label ?? '') as string).toLowerCase().includes(input.toLowerCase())
                    }}
                    optionRender={(option, _info) => {
                        const llmList = findModelListByBrandKey(brandValue);

                        if (llmList && llmList.length) {
                            const currentModel = llmList.find((model: {
                                llmModel: RawValueType | undefined;
                            }) => model.llmModel === option.value);

                            if (currentModel) {
                                return (
                                    <div style={{display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap'}}>
                                        <span>{option.label}</span>
                                        <div style={{display: 'flex', gap: '4px', flexWrap: 'wrap'}}>
                                            {currentModel.supportChat && (
                                                <Tag color="blue">文本聊天</Tag>
                                            )}
                                            {currentModel.multimodal && (
                                                <Tag color="purple">多模态</Tag>
                                            )}
                                            {currentModel.supportEmbed && (
                                                <Tag color="orange">向量化</Tag>
                                            )}
                                        </div>
                                    </div>
                                );
                            }
                        }

                        return <span>{option.label}</span>;
                    }}
                />
            )}
            <div style={{marginTop: 8}}>
                <Switch
                    size="small"
                    checked={isCustomInput}
                    onChange={(checked) => {
                        setIsCustomInput(checked);
                        // 切换时清空当前值
                        onChange?.('');
                        form.setFieldValue('llmModel', '');
                        form.setFieldValue('llmExtraConfig.isCustomInput', checked);
                        form.setFieldValue('isCustomInput', checked);

                    }}
                    checkedChildren="自定义"
                    unCheckedChildren="选择"
                />
            </div>
        </div>
    );
}

const Llms: React.FC<{ paramsToUrl: boolean }> = () => {

    const crudRef = useRef<{ openAddModal: () => void, onSearch: (values: any) => void, formReset: () => void }>(null);

    const [groupId, setGroupId] = useState('')
    const [treeData, setTreeData] = useState<any[]>([]);


    // 自定义表单渲染工厂函数
    const customFormRenderFactory = (position: "edit" | "search", columnConfig: any, form: FormInstance): JSX.Element | null => {
        if (position === 'search') return null;

        // 自定义模型名称字段
        if (columnConfig.key === 'llmModel') {
            return <ModelNameInput treeData={treeData} columnConfig={columnConfig} form={form}/>
        }

        if (columnConfig.key === 'supportChat') {

            const originOnValuesChange = columnConfig.onValuesChange;

            columnConfig.onValuesChange = (values: any, allChangeValues: any) => {
                if (originOnValuesChange) {
                    originOnValuesChange(values, allChangeValues);
                }
                if (values.supportChat) {
                    form.setFieldValue("supportFunctionCalling", values.supportChat);
                }
            }
        }

        return null;
    };

    return (
        <div style={{margin: "24px"}}>
            <KeywordSearchForm
                setIsEditOpen={() => {
                    crudRef.current?.openAddModal()
                }}
                addButtonText="新增大模型"
                columns={columns}
                tableAlias={"aiLlm"}
                onSearch={(values: any) => {
                    crudRef.current?.onSearch(values)
                }}/>

            <TreeClassifiedPage treeTableAlias="aiLlmBrand"
                                treeCardTitle={"接入平台"}
                                treeEditable={false}
                                onTreeDataLoaded={setTreeData}
                                treeTitleIconRender={(item) => {
                                    if (typeof item.icon === "string" && (item.icon.startsWith("http") || item.icon.startsWith("/"))) {
                                        return <img src={item.icon} alt={item.title} style={{width: 14, height: 14}}/>
                                    }
                                    return React.isValidElement(item.icon)
                                        ? item.icon
                                        : <div dangerouslySetInnerHTML={{__html: item.icon || ""}}
                                               style={{padding: "3px"}}/>
                                }}
                                treeCardExtra={
                                    <CustomLeftArrowIcon/>
                                }
                                onTreeSelect={setGroupId}>
                <CrudPage columnsConfig={columns} tableAlias="aiLlm" params={{brand: groupId}} key={groupId}
                          needHideSearchForm={true}
                          editLayout={{openType: "modal"}} ref={crudRef} formRenderFactory={customFormRenderFactory}/>
            </TreeClassifiedPage>
        </div>
    )
};

export default {
    path: "/ai/llms",
    element: Llms
};