## 限流器及限流算法

    https://juejin.im/post/5d1c978d51882555433429e6
    https://segmentfault.com/a/1190000014745556
    
## golang标准库限流器
    
    https://www.cyhone.com/articles/analisys-of-golang-rate/
    https://www.cyhone.com/articles/usage-of-golang-rate/
    
## redis+lua限流器使用介绍

    该限流器使用的是令牌桶算法
    lua脚本的实现借鉴golang标准库time/rate
    在 Golang 的 timer/rate 中的实现, 并没有单独维护一个 Timer，而是采用了 lazyload 的方式，
    直到每次消费之前才根据时间差更新 Token 数目，而且也不是用 BlockingQueue 来存放 Token，
    而是仅仅通过计数的方式。