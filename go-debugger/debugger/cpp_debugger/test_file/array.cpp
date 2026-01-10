#include <iostream>
#include <array>
#include <vector>
#include <string>

enum Color { RED, GREEN, BLUE };

struct Item {
    int id;
    float weight;
    Color color;
};

// 全局数组变量
int globalIntArr[3] = {10, 20, 30};
float globalFloatArr[2] = {1.1f, 2.2f};
char globalCharArr[5] = {'x', 'y', 'z', 'w', '\0'};
Color globalColorArr[3] = {RED, GREEN, BLUE};
Item globalItemArr[2] = { {1, 1.1f, RED}, {2, 2.2f, GREEN} };
int* globalPtrArr[2] = {&globalIntArr[0], &globalIntArr[1]};
const int globalConstArr[2] = {111, 222};
static double staticGlobalDoubleArr[2] = {3.14, 2.71};
std::array<int, 2> globalStdArr = {100, 200};
std::vector<float> globalVec = {3.3f, 4.4f, 5.5f};
int global2DArr[2][2] = { {7, 8}, {9, 10} };
const char* globalStrArr[2] = {"hello", "world"};

// 静态全局结构体数组
static Item staticGlobalItemArr[2] = { {3, 3.3f, BLUE}, {4, 4.4f, GREEN} };

int main() {
    // 局部数组变量
    int intArr[4] = {1, 2, 3, 4};
    float floatArr[3] = {1.5f, 2.5f, 3.5f};
    char charArr[6] = "hello";
    Color colorArr[2] = {GREEN, BLUE};
    Item itemArr[2] = { {5, 5.5f, RED}, {6, 6.6f, BLUE} };
    int* ptrArr[2] = {&intArr[0], &intArr[1]};
    const int constArr[2] = {1000, 2000};
    static int staticArr[2] = {7, 8};
    std::array<int, 3> stdArr = {5, 6, 7};
    std::vector<std::string> vecStr = {"foo", "bar"};
    std::vector<Item> vecItem = { {7, 7.7f, GREEN}, {8, 8.8f, RED} };
    int multiArr[2][3] = { {1,2,3}, {4,5,6} };
    char strArr[2][6] = {"hi", "ok"};
    int (*pIntArr)[4] = &intArr; // 指向数组的指针

    // 输出部分内容，防止编译器优化
    std::cout << intArr[0] << floatArr[1] << charArr[2] << colorArr[1]
              << itemArr[0].id << ptrArr[0][0] << constArr[1] << staticArr[1]
              << stdArr[2] << vecStr[1] << vecItem[0].id << multiArr[1][2]
              << strArr[1] << (*pIntArr)[3] << globalIntArr[1] << globalFloatArr[0]
              << globalCharArr[2] << globalColorArr[2] << globalItemArr[1].weight
              << globalPtrArr[1][0] << globalConstArr[1] << staticGlobalDoubleArr[1]
              << globalStdArr[1] << globalVec[2] << global2DArr[1][1]
              << globalStrArr[0] << staticGlobalItemArr[1].id
              << std::endl;

    return 0;
}