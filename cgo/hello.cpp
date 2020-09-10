// 表示该函数的链接符号遵循C语言的规则。
extern "C" {
    #include "hello.h"
}
#include <iostream>

void SaySomething(const char* s) {
    std::cout << s;
    std::cout << "\n";
}