import {Button, Divider, Input, message, Modal, Pagination} from "antd";
import {ModalProps} from "antd/es/modal/interface";
import React, { useState} from "react";
import {SearchOutlined} from "@ant-design/icons";
import {useNavigate} from "react-router-dom";
import {useCheckPermission} from "../../../hooks/usePermissions.tsx";
import {useGetManual} from "../../../hooks/useApis.ts";


export type BotDataItemProps = {
    title?: string,
    description?: string,
    icon?: string,
    isAdded?: boolean, //
    onButtonClick?: () => void,


}

export const BotDataItem: React.FC<BotDataItemProps> = ({title, description, icon = "/favicon.png", onButtonClick,isAdded = false, }) => {
    return (
        <div style={{
            display: "flex",
            gap: "20px",
            alignItems: "center",
            background: "#efefef",
            padding: "10px",
            borderRadius: "3px"
        }}>
            <div style={{"width": "40px"}}><img src={icon} alt="" style={{maxWidth: "100%"}}/>
            </div>
            <div style={{flexGrow: 1}}>
                <div>{title}</div>
                <div style={{color: "#999999"}}>{description}</div>
            </div>
            <div>
                <Button
                    onClick={() => {
                        if (!isAdded) {
                            onButtonClick?.()
                        }
                    }}
                    disabled={isAdded}
                    type={isAdded ? "default" : "primary"}
                >
                    {isAdded ? "已添加" : "添加"}
                </Button>
            </div>
        </div>
    )
}

export type BotModalProps = {
    tableAlias: string,
    tableTitle: string,
    goToPage: string,
    pageSize?: number,
    addedItems?: any[],
    addedItemsKeyField?: string,
    onSelectedItem?: (item: any) => void,
    onClose?: () => void
} & ModalProps

export const BotModal: React.FC<BotModalProps> = (props) => {
    const {
        pageSize = 10,
        addedItems = [],
        addedItemsKeyField = 'id',
         ...restProps
    } = props;

    const [items, setItems] = useState<any[]>([])
    const [keyword, setKeyword] = useState<string>('')
    const [currentPage, setCurrentPage] = useState<number>(1)
    const [total, setTotal] = useState<number>(0)
    const [loading, setLoading] = useState<boolean>(false)

    const queryPermission = useCheckPermission(`/api/v1/${props.tableAlias}/query`);


    const {
        doGet
    } = useGetManual(`/api/v1/${props.tableAlias}/page`);

    // const {result, doGet} = usePage(props.tableAlias, {pageSize});
    const navigate = useNavigate();

    const isItemAdded = (item: any) => {
        return addedItems.some(addedItem => {
                return addedItem[addedItemsKeyField] === item['id'];
        });
    };

    // 获取数据的通用方法
    const fetchData = (page: number = currentPage, searchKeyword: string = keyword) => {

        
        if (!queryPermission){
            return ;
        }

        setLoading(true)
        const params: any = {
            pageSize,
            pageNumber: page,
        }

        if (searchKeyword.trim()) {
            params.title = searchKeyword.trim()
            params.name = searchKeyword.trim()
            params.title__op = "like"
        }

        doGet({ params }).then((resp) => {
            if (resp?.data?.data?.records) {
                setItems(resp.data.data.records)
                setTotal(resp.data.data.totalRow || 0)
            } else {
                setItems([])
                setTotal(0)
                if (searchKeyword.trim()) {
                    message.error(`没有找到关于 "${searchKeyword}" 相关数据`)
                }
            }
        }).finally(() => {
            setLoading(false)
        })
    }


    // 搜索处理
    const searchHandler = () => {
        setCurrentPage(1)
        fetchData(1, keyword)
    }

    // 分页处理
    const handlePageChange = (page: number) => {
        setCurrentPage(page)
        fetchData(page, keyword)
    }

    // 重置搜索
    const resetSearch = () => {
        setKeyword('')
        setCurrentPage(1)
        fetchData(1, '')
    }

    return (
        <Modal
            title={<div style={{paddingLeft: "0"}}>添加{props.tableTitle}</div>}
            footer={null}
            {...restProps}
            width={"800px"}
            afterOpenChange={ (open) => {

                if (open ) {
                    fetchData(currentPage, '')
                }

            }}
            height={"700px"} // 增加高度以容纳分页器
        >
            <div style={{display: "flex", gap: "20px", marginTop: "20px", height: "580px"}}>
                <div style={{
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    gap: "10px",
                    background: "#efefef",
                    padding: "20px 10px",
                    borderRadius: "3px",
                    minWidth: "140px"
                }}>
                    <div>
                        <Button onClick={() => {
                            props.goToPage ? navigate(props.goToPage) : ''
                        }}>
                            创建{props.tableTitle}
                        </Button>
                    </div>
                    <Divider/>
                </div>

                <div style={{flexGrow: 1, display: "flex", flexDirection: "column"}}>
                    {/* 搜索区域 */}
                    <div style={{display: "flex", gap: "10px", marginBottom: "20px"}}>
                        <Input
                            size="middle"
                            placeholder="请输入关键字"
                            prefix={<SearchOutlined/>}
                            value={keyword}
                            onChange={(e) => setKeyword(e.target.value)}
                            style={{flexGrow: 1}}
                            onKeyDown={(e) => {
                                if (e.key === "Enter") {
                                    searchHandler()
                                }
                            }}
                        />
                        <Button onClick={searchHandler} loading={loading}>搜索</Button>
                        {keyword && (
                            <Button onClick={resetSearch}>重置</Button>
                        )}
                    </div>

                    {/* 内容区域 */}
                    <div style={{
                        flex: 1,
                        display: "flex",
                        flexDirection: "column",
                        gap: "15px",
                        minHeight: "400px",
                        overflow: "auto"
                    }}>
                        {loading ? (
                            <div style={{
                                display: "flex",
                                justifyContent: "center",
                                alignItems: "center",
                                height: "200px",
                                color: "#999"
                            }}>
                                加载中...
                            </div>
                        ) : items.length === 0 ? (
                            <div style={{
                                display: "flex",
                                justifyContent: "center",
                                alignItems: "center",
                                height: "200px",
                                color: "#999"
                            }}>
                                {keyword ? `没有找到关于 "${keyword}" 的数据` : "暂无数据"}
                            </div>
                        ) : (
                            items.map((item) => (
                                <BotDataItem
                                    key={item.id}
                                    onButtonClick={() => props.onSelectedItem?.(item)}
                                    title={item.title || item.name}
                                    description={item.description}
                                    icon={item.icon}
                                    isAdded={isItemAdded(item)}
                                />
                            ))
                        )}
                    </div>

                    {/* 分页区域 */}
                    {total > 0 && (
                        <div style={{
                            display: "flex",
                            justifyContent: "center",
                            marginTop: "20px",
                            borderTop: "1px solid #f0f0f0",
                            paddingTop: "16px"
                        }}>
                            <Pagination
                                current={currentPage}
                                total={total}
                                pageSize={pageSize}
                                onChange={handlePageChange}
                                showSizeChanger={false}
                                showQuickJumper
                                showTotal={(total, range) =>
                                    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`
                                }
                                size="small"
                            />
                        </div>
                    )}
                </div>
            </div>
        </Modal>
    );
}