package tech.aiflowy.ai.node;


import com.agentsflex.core.chain.ChainNode;
import com.alibaba.fastjson.JSONObject;
import dev.tinyflow.core.Tinyflow;
import dev.tinyflow.core.parser.BaseNodeParser;

/**
 * Sql查询节点解析
 *
 * @author tao
 * @date 2025-05-21
 */
public class SqlNodeParser extends BaseNodeParser {


    @Override
    public ChainNode parse(JSONObject jsonObject, Tinyflow tinyflow) {

        JSONObject data = getData(jsonObject);
        String sql = data.getString("sql");
        SqlNode sqlNode = new SqlNode(sql);
        addParameters(sqlNode, data);
        addOutputDefs(sqlNode, data);
        return sqlNode;
    }

    public String getNodeName() {
        return "sql-node";
    }
}
