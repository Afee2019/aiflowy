export default {
    'make-file': {
        title: '文件生成',
        group: 'base',
        description: '生成 Word、PDF、HTML 等文件供用户下载',
        icon: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M14 13.5V8C14 5.79086 12.2091 4 10 4C7.79086 4 6 5.79086 6 8V13.5C6 17.0899 8.91015 20 12.5 20C16.0899 20 19 17.0899 19 13.5V4H21V13.5C21 18.1944 17.1944 22 12.5 22C7.80558 22 4 18.1944 4 13.5V8C4 4.68629 6.68629 2 10 2C13.3137 2 16 4.68629 16 8V13.5C16 15.433 14.433 17 12.5 17C10.567 17 9 15.433 9 13.5V8H11V13.5C11 14.3284 11.6716 15 12.5 15C13.3284 15 14 14.3284 14 13.5Z"></path></svg>',
        sortNo: 802,
        parametersAddEnable: true,
        outputDefsAddEnable: true,
        forms: [
            {
                type: 'heading',
                label: '文件设置',
            },
            {
                type: 'select',
                label: '文件类型',
                description: '请选择生成的文件类型',
                name: 'suffix',
                defaultValue: 'docx',
                options: [
                    {
                        label: 'docx',
                        value: 'docx'
                    }
                ]
            }],
        parameters: [
            {
                name: 'content',
                nameDisabled: true,
                title: '内容',
                dataType: 'String',
                required: true,
                description: '内容',
                deleteDisabled: true,
            },
        ],
        outputDefs: [
            {
                name: 'url',
                nameDisabled: true,
                title: '文件下载地址',
                dataType: 'String',
                dataTypeDisabled: true,
                required: true,
                parametersAddEnable: false,
                description: '生成后的文件地址',
                deleteDisabled: true,
            },
        ]
    }
}