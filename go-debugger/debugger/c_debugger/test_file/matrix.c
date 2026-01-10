#include <stdio.h>

// 打印二维数组
void print2DArray(int arr[][3], int rows) {
    for (int i = 0; i < rows; i++) {
        for (int j = 0; j < 3; j++) {
            printf("%d ", arr[i][j]);
        }
        printf("\n"); // 每行结束后换行
    }
}

int main() {
    // 1. 直接初始化二维数组
    int matrix[2][3] = {
        {1, 2, 3},
        {4, 5, 6}
    };

    print2DArray(matrix, 2);
    return 0;
}
