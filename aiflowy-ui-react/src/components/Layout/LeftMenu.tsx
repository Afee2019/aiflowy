import React from 'react';
import { Menu } from 'antd';
import { useNavigate } from "react-router-dom";
import Sider from "antd/es/layout/Sider";
import logo from "/favicon.png";
import { useMenus } from "../../hooks/useMenus.tsx";

const LeftMenu: React.FC<{ collapsed: boolean }> = ({ collapsed }) => {
    const navigate = useNavigate();
    const { loading, menuItems, selectItems } = useMenus();
    const selectMenuKeys: string[] = selectItems.map((item) => item.key as string);

    return (
        <Sider
            width={280}
            collapsed={collapsed}
            style={{
                background: "#f5f5f5",
                overflow: 'hidden',
                height: '100vh',
                position: 'sticky',
                top: 0,
                left: 0
            }}
        >
            <style>
                {`
                /* 自定义滚动条样式 */
                ::-webkit-scrollbar {
                    width: 10px; /* 滚动条宽度 */
                    height: 6px; /* 滚动条高度 */
                }
                ::-webkit-scrollbar-thumb {
                    background: rgba(0, 0, 0, 0.2);
                    border-radius: 3px;
                }
                ::-webkit-scrollbar-thumb:hover {
                    background: rgba(0, 0, 0, 0.4);
                }
                ::-webkit-scrollbar-track {
                    background: transparent;
                }
                `}
            </style>

            <div
                style={{
                    display: "flex",
                    flexDirection: "column",
                    height: "100%",
                    overflow: 'hidden'
                }}
            >
                <div>
                    <div
                        style={{
                            width: "100%",
                            margin: "auto",
                            background: "#fff",
                            padding: "5px 5px 5px 20px",
                            display: "flex",
                        }}>
                        <div>
                            <img alt="AIFlowy" src={logo} style={{ height: "38px" }} />
                        </div>

                        {!collapsed &&
                            <div style={{
                                fontSize: "28px",
                                height: "38px",
                                lineHeight: "38px",
                                marginLeft: "10px",
                                fontWeight: "bold",
                            }}>
                                AIFlowy
                            </div>
                        }
                    </div>
                </div>

                <div
                    style={{
                        marginTop: "1px",
                        background: "#fff",
                        flex: "1 1 auto",
                        overflowY: 'auto',
                        overflowX: 'hidden',
                        height: 'calc(100vh - 50px)',
                        scrollbarWidth: 'thin', // Firefox支持
                    }}
                    onScroll={(e) => {
                        e.stopPropagation();
                    }}
                >
                    {!loading ? (
                        <Menu
                            mode="inline"
                            defaultSelectedKeys={selectMenuKeys}
                            defaultOpenKeys={selectMenuKeys}
                            items={menuItems}
                            onClick={(item) => {
                                navigate(item.key)
                            }}
                            style={{
                                borderRight: 'none',
                                // 确保菜单项不会太靠近滚动条
                                paddingRight: '4px'
                            }}
                        />
                    ) : <>loading...</>}
                </div>
            </div>
        </Sider>
    );
};

export default LeftMenu;