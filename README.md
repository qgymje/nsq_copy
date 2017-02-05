# learning by copying source code.


### 2017/2/5 Sun

##### dirlock 

1. 知识点涉及到c语言的flock函数:
```int flock(int fd, int operation)```
operation是系统定义的整型常量:
LOCK_SH: share,共享锁,其实程序可以访问这个这个文件
LOCK_EX: exclusive,排它锁，独享锁, 互斥锁, 其它程序无法访问此文件
LOCK_UN: unlock
LOCK_NB: 当无法建立锁时,操作马上返回而不会等待,通常与LOCK_EX做OR操作组合

[参考地址](http://blog.csdn.net/loophome/article/details/49681217)

2. atomic.Value 类型的
```Store(x interface{})```
```Load() interface{}```

[stupid gopher tricks](https://www.youtube.com/watch?v=UECh7X07m6E&t=22m40s)
提到适合存储global configurations
