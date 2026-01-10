#include <stdio.h>
#include <stdlib.h>

// 定义链表节点结构
typedef struct Node {
    int data;
    struct Node* next;
} Node;

// 主函数，程序入口
int main() {
    // 创建三个节点
    Node* node1 = (Node*)malloc(sizeof(Node));
    Node* node2 = (Node*)malloc(sizeof(Node));
    Node* node3 = (Node*)malloc(sizeof(Node));

    // 为节点赋值
    node1->data = 1;
    node2->data = 2;
    node3->data = 3;

    // 连接节点形成链表
    node1->next = node2;
    node2->next = node3;
    node3->next = NULL;

    // 打印链表
    Node* current = node1;
    while (current != NULL) {
        printf("%d ", current->data);
        current = current->next;
    }
    printf("\n");

    // 释放链表内存
    free(node1);
    free(node2);
    free(node3);

    return 0;
}