package tech.aiflowy.common.ai;

import com.agentsflex.core.document.Document;
import com.agentsflex.core.document.DocumentSplitter;
import com.agentsflex.core.document.id.DocumentIdGenerator;
import com.alibaba.fastjson2.JSON;
import com.alibaba.fastjson2.TypeReference;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class ExcelDocumentSplitter implements DocumentSplitter {
    private int rowsPerChunk;
    private boolean includeHeader;

    public ExcelDocumentSplitter(int rowsPerChunk) {
        if (rowsPerChunk <= 0) {
            throw new IllegalArgumentException("rows must be greater than 0");
        }
        this.rowsPerChunk = rowsPerChunk;
    }

    public ExcelDocumentSplitter(int rowsPerChunk, boolean includeHeader) {
        if (rowsPerChunk <= 0) {
            throw new IllegalArgumentException("rows must be greater than 0");
        }
        this.rowsPerChunk = rowsPerChunk;
        this.includeHeader = includeHeader;
    }

    @Override
    public List<Document> split(Document document, DocumentIdGenerator idGenerator) {
        if (document == null || document.getContent() == null) {
            return Collections.emptyList();
        }

        // 解析JSON数据为表格结构
        List<List<String>> tableData = JSON.parseObject(document.getContent(),
                new TypeReference<List<List<String>>>() {});

        if (tableData == null || tableData.isEmpty()) {
            return Collections.emptyList();
        }

        List<Document> chunks = new ArrayList<>();
        List<String> headers = includeHeader ? tableData.get(0) : null;
        int startRow = includeHeader ? 1 : 0;

        // 按照指定行数分割数据
        for (int i = startRow; i < tableData.size(); i += rowsPerChunk) {
            int endRow = Math.min(i + rowsPerChunk, tableData.size());

            // 构建当前分块的Markdown表格
            StringBuilder sb = new StringBuilder();

            // 添加表头（如果包含）
            if (headers != null) {
                sb.append("| ").append(String.join(" | ", headers)).append(" |\n");
                sb.append("|");
                for (int j = 0; j < headers.size(); j++) {
                    sb.append(" --- |");
                }
                sb.append("\n");
            }

            // 添加数据行
            for (int j = i; j < endRow; j++) {
                List<String> row = tableData.get(j);
                sb.append("| ").append(String.join(" | ", row)).append(" |\n");
            }

            // 创建新文档
            Document newDocument = new Document();
            newDocument.addMetadata(document.getMetadataMap());
            newDocument.setContent(sb.toString());

            if (idGenerator != null) {
                newDocument.setId(idGenerator.generateId(newDocument));
            }

            chunks.add(newDocument);
        }

        return chunks;
    }

    public int getRowsPerChunk() {
        return rowsPerChunk;
    }

    public void setRowsPerChunk(int rowsPerChunk) {
        this.rowsPerChunk = rowsPerChunk;
    }

    public boolean isIncludeHeader() {
        return includeHeader;
    }

    public void setIncludeHeader(boolean includeHeader) {
        this.includeHeader = includeHeader;
    }
}