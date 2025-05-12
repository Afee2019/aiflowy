package tech.aiflowy.common.ai;

import com.agentsflex.core.document.Document;
import com.agentsflex.core.document.DocumentParser;
import com.alibaba.excel.EasyExcel;
import com.alibaba.excel.read.listener.PageReadListener;

import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class ExcelDocumentParser implements DocumentParser {
    @Override
    public Document parse(InputStream inputStream) {
        List<List<String>> tableData = new ArrayList<>();
        List<String> headerRow = new ArrayList<>();

        // 使用 EasyExcel 读取 Excel 输入流
        EasyExcel.read(inputStream, new PageReadListener<Map<Integer, String>>(dataList -> {
                    for (Map<Integer, String> row : dataList) {
                        List<String> rowData = new ArrayList<>();
                        for (int i = 0; i < row.size(); i++) {  // 注意这里是 i < row.size()
                            rowData.add(row.getOrDefault(i, ""));  // 防止 null 值
                        }

                        if (headerRow.isEmpty()) {
                            // 第一行作为表头
                            headerRow.addAll(rowData);
                            tableData.add(headerRow);
                        } else {
                            // 添加数据行
                            tableData.add(rowData);
                        }
                    }
                }))
                .headRowNumber(0)  // 关键：不要跳过任何行
                .sheet()           // 默认第一个 sheet
                .doRead();
        String plainText = generateMarkdownTable(tableData);

        // 创建并返回 Document 对象
        return new Document(plainText);
    }


    private static String generateMarkdownTable(List<List<String>> tableData) {
        if (tableData == null || tableData.isEmpty()) {
            return "表格数据为空";
        }

        StringBuilder sb = new StringBuilder();

        // 表头
        List<String> headers = tableData.get(0);
        sb.append("| ").append(String.join(" | ", headers)).append(" |\n");

        // 分隔线
        sb.append("|");
        for (int i = 0; i < headers.size(); i++) {
            sb.append(" --- |");
        }
        sb.append("\n");

        // 数据行
        for (int i = 1; i < tableData.size(); i++) {
            List<String> row = tableData.get(i);
            sb.append("| ").append(String.join(" | ", row)).append(" |\n");
        }

        return sb.toString();
    }
}
