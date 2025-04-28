export default {
    'make-file': {
        title: '生成文件',
        group: 'base',
        description: '生成文件',
        icon: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M21 15.2426V21.0082C21 21.556 20.5551 22 20.0066 22H3.9934C3.44476 22 3 21.5511 3 20.9925V9H9C9.55228 9 10 8.55228 10 8V2H20.0017C20.5531 2 21 2.45531 21 2.9918V6.75736L12.0012 15.7562L11.995 19.995L16.2414 20.0012L21 15.2426ZM21.7782 8.80761L23.1924 10.2218L15.4142 18L13.9979 17.9979L14 16.5858L21.7782 8.80761ZM3 7L8 2.00318V7H3Z"></path></svg>',
        sortNo: 802,
        parametersAddEnable: true,
        outputDefsAddEnable: true,
        forms: [{
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
                required: true,
                parametersAddEnable: false,
                description: '生成后的文件地址',
            },
        ]
    }
}