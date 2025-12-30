# 编译问题分析：Java 版本冲突

> 文档编号：002
> 问题发生日期：2025-12-29
> 状态：已解决
> 更新：2025-12-29 - 修正根因分析

---

## 1. 问题现象

执行 Maven 编译命令时失败：

```bash
$ mvn clean package -DskipTests
```

**错误信息：**
```
[ERROR] Failed to execute goal org.apache.maven.plugins:maven-compiler-plugin:3.13.0:compile
(default-compile) on project aiflowy-common-base:
Fatal error compiling: 无效的目标发行版: 17
```

**编译日志关键信息：**
```
[INFO] Compiling 49 source files with javac [debug target 17] to target/classes
```

---

## 2. 环境信息

| 组件 | 版本 |
|------|------|
| JDK | 1.8.0_421 (Oracle) |
| Maven | 3.9.9 |
| OS | macOS Darwin 25.2.0 (aarch64) |
| maven-compiler-plugin | 3.13.0 |

**项目配置（pom.xml）：**
```xml
<properties>
    <java.version>8</java.version>
    <maven.compiler.source>8</maven.compiler.source>
    <maven.compiler.target>8</maven.compiler.target>
</properties>
```

---

## 3. 问题分析

### 3.1 现象矛盾点

- 本机 JDK 版本是 1.8
- 项目 pom.xml 明确声明 Java 8
- 但编译时却尝试使用 Java 17 作为目标版本

### 3.2 初步诊断（误判）

最初怀疑是 `mybatis-parent:48` 通过依赖链污染了编译器属性。

**依赖链分析：**
```
aiflowy-common-base
    └── mybatis-flex-core:1.11.3
            └── mybatis:3.5.19
                    └── (parent) mybatis-parent:48  ← 定义了 java.version=17
```

但经过进一步测试，发现即使是**完全空的项目**（无任何依赖），effective-pom 中也显示 Java 17。

### 3.3 真正的根因

检查 Maven 全局配置文件：

```bash
$ grep -A15 "jdk-17" /Users/shawn/bin/apache-maven-3.9.9/conf/settings.xml
```

发现：

```xml
<profile>
    <id>jdk-17</id>
    <activation>
        <activeByDefault>true</activeByDefault>  <!-- 问题根源！ -->
        <jdk>17</jdk>
    </activation>
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
        <maven.compiler.compilerVersion>17</maven.compiler.compilerVersion>
    </properties>
</profile>
```

**问题根源：** `activeByDefault=true` 导致该 profile 无条件激活，覆盖了所有项目的编译器配置。

### 3.4 验证测试

| 测试场景 | effective-pom 中的 maven.compiler.source |
|----------|------------------------------------------|
| 空项目（无依赖） | 17 ❌ |
| 依赖 guava（无 mybatis） | 17 ❌ |
| 依赖 mybatis | 17 ❌ |

所有场景都显示 Java 17，证明问题与依赖无关，是 Maven 全局配置导致。

---

## 4. 解决方案

### 4.1 根本修复：修改 Maven settings.xml

**文件：** `/Users/shawn/bin/apache-maven-3.9.9/conf/settings.xml`

**修改内容：** 注释掉 `activeByDefault`

```xml
<profile>
    <id>jdk-17</id>
    <activation>
        <!-- 移除 activeByDefault，避免覆盖所有项目的编译器设置 -->
        <!-- <activeByDefault>true</activeByDefault> -->
        <jdk>17</jdk>  <!-- 仅当实际使用 JDK 17 时激活 -->
    </activation>
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
        <maven.compiler.compilerVersion>17</maven.compiler.compilerVersion>
    </properties>
</profile>
```

### 4.2 项目级修复（可选保留）

在根 `pom.xml` 中显式配置 `maven-compiler-plugin`：

