import React, {useState} from "react";
import {Button, Col, Empty, List, Modal, Pagination, Row, Spin} from "antd";
import {ModalProps} from "antd/es/modal/interface";

import './pluginTool.css'
import {useGetManual} from "../../../hooks/useApis.ts";
import Search from "antd/es/input/Search";
import {useNavigate} from "react-router-dom";

import {AccordionItem} from '../../../components/Accordion/Accordion.tsx';
import '../../../components/Accordion/accordion.css'
import {useCheckPermission} from "../../../hooks/usePermissions.tsx";

// 定义数据类型接口
interface Tool {
    id: string;
    name: string;
    description: string;
}

interface PluginItem {
    icon: string;
    name: string;
    description: string;
    tools: Tool[];
}

interface PageData {
    data: {
        totalRow: number;
        records: PluginItem[];
    };
}

interface PageParams {
    pageNumber: number;
    pageSize: number;
    name: string;
}

export type PluginToolsProps = {
    selectedItem?: any[],
    goToPage: string,
    onSelectedItem?: (item: any) => void,
    onRemoveItem?: (item: any) => void,
    onClose?: () => void
} & ModalProps

export const PluginTools: React.FC<PluginToolsProps> = (props) => {

    const selected = props.selectedItem;
    const [pageParams, setPageParams] = useState<PageParams>({
        pageNumber: 1,
        pageSize: 10,
        name: ''
    });
    const [selectedIndex, setSelectedIndex] = useState(-1);

    const {doGet: getPageData} = useGetManual("/api/v1/aiPlugin/page");
    const queryPermission = useCheckPermission("/api/v1/aiPlugin/query")

    // 为 pageData 添加明确的类型定义
    const [pageData, setPageData] = useState<PageData | undefined>(undefined)
    const [pageDataLoading, setPageDataLoading] = useState(false);

    const fetchData = async (params: PageParams) => {
        if (!queryPermission){
            return;
        }

        setPageDataLoading(true)
        try {
            const resp = await getPageData({
                params: params // 使用传入的params而不是pageParams
            });
            setPageData(resp.data as PageData) // 明确类型转换
        } catch (error) {
            console.error('Failed to fetch data:', error);
        } finally {
            setPageDataLoading(false)
        }
    }

    const navigate = useNavigate();
    const goToPage = props.goToPage;

    const onSearch = (value: string) => {
        const newParams = {
            ...pageParams,
            pageNumber: 1, // 搜索时重置到第一页
            name: value
        };
        setPageParams(newParams);
        fetchData(newParams);
    };

    const changeCurrent = (page: number) => {
        const newParams = {
            ...pageParams,
            pageNumber: page,
        };
        setPageParams(newParams);
        fetchData(newParams);
    }

    return (
        <Modal title={'选择插件'} footer={null} {...props} width={"1000px"}
               height={"600px"}
               afterOpenChange={async (open) => {
                   if (open) {
                       await fetchData(pageParams);
                   }
               }}
        >
            <Row gutter={16} style={{width: "100%"}}>
                <Col span={6} style={{backgroundColor: "#f5f5f5", paddingTop: "10px"}}>
                    <Search style={{marginBottom: "10px"}} placeholder="搜索" onSearch={onSearch}/>
                    <Button block type={"primary"} onClick={() => {
                        goToPage ? navigate(goToPage) : ''
                    }}>创建插件</Button>
                </Col>
                <Col span={18}>
                    <Spin spinning={pageDataLoading}>
                        <div>
                            {(pageData?.data?.totalRow ?? 0) > 0 ? pageData?.data.records.map((item: PluginItem, index: number) => {
                                return (
                                    <AccordionItem
                                        key={index}
                                        icon={item.icon}
                                        title={item.name}
                                        description={item.description}
                                        isActive={selectedIndex === index}
                                        clickItem={() => {
                                            if (selectedIndex === index) {
                                                setSelectedIndex(-1)
                                            } else {
                                                setSelectedIndex(index)
                                            }
                                        }}
                                    >
                                        <List
                                            dataSource={item.tools}
                                            renderItem={(tool: Tool, toolIndex) => (
                                                <List.Item
                                                    key={toolIndex}
                                                    actions={selected?.includes(tool.id) ? [<Button color="danger"
                                                                                                    variant="outlined"
                                                                                                    onClick={() => {
                                                                                                        props.onRemoveItem?.(tool)
                                                                                                    }
                                                                                                    }>
                                                        移除
                                                    </Button>] : [<Button onClick={() => {
                                                        props.onSelectedItem?.(tool)
                                                    }
                                                    }>
                                                        选择
                                                    </Button>]}
                                                >
                                                    <List.Item.Meta
                                                        title={tool.name}
                                                        description={tool.description}
                                                    />
                                                </List.Item>
                                            )}
                                        />
                                    </AccordionItem>
                                )
                            }) : <Empty style={{height: "100%"}} image={Empty.PRESENTED_IMAGE_SIMPLE}/>
                            }
                            <div style={{padding: "10px", backgroundColor: "#f5f5f5"}}>
                                <Pagination
                                    onChange={changeCurrent}
                                    align="end"
                                    current={pageParams.pageNumber}
                                    pageSize={pageParams.pageSize}
                                    total={pageData?.data?.totalRow ?? 0}
                                />
                            </div>
                        </div>
                    </Spin>
                </Col>
            </Row>
        </Modal>
    )
}