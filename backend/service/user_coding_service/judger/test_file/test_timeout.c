#include <stdio.h>
#include <unistd.h>

int main() {
    int num1, num2, sum;
    scanf("%d", &num1);
    scanf("%d", &num2);
    sleep(2);
    sum = num1 + num2;
    printf("%d\n",sum);

    return 0;
}
