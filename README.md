# ZhaJinHuaGameOfGolang
扎金花Demo。目的：golang语言学习实战。。

简单的一下吧，扎金花Demo。目的：golang语言学习实战。。

首先本来客户端和服务器的代码在一起的，都是用C#编写。

这里用Golang重写了服务器端，但是有前提：
1.客户端那边逻辑等修改不了。
2.连传输的协议都无法修改。在客户端，协议的传输的消息体导的是dll。
3.本来客户端和服务器的代码在一起的，都是用C#编写，是一个完美的成品，这里只是Golang重写了服务器端，所以客户端那边没人改。
4.客户端Net（网络部分）序列化部分修改了。

问题：
已经定好的协议，因为是C#写的，他们的传输序列化使用的C#自带的序列化方式，
所以这里用golang写，序列化需要自己重写（客户端的部分也是），但是协议你无法改（也就是需求是死的）。

进度：
目前已经实现了登陆和房间匹配功能。客户端和服务器的连接，可以正常单机游戏。联网部分实现了房间匹配功能，
战斗模块还没实现。 因为战斗模块的传输消息体太复杂了！之前用c#没有语言传输和序列化，现在用golang实现起来比较麻烦。

收获和不足：
1.熟悉了golang语言，整体Demo开发，c#也了解了下。
2.自定义序列化，比较麻烦，跨语言，比如c#之间，可以直接List，MAP封装成一个类然后序列后传递。但是设计跨语言c#和
golang就不行了，（这里又不能用json和protubof等，因为客户端那边是死的改不了，客户端的小伙伴比较忙也没空写Demo）。
自己序列化，感觉还是封装、抽象度不够高，设计的不够好。
3.golang的设计模式，上升到设计模式了，之前搞c++，设计模式面向对象的思想。但是go语言这个。。总觉得开始在设计上怪怪的，
有种是“不是这样做设计太简单”的感觉。这个还要慢慢学习。