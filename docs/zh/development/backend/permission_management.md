# 权限管理

AIFlowy 是基于 RBAC 的权限控制，用户可以精确控制到按钮级别的权限。

权限框架为 [sa-token（v1.40.0）](https://sa-token.cc/)

## crud 权限拦截

基础的 crud 接口被提取到 `tech.aiflowy.common.web.controller.BaseCurdController` 中去了。

其中包含了 `list`、`page`、`detail`、`save`、`update`、`remove`、`removeBatch`、`intelligentFilling` 等接口。

按照功能划分为 `query`、`save`、`remove` 三种权限。新增功能模块时，需要在菜单中配置这几样基础权限，否则会被拦截器拦截并提示 `无权操作`。

统一处理 crud 权限的拦截器为：`tech.aiflowy.auth.config.CurdInterceptor`

代码如下：
```java
public class CurdInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception {

        String groupName = "";
        // 检查handler是否是HandlerMethod类型
        if (handler instanceof HandlerMethod) {
            HandlerMethod handlerMethod = (HandlerMethod) handler;

            // 获取类上的特定注解
            UsePermission classAnnotation = handlerMethod.getBeanType().getAnnotation(UsePermission.class);
            if (classAnnotation != null) {
                // 处理注解逻辑
                groupName = classAnnotation.moduleName();
            }
            // 有此注解，交给 sa token 自行判断
            SaCheckPermission saCheckPermission = handlerMethod.getMethodAnnotation(SaCheckPermission.class);
            if (saCheckPermission != null) {
                return true;
            }
        }
        String requestUri = request.getRequestURI();
        // 查询
        String finalGroupName = groupName;
        SaRouter.match("/**/list",
                "/**/page",
                "/**/detail",
                "/**/intelligentFilling"
        ).check(r -> {
            checkBaseCurd(requestUri, finalGroupName, "query");
        });
        // 保存
        SaRouter.match("/**/save",
                "/**/update"
        ).check(r -> {
            checkBaseCurd(requestUri, finalGroupName, "save");
        });
        // 删除
        SaRouter.match("/**/remove",
                "/**/removeBatch"
        ).check(r -> {
            checkBaseCurd(requestUri, finalGroupName, "remove");
        });

        return true;
    }

    private void checkBaseCurd(String uri, String groupName, String permission) {
        int idx = uri.lastIndexOf("/");
        String per = uri.substring(0,idx + 1) + permission;
        if (StrUtil.isNotEmpty(groupName)) {
            // 如果指定了继承的模块，就改为该模块的权限校验
            per = groupName + "/" + permission;
        }
        StpUtil.checkPermission(per);
    }
}
```

## 权限继承

某些功能模块存在子项，比如定时任务模块，它的子项则是定时任务日志模块。

但是日志模块又不需要出现在菜单中，因此没法赋予用户日志模块的权限，访问接口时会被拦截器拦截。

这时就可以用权限继承 `@UsePermission` 来让用户获得定时任务的权限时同时具备日志的权限。

在定时任务日志的控制器类上加上注解即可，代码如下：
```java
@RestController
@RequestMapping("/api/v1/sysJobLog")
@UsePermission(moduleName = "/api/v1/sysJob") // 继承定时任务的权限
public class SysJobLogController extends BaseCurdController<SysJobLogService, SysJobLog> {
    public SysJobLogController(SysJobLogService service) {
        super(service);
    }

    @Override
    protected Result onSaveOrUpdateBefore(SysJobLog entity, boolean isSave) {
        LoginAccount loginUser = SaTokenUtil.getLoginAccount();
        if (isSave) {
            commonFiled(entity,loginUser.getId(),loginUser.getTenantId(), loginUser.getDeptId());
        }
        return super.onSaveOrUpdateBefore(entity, isSave);
    }
}
```

## 业务接口的权限控制

除了基础的增删改查接口，平时会写很多业务接口，如果需要对这些接口进行权限控制。

直接使用 `sa-token` 的 `@SaCheckPermission` 注解即可。

<span style="color: red;font-weight: bold;">需要注意的是，由于拦截器拦截 crud 时是通过接口后缀来判断的。因此写业务接口时要避免使用这些后缀：
`list`、`page`、`detail`、`save`、`update`、`remove`、`removeBatch`、`intelligentFilling`</span>