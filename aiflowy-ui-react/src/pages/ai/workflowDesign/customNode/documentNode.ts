export default {
    'document-node': {
        title: '获取文档中的内容',
        group: 'base',
        description: '提取PDF中的文字内容',
        icon: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M9 2.00318V2H19.9978C20.5513 2 21 2.45531 21 2.9918V21.0082C21 21.556 20.5551 22 20.0066 22H3.9934C3.44476 22 3 21.5501 3 20.9932V8L9 2.00318ZM5.82918 8H9V4.83086L5.82918 8ZM11 4V9C11 9.55228 10.5523 10 10 10H5V20H19V4H11Z"></path></svg>',
        sortNo: 801,
        parametersAddEnable: false,
        outputDefsAddEnable: false,
        parameters: [
            {
                name: 'fileUrl',
                nameDisabled: true,
                title: '文档地址',
                dataType: 'File',
                required: true,
                description: '文档的url地址',
            },
        ],
        outputDefs: [
            {
                name: 'content',
                nameDisabled: true,
                title: '解析后的文本',
                dataType: 'String',
                required: true,
                parametersAddEnable: false,
                description: '解析后的文本内容',
            },
        ]
    }
}