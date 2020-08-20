## 知识整理

### Go语言基础

#### 基本数据结构
##### slice
1. 数据结构

   编译期间的切片是 Slice 类型的，但是在运行时切片由如下的 SliceHeader 结构体表示，一个三元组结构。
   其中 Data 字段是指向数组的指针，Len 表示当前切片的长度，而 Cap 表示当前切片的容量，也就是 Data 数组的大小。
```
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```
   Data 作为一个指针指向的数组是一片连续的内存空间，这片内存空间可以用于存储切片中保存的全部元素，
   数组中的元素只是逻辑上的概念，底层存储其实都是连续的，所以我们可以将切片理解成一片连续的内存空间加上长度与容量的标识。
   ![slice内存结构](./images/slice-1.jpg)

2. 追加和扩容

   在分配内存空间之前需要先确定新的切片容量，Go 语言根据切片的当前容量选择不同的策略进行扩容：
   - 如果期望容量大于当前容量的两倍就会使用期望容量；
   - 如果当前切片的长度小于 1024 就会将容量翻倍；
   - 如果当前切片的长度大于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量；

3. 拷贝切片

   当我们使用 copy(a, b) 的形式对切片进行拷贝时，编译期间的 cmd/compile/internal/gc.copyany 函数也会分两种情况进行处理，
   如果当前 copy 不是在运行时调用的，copy(a, b) 会被直接转换成下面的代码：
```
n := len(a)
if n > len(b) {
    n = len(b)
}
if a.ptr != b.ptr {
    memmove(a.ptr, b.ptr, n*sizeof(elem(a))) 
}
```
   其中 memmove 会负责对内存进行拷贝，在其他情况下，编译器会使用 runtime.slicecopy 函数替换运行期间调用的 copy，例如：go copy(a, b)：
```
func slicecopy(to, fm slice, width uintptr) int {
	if fm.len == 0 || to.len == 0 {
		return 0
	}
	n := fm.len
	if to.len < n {
		n = to.len
	}
	if width == 0 {
		return n
	}
	...

	size := uintptr(n) * width
	if size == 1 {
		*(*byte)(to.array) = *(*byte)(fm.array)
	} else {
		memmove(to.array, fm.array, size)
	}
	return n
}
```
   上述函数的实现非常直接，两种不同的拷贝方式一般都会通过 memmove 将整块内存中的内容拷贝到目标的内存区域中.
   ![slice-copy](./images/slice-2.jpg)
   相比于依次对元素进行拷贝，这种方式能够提供更好的性能，但是需要注意的是，哪怕使用 memmove 对内存成块进行拷贝，
   但是这个操作还是会占用非常多的资源，在大切片上执行拷贝操作时一定要注意性能影响。

4. 小结

   切片的很多功能都是在运行时实现的了，无论是初始化切片，还是对切片进行追加或扩容都需要运行时的支持，
   需要注意的是在遇到大切片扩容或者复制时可能会发生大规模的内存拷贝，一定要在使用时减少这种情况的发生避免对程序的性能造成影响。


[参考](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-array-and-slice/)

##### map

