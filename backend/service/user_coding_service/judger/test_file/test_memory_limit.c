#include <stdio.h>
#include <stdlib.h>

int main() {
    int num1, num2, sum;
    scanf("%d", &num1);
    scanf("%d", &num2);
    size_t size = 10 * 1024 * 1024;  // 10MB
    void* memory = malloc(size);
    if (memory == NULL) {
        printf("内存申请失败\n");
        return 1;
    }
    free(memory);  // 释放内存
    sum = num1 + num2;
    printf("%d\n",sum);
    return 0;
}
