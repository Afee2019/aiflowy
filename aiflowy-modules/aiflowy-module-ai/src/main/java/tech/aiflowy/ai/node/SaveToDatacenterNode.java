package tech.aiflowy.ai.node;

import com.agentsflex.core.chain.Chain;
import com.agentsflex.core.chain.node.BaseNode;
import com.alibaba.fastjson2.JSONArray;
import com.alibaba.fastjson2.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.aiflowy.common.constant.Constants;
import tech.aiflowy.common.entity.LoginAccount;
import tech.aiflowy.common.util.SpringContextUtil;
import tech.aiflowy.datacenter.service.DatacenterTableService;

import java.math.BigInteger;
import java.util.HashMap;
import java.util.Map;

public class SaveToDatacenterNode extends BaseNode {

    private static final Logger log = LoggerFactory.getLogger(SaveToDatacenterNode.class);

    private BigInteger tableId;

    public SaveToDatacenterNode() {
    }

    public SaveToDatacenterNode(BigInteger tableId) {
        this.tableId = tableId;
    }

    @Override
    protected Map<String, Object> execute(Chain chain) {

        Map<String, Object> map = chain.getParameterValues(this);
        JSONObject json = new JSONObject(map);

        Map<String, Object> res = new HashMap<>();

        // 默认为未知来源
        LoginAccount account = defaultAccount();
        Object cache = chain.getMemory().get(Constants.LOGIN_USER_KEY);
        if (cache != null) {
            account = (LoginAccount) cache;
        }

        DatacenterTableService service = SpringContextUtil.getBean(DatacenterTableService.class);

        JSONArray saveList = json.getJSONArray("saveList");

        int successRows = 0;
        for (Object object : saveList) {
            JSONObject obj = new JSONObject((com.alibaba.fastjson.JSONObject) object);
            obj.put("table_id", tableId);
            try {
                service.saveValue(tableId, obj, account);
            } catch (Exception e) {
                log.error("工作流保存数据到数据中枢失败，表ID：{}，具体值：{}", tableId, obj, e);
                throw e;
            }
            successRows++;
        }

        res.put("successRows", successRows);
        return res;
    }

    public BigInteger getTableId() {
        return tableId;
    }

    public void setTableId(BigInteger tableId) {
        this.tableId = tableId;
    }

    private LoginAccount defaultAccount() {
        LoginAccount account = new LoginAccount();
        account.setId(new BigInteger("0"));
        account.setDeptId(new BigInteger("0"));
        account.setTenantId(new BigInteger("0"));
        return account;
    }
}
