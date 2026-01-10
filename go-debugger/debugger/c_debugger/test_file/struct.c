#include <stdio.h>

// 定义日期结构体
struct Date {
    int year;
    int month;
    int day;
};

// 定义学生信息结构体，嵌套 Date 结构体
typedef struct Student {
    char name[50];
    int id;
    struct Date birthdate;
} Student;

// 打印学生信息的函数
void printStudentInfo(struct Student s) {
    printf("学生姓名: %s\n", s.name);
    printf("学生 ID: %d\n", s.id);
    printf("出生日期: %d-%d-%d\n", s.birthdate.year, s.birthdate.month, s.birthdate.day);
}

Student globalstudent = {
   .name = "张三",
   .id = 12345,
   .birthdate = {2005, 3, 15}
};
Student* globalstudentPtr = &globalstudent;

int main() {
    // 初始化一个学生信息结构体变量
    struct Student student = {
       .name = "张三",
       .id = 12345,
       .birthdate = {2005, 3, 15}
    };
    Student* studentPtr = &student;
    printStudentInfo(student);
    printStudentInfo(globalstudent);

    return 0;
}