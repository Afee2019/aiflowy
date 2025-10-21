import React, { useEffect, useState } from 'react';
import {Button, Input, message, Table, Upload, UploadProps} from 'antd';
import type { TableColumnsType } from 'antd';
import {usePostManual} from "../../../hooks/useApis.ts";
import JsonViewer from "./JsonViewer.tsx";
import {PlayCircleOutlined, UploadOutlined} from "@ant-design/icons";
import {Uploader} from "../../../components/Uploader";

export interface DataType {
    key: string;
    name: string;
    type?: string;
    required?: boolean; // 新增字段表示是否必填
    defaultValue: string;
    children?: DataType[];
}



interface CustomTableProps {
    data: DataType[];
    pluginToolTitle: string;
    pluginToolId: string;
}

const CustomTable: React.FC<CustomTableProps> = ({ data, pluginToolTitle, pluginToolId }) => {
    const [tableData, setTableData] = useState<DataType[]>([]);
    const {doPost: doPostPluginToolTest} = usePostManual('/api/v1/aiPluginTool/test')
    const [pluginTestResultJsonString, setPluginTestResultJsonString] = useState<string>('');
    const [pluginTestRequestJsonString, setPluginTestRequestJsonString] = useState<string>('');
    const [activeView, setActiveView] = useState<'request' | 'response'>('response');
    const [errors, setErrors] = useState<Record<string, string>>({});

    // 初始化数据
    useEffect(() => {
        setErrors({})
        setTableData(data);
    }, [data]);
    const columns = (
        handleInputChange: (key: React.Key, value: string) => void
    ): TableColumnsType<DataType> => [
        {
            title: '参数名称',
            dataIndex: 'name',
            key: 'name',
            width: 200,
            render: (text: string, record: DataType) => (
                <span>
                {text}
                    {record.required && <span style={{ color: 'red', marginLeft: 4 }}>*</span>}
            </span>
            ),
        },
        {
            title: '参数值',
            dataIndex: 'defaultValue',
            key: 'defaultValue',
            width: 100,
            render: (text: string, record: DataType) => {
                console.log('record', record)
                return record?.type === 'Object' || record?.type === 'Array' ? (
                    <span>--</span>
                )
                : record?.type ==='File' ? (
                    <div style={{ display: 'flex', flexDirection: 'column' }}>
                        <Uploader maxCount={1} onChange={(file: any) => {
                            console.log('file', file)
                            if (file?.file?.response?.errorCode === 0){
                                handleInputChange(record.key, file.file.response.path)
                            }
                        }}>
                            <Button icon={<UploadOutlined />}>选择文件</Button>
                        </Uploader>
                    </div>
                )
                        :
                (
                    <div style={{ display: 'flex', flexDirection: 'column' }}>
                        <Input
                            value={text}
                            variant="filled"
                            size="middle"
                            onChange={(e) => handleInputChange(record.key, e.target.value)}
                        />
                        {/* 错误提示 */}
                        {errors[record.key as string] && (
                            <span style={{ color: 'red', fontSize: '12px', marginTop: 4 }}>
                            {errors[record.key as string]}
                        </span>
                        )}
                    </div>
                )
            },
        },
    ];
    // 处理 input change
    const handleInputChange = (key: React.Key, value: string) => {
        const updateData = (dataList: DataType[]): DataType[] => {
            return dataList.map((item) => {
                if (item.key === key) {
                    return { ...item, defaultValue: value };
                }
                if (item.children) {
                    return { ...item, children: updateData(item.children) };
                }
                return item;
            });
        };

        setTableData(updateData(tableData));
    };

    // 提交函数：打印或发送数据
    const handleSubmit = () => {
        const newErrors: Record<string, string> = {};

        const missingRequired = tableData.filter(item => {
            if (item.required && !item.defaultValue) {
                newErrors[item.key as string] = '此字段不能为空';
                return true;
            }
            return false;
        });

        if (missingRequired.length > 0) {
            setErrors(newErrors);
            return;
        }

        // 清空错误
        setErrors({});

        // 提交数据
        setPluginTestRequestJsonString(JSON.stringify(tableData));
        doPostPluginToolTest({
            data: {
                inputData: tableData,
                pluginToolId: pluginToolId
            }
        }).then(res => {
            setPluginTestResultJsonString(JSON.stringify(res.data.data));
            setActiveView('response');
        });
    };

    return (
        <>
            <div style={{ width: '100%', height: '100%', display: 'flex', flexDirection: 'row'}}>
                <div style={{ width: '50%', flexShrink: 0 }}>
                    <h2>{pluginToolTitle} 输入参数</h2>
                    <Table<DataType>
                        columns={columns(handleInputChange)}
                        dataSource={tableData}
                        pagination={false}
                        scroll={{ x: 300 }}
                    />
                    <div
                        style={{
                            width: '100%',
                            display: 'flex',
                            flexDirection: 'row',
                            justifyContent: 'flex-end',
                            marginTop: '20px',
                            padding: '10px',
                        }}
                    >
                        <Button type="primary" size="middle" onClick={handleSubmit} icon={<PlayCircleOutlined />}>
                            运行
                        </Button>
                    </div>
                </div>
                <div style={{ width: '50%', flexShrink: 0 }}>
                    <h2>{pluginToolTitle} 测试结果</h2>
                    { pluginTestRequestJsonString ?
                    <div style={{ width: '100%', height: '100%', overflow: 'auto'}}>
                        <div style={{display: 'flex', flexDirection: 'row', gap: '5px', paddingLeft:  '10px'}}>
                            <h2
                                style={{
                                    cursor: 'pointer',
                                    margin:  '0',
                                    color: activeView === 'request' ? '#1890ff' : 'inherit'
                                }}
                                onClick={() => setActiveView('request')}
                            >
                                Request
                            </h2>
                            <h2
                                style={{
                                    cursor: 'pointer',
                                    margin:  '0',
                                    color: activeView === 'response' ? '#1890ff' : 'inherit'
                                }}
                                onClick={() => setActiveView('response')}
                            >
                                Response
                            </h2>
                        </div>
                        {activeView === 'request' ? (
                            <JsonViewer jsonString={pluginTestRequestJsonString} />
                        ) : (
                            <JsonViewer jsonString={pluginTestResultJsonString} />
                        )}
                    </div> :
                    <div style={{width: '100%', height: '100%', display: 'flex', justifyContent: 'center',
                        alignItems: 'center', fontSize: '20px', color: '#B0B2B3'}}>调试结果展示页面</div>
                    }

                </div>
            </div>
        </>
    );
};

export default CustomTable;