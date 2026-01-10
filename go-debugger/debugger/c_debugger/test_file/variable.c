#include <stdio.h>
#include <stdlib.h>
// 定义枚举类型
typedef enum {
   RED,
   GREEN,
   BLUE
} Color;
// 定义结构体类型
typedef struct {
   int id;
   float weight;
   Color color;
} Item;
// 定义联合体类型
typedef union {
   int ival;
   float fval;
   char cval;
} Value;
typedef int * PTR_INT;
// 全局变量
int globalInt = 10;
float globalFloat = 3.14;
char globalChar = 'A';
Item globalItem = {1, 65.5, RED};
Item* globalItemPtr = &globalItem;
// 静态全局变量
static int staticGlobalInt = 20;
// 函数声明
void manipulateLocals();
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
   // 结构体局部变量
   Item localItem;
   localItem.id = 2;
   localItem.weight = 42.0;
   localItem.color = GREEN;
   // 局部枚举变量
   Color localColor = BLUE;
   // 局部联合体变量
   Value localValue;
   localValue.ival = 123;
   // 输出局部变量的值
    printf("localInt: %d, localChar: %c, staticLocalFloat: %.2f\n", localInt, localChar, staticLocalFloat);
    printf("localItem: id=%d, weight=%.1f, color=%d\n", localItem.id, localItem.weight, localItem.color);
    printf("localColor: %d, localValue: %d\n", localColor, localValue.ival);
}                         // 56
void manipulatePointers() { // 57
   // 动态分配的变量   // 58
   PTR_INT dynamicInt = (int*) malloc(sizeof(int));
   *dynamicInt = 30;
   // 指针变量
   int* ptrToInt = &globalInt;
   Item* ptrToItem = &globalItem;
   Color* ptrToColor = (Color*) malloc(sizeof(Color));
   *ptrToColor = BLUE;
   // 数组变量
   int intArray[3] = { 1, 2, 3 };
   float floatArray[] = { 1.1f, 2.2f, 3.3f };
   Color colorArray[3] = { RED, GREEN, BLUE };
   // 字符串
   char* string = "Hello, World!";
	// 空指针
	Item* nilPoit;
	int* intArrayPtr = &intArray;
   // 清理动态分配的内存
   free(dynamicInt);
   free(ptrToColor);
}
