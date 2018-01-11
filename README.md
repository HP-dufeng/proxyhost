# proxyhost

无论你是开发、QA还是产品，在目前的开发及测试中都离不开host。host为开发和测试带来了诸多便利，但也有很多烦恼：

+ 切换host总是不能立即生效，换一套host环境经常需要重启浏览器。
+ 经常在一套host环境下比如（betaA）想对比一下线上的情况，做不到。
+ 由于经常在多套host环境中切换，导致系统host环境乱七八糟，发现问题都不确定到底在哪套环境下。

proxyhost就是为了解决烦人的host问题而诞生的 , 它采用了 **沙箱机制**，在一个独立的浏览器进程中使用host。

+ 正常浏览器不受任何影响，只需配置，可以同时启动多个浏览器，同时查看各个发布环境功能。

### 命令行运行
proxy -url=http://automata.cefcfco.com:6789



## Resources

* [`github.com/liyangready/multiple-host`][9] - 虚拟host解决方案，轻松实现两套host环境