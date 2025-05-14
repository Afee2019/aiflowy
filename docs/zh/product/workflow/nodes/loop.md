# 循环

循环节点指的是按一定次数循环执行某个逻辑，然后将每次循环的执行结果以数组的形式返回。

## 示例
> 图片较大，看不清可右键选择在新标签页中查看。

指定 `loopVar` 参数，可以是固定值，也可以是引用值。

连接循环体的节点会根据指定的 `loopVar` 参数，循环执行。

最终循环节点会将所有结果输出为数组。

循环体当中有两个变量，一个是 `loopItem`，一个是 `index`。

如果循环体的上一个节点也是循环体，则 `loopItem` 会是上一个节点输出的数组中的内容。

![loopConfigDemo.png](../resource/loopConfigDemo.png)

执行结果：

![loopRunDemo.png](../resource/loopRunDemo.png)