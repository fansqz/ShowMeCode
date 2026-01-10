#include <stdio.h>

// 二分查找函数
int binarySearch(int arr[], int left, int right, int target) {
    while (left <= right) {
        int mid = left + (right - left) / 2;

        // 检查中间元素是否为目标元素
        if (arr[mid] == target)
            return mid;

        // 如果中间元素小于目标元素，则在右半部分继续查找
        else if (arr[mid] < target)
            left = mid + 1;

        // 如果中间元素大于目标元素，则在左半部分继续查找
        else
            right = mid - 1;
    }

    // 目标元素不存在于数组中
    return -1;
}

int main() {
    int arr[] = {2, 4, 6, 8, 10, 12, 14, 16, 18, 20};
    int n = sizeof(arr) / sizeof(arr[0]);
    int target = 12;

    int result = binarySearch(arr, 0, n - 1, target);

    if (result != -1)
        printf("元素 %d 在数组中的索引是 %d\n", target, result);
    else
        printf("元素 %d 不在数组中\n", target);

    return 0;
}