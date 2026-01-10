#include <iostream>
#include <string>
#include <memory>
#include <array>
#include <vector>
// 定义枚举类型，使用 enum class 增强类型安全性
enum class Color {
    RED,
    GREEN,
    BLUE
};

// 定义结构体类型，C++ 中 struct 可直接使用
struct Item {
    int id;
    float weight;
    Color color;
};

// 定义联合体类型
union Value {
    int ival;
    float fval;
    char cval;
};

// 使用 using 定义类型别名
using PTR_INT = std::unique_ptr<int>;

// 全局变量
int globalInt = 10;
float globalFloat = 3.14;
char globalChar = 'A';
Item globalItem {1, 65.5, Color::RED};
std::unique_ptr<Item> globalItemPtr = std::make_unique<Item>(globalItem);

// 静态全局变量
static int staticGlobalInt = 20;

// 函数声明
void manipulateLocals(int argint);
void manipulatePointers();

int main() {
    manipulateLocals(2);
    manipulatePointers();
    return 0;
}

void manipulateLocals(int argint) {
    // 局部变量
    int localInt = 5;
    char localChar = 'G';
    // 静态局部变量
    static float staticLocalFloat = 6.78;
    // 结构体局部变量，使用统一初始化语法
    Item localItem {2, 42.0, Color::GREEN};
    // 局部枚举变量
    Color localColor = Color::BLUE;
    // 局部联合体变量
    Value localValue;
    localValue.ival = 123;

    // 输出局部变量的值，使用 std::cout 进行输出
    std::cout << "localInt: " << localInt << ", localChar: " << localChar
              << ", staticLocalFloat: " << staticLocalFloat << std::endl;
    std::cout << "localItem: id=" << localItem.id << ", weight=" << localItem.weight
              << ", color=" << static_cast<int>(localItem.color) << std::endl;
    std::cout << "localColor: " << static_cast<int>(localColor)
              << ", localValue: " << localValue.ival << std::endl;
}

void manipulatePointers() {
    // 动态分配的变量，使用 std::unique_ptr 管理内存
    PTR_INT dynamicInt = std::make_unique<int>(30);
    // 指针变量
    int* ptrToInt = &globalInt;
    Item* ptrToItem = globalItemPtr.get();
    auto ptrToColor = std::make_unique<Color>(Color::BLUE);
    // 数组变量
    std::vector<int> intArray {1, 2, 3};
    std::array<float, 3> floatArray {1.1f, 2.2f, 3.3f};
    std::array<Color, 3> colorArray {Color::RED, Color::GREEN, Color::BLUE};
    // 字符串，使用 std::string 类型
    std::string str = "Hello, World!";
    std::cout << "String: " << str << std::endl;

    // 空指针，使用 nullptr
    Item* nilPoit = nullptr;
    int* intArrayPtr = intArray.data();

    // 无需手动清理内存，std::unique_ptr 会自动管理
}