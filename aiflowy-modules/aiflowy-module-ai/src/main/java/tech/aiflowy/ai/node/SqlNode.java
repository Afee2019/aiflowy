package tech.aiflowy.ai.node;
import com.agentsflex.core.chain.Chain;
import com.agentsflex.core.chain.node.BaseNode;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jfinal.template.stat.ast.For;
import com.mybatisflex.core.row.Db;
import com.mybatisflex.core.row.Row;
import net.sf.jsqlparser.JSQLParserException;
import net.sf.jsqlparser.parser.CCJSqlParserUtil;
import net.sf.jsqlparser.statement.Statement;
import net.sf.jsqlparser.statement.select.Select;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.util.StringUtils;
import tech.aiflowy.common.web.exceptions.BusinessException;

import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;


/**
 * SQL查询节点
 *
 * @author tao
 * @date 2025-05-21
 */
public class SqlNode extends BaseNode {

    private final String sql;

    private static final Logger logger = LoggerFactory.getLogger(SqlNode.class);

    public SqlNode(String sql) {
        this.sql = sql;
    }

    @Override
    protected Map<String, Object> execute(Chain chain) {

        Map<String, Object> map = chain.getParameterValues(this);
        Map<String, Object> res = new HashMap<>();
        Map<String, Object> formatSqlMap = formatSql(sql);
        String formatSql = (String)formatSqlMap.get("replacedSql");

        Statement statement = null;
        try {
            statement  = CCJSqlParserUtil.parse(formatSql);

        } catch (JSQLParserException e) {
            logger.error("sql 解析报错：",e);
            throw new BusinessException("SQL解析失败，请确认SQL语法无误");
        }

        if (!(statement instanceof Select)) {
            logger.error("sql 解析报错：statement instanceof Select 结果为false");
            throw new BusinessException("仅支持查询语句！");
        }

        List<String> paramNames = (List<String>) formatSqlMap.get("paramNames");

        List<Object> paramValues = new ArrayList<>();
        paramNames.forEach(paramName -> {
            Object o = map.get(paramName);
            paramValues.add(o);
        });

        List<Row> rows = Db.selectListBySql(formatSql, paramValues.toArray());

        res.put("queryData", rows);
        return res;
    }

    private Map<String,Object> formatSql(String sql) {

        if (!StringUtils.hasLength(sql)){
            logger.error("sql解析报错：sql为空");
            throw new BusinessException("sql 不能为空！");
        }

        // 用来提取参数名
        Pattern pattern = Pattern.compile("\\{\\{([^}]+)}}");
        Matcher matcher = pattern.matcher(sql);

        List<String> paramNames = new ArrayList<>();

        // 构建替换后的 SQL
        StringBuffer replacedSql = new StringBuffer();
        while (matcher.find()) {
            paramNames.add(matcher.group(1)); // 获取 {{...}} 中的内容
            matcher.appendReplacement(replacedSql, "?");
        }

        matcher.appendTail(replacedSql);
        HashMap<String, Object> formatSqlMap = new HashMap<>();
        String formatSql = replacedSql.toString();

        if (formatSql.endsWith(";") || formatSql.endsWith("；")) {
            formatSql = formatSql.substring(0, formatSql.length() - 1);
        }

        formatSql =  formatSql.replace("“" ,"\"").replace("”","\"");

        logger.info("Replaced SQL: {}", replacedSql);
        logger.info("Parameter names: {}", paramNames);
        formatSqlMap.put("replacedSql",formatSql );
        formatSqlMap.put("paramNames", paramNames);
        return   formatSqlMap;
    }


}
