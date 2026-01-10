import { languageConstants } from '@/constants/languages';

const defaultCCode = `#include <stdio.h>
int main() {
    int a, b;
	scanf("%d", &a);
    printf("a * a = %d\n", a * a);
    return 0;
}`;

const defaultGoCode = `package main
import (
	"fmt"
)
func main() {
	a := 0
	b := 0
	for {
		n, _ := fmt.Scan(&a, &b)
		if n == 0 {
			break
		} else {
			fmt.Printf("%d\n", a+b)
		}
	}
}
`;

const defaultJavaCode = `import java.util.Scanner;
public class Main {
    public static void main(String[] args) {
        Scanner in = new Scanner(System.in);
        while (in.hasNextInt()) {
            int a = in.nextInt();
            int b = in.nextInt();
            System.out.println(a + b);
        }
    }
}`;

export const defaultCodeMap = new Map<languageConstants, string>([
  [languageConstants.C, defaultCCode],
  [languageConstants.GO, defaultGoCode],
  [languageConstants.Java, defaultJavaCode],
]);
