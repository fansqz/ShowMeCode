# 可视化

<div style="border: 2px solid #FF9800; border-radius: 5px; background-color: #FFF3E0; padding: 15px; margin-top: 20px;">   <i style="color: #FF9800; font-size: 24px; margin-bottom: 10px;" class="fas fa-exclamation-circle"></i><p style="color: #444;">&#9997;ShowMeCode的调试支持运行时数据结构的可视化。想要该功能需要一定配置，下面是关于可视化配置的说明。</p> </div>

## 数组可视化

我们在ShowMeCode中写入如下代码（`C`语言）：

```go
#include <stdio.h>

int main() {
    int arr[] = {64, 34, 25, 12, 22, 11, 90};
    int n = sizeof(arr) / sizeof(arr[0]);
    for (int i = 0; i < n; i++) {
        printf("%d ", arr[i]);
    }
    printf("\n");
    return 0;
}    
```

我们使用了`arr`这个数组，并且数组使用了`i`、`j`这些变量进行取值。我们进行以下配置（多个取值变量之间使用`,`分割）：

![Description](https://fan-code-image-1308043430.cos.ap-guangzhou.myqcloud.com/problem/file/9c7df3fd-113b-11f0-8df6-00155d03b0f2.image.png)

打开可视化开关
![Description](https://fan-code-image-1308043430.cos.ap-guangzhou.myqcloud.com/problem/file/dda9b38b-2779-11f0-901a-00163e06661f.image.png)
添加断点
![Description](https://fan-code-image-1308043430.cos.ap-guangzhou.myqcloud.com/problem/file/fe0ce02b-2779-11f0-901a-00163e06661f.image.png)
开始调试！！！
![Description](https://fan-code-image-1308043430.cos.ap-guangzhou.myqcloud.com/problem/file/2559c1e4-277a-11f0-901a-00163e06661f.image.png)


&nbsp;

## 链表可视化

我们在ShowMeCode中定义如下结构体（`C`语言）：

```go
typedef struct ListNode {
    int val;
    struct ListNode *prev;
    struct ListNode *next;
} ListNode;
```

我们定义了`Node`结构体作为链表节点的结构体，其中`prev`指针指向链表当前节点上一个节点，`next`指向当前节点的下一个节点，那么我们进行以下配置：
![Description](https://fan-code-image-1308043430.cos.ap-guangzhou.myqcloud.com/problem/file/a6bf7ebf-116f-11f0-8df6-00155d03b0f2.image.png)

添加断点，并开始调试！！！

![Description](https://fan-code-image-1308043430.cos.ap-guangzhou.myqcloud.com/problem/file/c6d7bb8f-116f-11f0-8df6-00155d03b0f2.image.png)

<div style="border-left: 2px solid #FF9800; border-radius: 0 5px 5px 0; background-color: #F5F5F5; padding: 15px; margin-top: 20px;">   <i style="color: #FF9800; font-size: 24px; margin-bottom: 10px;" class="fas fa-exclamation-circle"></i>&#128075;&#128075;&#128075;<p style="color: #444;">以上完成了可视化学习，其他语言的配置同理，和C语言的配置方式是一样的</p> </div>