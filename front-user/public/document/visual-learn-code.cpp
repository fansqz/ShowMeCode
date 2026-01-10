#include <iostream>
#include <vector>

int main() {
    // 使用vector代替C风格数组，更安全和灵活
    std::vector<int> arr = {64, 34, 25, 12, 22, 11, 90};
    
    // 输出数组元素
    for (int num : arr) {
        std::cout << num << " ";
    }
    std::cout << std::endl;
    
    return 0;
}