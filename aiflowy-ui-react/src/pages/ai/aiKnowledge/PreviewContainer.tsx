// components/PreviewContainer.tsx
import React from 'react';
import { List, Button } from 'antd';
import { FileTextOutlined } from '@ant-design/icons';
import '../style/PreviewContainer.less'; // 引入样式文件

interface PreviewItem {
    sorting: string;
    content: string;
    score: string;
}

interface PreviewContainerProps {
    data?: PreviewItem[];
    total?: number;
    loading?: boolean;
    confirmImport?: boolean;
    disabledConfirm?: boolean;
    onCancel?: () => void;
    onConfirm?: () => void;
    isSearching?: boolean;
}

const PreviewContainer: React.FC<PreviewContainerProps> = ({
                                                               total = 0,
                                                               data = [],
                                                               loading = false,
                                                               confirmImport = false,
                                                               disabledConfirm = false,
                                                               onCancel = () => {},
                                                               onConfirm = () => {},
                                                               isSearching = false,
                                                           }) => {
    return (
        <div className="preview-container">
            <div className="preview-header">
                <h3>
                    <FileTextOutlined />
                    {isSearching ? '检索结果' : '文档预览'}
                </h3>
                {data.length > 0 && (
                    <span className="preview-stats">
            共 {total} 个分段
          </span>
                )}
            </div>

            <div className="preview-content">
                    <List
                        className="preview-list"
                        itemLayout="horizontal"
                        dataSource={data}
                        loading={loading}
                        renderItem={(item: PreviewItem, index) => (
                            <List.Item>
                                <List.Item.Meta
                                    avatar={<div className="segment-badge">{item.sorting ?? index + 1}</div>}
                                    title={
                                        isSearching ? (
                                            <a>{`相似度 ${item?.score ?? 'N/A'}`}</a>
                                        ) : undefined
                                    }
                                    description={<div style={{ whiteSpace:'normal' }}>
                                        {item.content}
                                    </div>}
                                />
                            </List.Item>
                        )}
                    />
            </div>

            {confirmImport && (
                <div className="preview-actions">
                    <div className="action-buttons">
                        <Button
                            style={{
                                minWidth: '100px',
                                height: '36px'
                            }}
                            onClick={onCancel}
                        >
                            取消导入
                        </Button>
                        <Button
                            type="primary"
                            style={{
                                minWidth: '100px',
                                height: '36px'
                            }}
                            loading={disabledConfirm}
                            onClick={onConfirm}
                        >
                            确认导入
                        </Button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default PreviewContainer;
