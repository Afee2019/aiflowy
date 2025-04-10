# 配置文件

## 数据库配置

```yml
spring:
  datasource:
    url: jdbc:mysql://127.0.0.1:3306/aiflowy?useInformationSchema=true&characterEncoding=utf-8
    username: root
    password: 123456
```

## 本地文件配置
```yml
spring:
  servlet:
    multipart:
      max-file-size: 100MB
      max-request-size: 100MB
  web:
    resources:
      # 此处要和下面的 aiflowy.storage.local.root 一致。
      # 如 root: 'D://files'，那这里就该是 static-locations: file:D://files
      static-locations: classpath:/public
  mvc:
    static-path-pattern: /static/**
aiflowy:
  storage:
    local:
      # 默认存储在classpath下的public目录
      # target/public 下
      root: ''
```

默认存储在本地。

另外，我们也可以去实现自己的存储类型，只需要编写一个类，实现 `FileStorageService` 接口，并通过 `@Component` 注解为当前的实现类型取个名字，例如：

```java
@Component("myStorage")
public class MyFileStorageServiceImpl implements FileStorageService {

    @Override
    public String save(MultipartFile file) {
        // 在这里，去实现你的文件存储逻辑
    }

    @Override
    public InputStream readStream(String path) throws IOException {
        // 在这里，去实现你的文件存储逻辑
    }
    
}
```
此时，我们添加如下配置，即可把当前 APP 的存储类型修改为你自己的实现类：
```yml
aiflowy:
  storage:
    type: myStorage
```

## 其他配置
```yml
aiflowy:
  # ollama 服务地址
  ollama:
    host: http://127.0.0.1:11434
  # 不进行登录拦截的路径
  login:
    excludes: /api/v1/auth/**, /static/**
```