1. 哈希表原理
   对key进行哈希，得到一个地址，直接在该地址存放value。
   实现哈希表的关键点在于如何选择哈希函数，哈希函数的选择在很大程度上能够决定哈希表的读写性能。
   
   冲突解决：
   ```
   开放寻址法
   开放寻址法2是一种在哈希表中解决哈希碰撞的方法，这种方法的核心思想是对数组中的元素依次探测和比较以判断目标键值对是否存在于哈希表中，
   如果我们使用开放寻址法来实现哈希表，那么在支撑哈希表的数据结构就是数组，不
   过因为数组的长度有限，存储 (author, draven) 这个键值对时会从如下哈希位置开始往下遍历，
   当我们向当前哈希表写入新的数据时发生了冲突，就会将键值对写入到下一个不为空的位置
   ```
   
   ```
   拉链法(链地址法)
   与开放地址法相比，拉链法是哈希表中最常见的实现方法，大多数的编程语言都用拉链法实现哈希表，
   它的实现比较开放地址法稍微复杂一些，但是平均查找的长度也比较短，各个用于存储节点的内存都是动态申请的，可以节省比较多的存储空间。

   实现拉链法一般会使用数组加上链表，不过有一些语言会在拉链法的哈希中引入红黑树以优化性能，
   拉链法会使用链表数组作为哈希底层的数据结构，我们可以将它看成一个可以扩展的『二维数组』
   
   当我们需要将一个键值对 (Key6, Value6) 写入哈希表时，键值对中的键 Key6 都会先经过一个哈希函数，哈希函数返回的哈希会帮助我们选择一个桶，
   和开放地址法一样，选择桶的方式就是直接对哈希返回的结果取模，
   选择了 2 号桶之后就可以遍历当前桶中的链表了，在遍历链表的过程中会遇到以下两种情况：
   
   找到键相同的键值对 —— 更新键对应的值；
   没有找到键相同的键值对 —— 在链表的末尾追加新键值对
   ```
   

2. 数据结构
   
   ```
   type hmap struct {
   	count     int
   	flags     uint8
   	B         uint8
   	noverflow uint16
   	hash0     uint32
   
   	buckets    unsafe.Pointer
   	oldbuckets unsafe.Pointer
   	nevacuate  uintptr
   
   	extra *mapextra
   }
   
   // count 表示当前哈希表中的元素数量；
   // B 表示当前哈希表持有的 buckets 数量，但是因为哈希表中桶的数量都 2 的倍数，所以该字段会存储对数，也就是 len(buckets) == 2^B；
   // hash0 是哈希的种子，它能为哈希函数的结果引入随机性，这个值在创建哈希表时确定，并在调用哈希函数时作为参数传入；
   // oldbuckets 是哈希在扩容时用于保存之前 buckets 的字段，它的大小是当前 buckets 的一半；
   
   
   // 桶的结构体 bmap 在 Go 语言源代码中的定义只包含一个简单的 tophash 字段，
   // tophash 存储了键的哈希的高 8 位，通过比较不同键的哈希的高 8 位可以减少访问键值对次数以提高性能
   type bmap struct {
   	tophash [bucketCnt]uint8
   }
   
   
   // 编译后结构体
   // bmap 结构体其实不止包含 tophash 字段，由于哈希表中可能存储不同类型的键值对并且 Go 语言也不支持泛型，
   // 所以键值对占据的内存空间大小只能在编译时进行推导，这些字段在运行时也都是通过计算内存地址的方式直接访问的，
   // 所以它的定义中就没有包含这些字段，但是我们能根据编译期间的 cmd/compile/internal/gc.bmap 函数对它的结构重建
   type bmap struct {
       topbits  [8]uint8
       keys     [8]keytype
       values   [8]valuetype
       pad      uintptr
       overflow uintptr
   }
   
   ```
   ![map内存结构](./images/map-1.jpg)
   
3. 读写操作


4. 扩容

##### context
##### channel
##### select
##### timer

#### 调度器

#### 系统监控 sysmon

#### 内存模型

#### 垃圾回收机制

### 延伸阅读
[参考](https://draveness.me/golang/docs)



----
### 常用后端组件

#### mysql

#### redis

#### mq
##### rabbitMq
##### kafka

#### 分布式
##### etcd
##### consul
##### raft
##### zk



### 开发运维

#### linux
##### 服务器各项性能指标
##### 常用命令

#### docker

#### k8s
 

### 项目经验

#### 通用网关项目
##### 微服务框架
##### grpc
##### 服务注册、发现
##### 服务限流
##### 负载均衡
##### 服务降级
##### 服务监控
##### 链路追踪

#### 消息网关
##### 长连接
##### websocket
##### 消息队列 延时消息
##### 模块解耦
##### 模块解耦

#### 交易所
##### 应用特性
##### 业内常用解决方案
##### 有状态分布式应用
##### raft协议具体实现


