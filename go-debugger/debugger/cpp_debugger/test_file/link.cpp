#include <iostream>

// 定义链表节点结构
class Node {
public:
    int data;
    Node* next;

    Node(int value) : data(value), next(nullptr) {}
};

int main() {
    // 创建三个节点
    Node* node1 = new Node(1);
    Node* node2 = new Node(2);
    Node* node3 = new Node(3);

    // 连接节点形成链表
    node1->next = node2;
    node2->next = node3;
    node3->next = nullptr;

    // 打印链表
    Node* current = node1;
    while (current != nullptr) {
        std::cout << current->data << " ";
        current = current->next;
    }
    std::cout << std::endl;

    // 清理内存
    delete node3;
    delete node2;
    delete node1;

    return 0;
} 