```xml
<build>
    <plugins>
        <plugin>
            <groupId>org.apache.maven.plugins</groupId>
            <artifactId>maven-compiler-plugin</artifactId>
            <version>3.13.0</version>
            <configuration>
                <source>8</source>
                <target>8</target>
                <encoding>UTF-8</encoding>
            </configuration>
        </plugin>
    </plugins>
</build>
```

此配置作为防御性措施保留，可防止类似的全局配置问题。

---

## 5. 验证结果

### 5.1 修复 settings.xml 后

```bash
$ mvn help:effective-pom | grep maven.compiler.source
    <maven.compiler.source>8</maven.compiler.source>
```

### 5.2 编译测试

```bash
$ mvn clean package -DskipTests
[INFO] Compiling 49 source files with javac [debug target 8] to target/classes
...
[INFO] BUILD SUCCESS
```

---

## 6. 问题定性

| 问题 | 答案 |
|------|------|
| 是否是上游项目（AIFlowy/MyBatis）的 Bug？ | **否** |
| 是否是 Maven 的 Bug？ | **否** |
| 是否是本地环境特有问题？ | **是** |
| 其他开发者会遇到吗？ | **不会**（除非有类似的全局配置） |
| 使用 JDK 17 会有问题吗？ | **不会**（JDK 17 可编译 target=17） |
| 需要给上游创建 Issue/PR 吗？ | **不需要** |

---

## 7. 经验总结

### 7.1 教训

1. **全局配置的影响范围**：Maven settings.xml 中的 profile 配置会影响所有项目，使用 `activeByDefault=true` 需谨慎。

2. **排查顺序**：遇到编译器版本问题时，应按以下顺序排查：
   - ① Maven 全局 settings.xml
   - ② 用户级 ~/.m2/settings.xml
   - ③ 项目 pom.xml
   - ④ 依赖的父 POM

3. **不要过早下结论**：初步分析时误判为 mybatis-parent 的问题，但通过空项目测试验证后发现真正原因。

### 7.2 诊断命令速查

```bash
# 查看 effective-pom 中的编译器配置
mvn help:effective-pom | grep maven.compiler

# 查看激活的 profiles
mvn help:active-profiles

# 查看 Maven 全局配置
cat $MAVEN_HOME/conf/settings.xml | grep -A20 "<profiles>"

# 测试空项目（排除依赖影响）
mkdir test && cd test
echo '<project><modelVersion>4.0.0</modelVersion><groupId>t</groupId><artifactId>t</artifactId><version>1</version></project>' > pom.xml
mvn help:effective-pom | grep maven.compiler
```

### 7.3 最佳实践

```xml
<!-- settings.xml 中的 JDK profile 应使用 JDK 版本作为激活条件，而非 activeByDefault -->
<profile>
    <id>jdk-17</id>
    <activation>
        <jdk>17</jdk>  <!-- 仅当 JAVA_HOME 指向 JDK 17 时激活 -->
    </activation>
    <properties>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
    </properties>
</profile>

<profile>
    <id>jdk-8</id>
    <activation>
        <jdk>1.8</jdk>  <!-- 仅当 JAVA_HOME 指向 JDK 8 时激活 -->
    </activation>
    <properties>
        <maven.compiler.source>8</maven.compiler.source>
        <maven.compiler.target>8</maven.compiler.target>
    </properties>
</profile>
```

---

## 8. 修改记录

| 日期 | 修改内容 |
|------|----------|
| 2025-12-29 | 初始版本，误判为 mybatis-parent 问题 |
| 2025-12-29 | 修正根因为 Maven settings.xml 全局配置问题 |

---

## 9. 影响范围

| 影响项 | 说明 |
|--------|------|
| 修改文件 | `/Users/shawn/bin/apache-maven-3.9.9/conf/settings.xml` |
| 影响范围 | 本机所有 Maven 项目 |
| 向后兼容 | 是（JDK 17 项目仍可通过 `<jdk>17</jdk>` 激活） |

---

*文档编号：002 - 编译问题分析：Java 版本冲突*
