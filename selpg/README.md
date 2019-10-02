# 使用golang开发Linux命令行使用程序selpg

## 文章目录
[TOC]

## 概述

> CLI（Command Line Interface）实用程序是Linux下应用开发的基础。正确的编写命令行程序让应用与操作系统融为一体，通过shell或script使得应用获得最大的灵活性与开发效率。Linux提供了cat、ls、copy等命令与操作系统交互；go语言提供一组实用程序完成从编码、编译、库管理、产品发布全过程支持；容器服务如docker、k8s提供了大量实用程序支撑云服务的开发、部署、监控、访问等管理任务；git、npm等都是大家比较熟悉的工具。尽管操作系统与应用系统服务可视化、图形化，但在开发领域，CLI在编程、调试、运维、管理中提供了图形化程序不可替代的灵活性与效率



## 原版C语言selpg分析

### 如何使用

> selpg是从文本输入选择页范围的实用程序。该输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入。selpg 首先处理所有的命令行参数。在扫描了所有的选项参数（也就是那些以连字符为前缀的参数）后，如果 selpg 发现还有一个参数，则它会接受该参数为输入文件的名称并尝试打开它以进行读取。如果没有其它参数，则 selpg 假定输入来自标准输入。

在命令行中输入的指令为如下形式，：

~~~
selpg -sstart_page -eend_page [-lline | -f ] [-d dstFile] filename
~~~

我们来看每个参数的含义。

首先是两个必选参数：

- -sstartPage：例如`-s10`，表示从第10页开始打印。
- -eendPage：例如`-e50`，表示打印到50页。

```
selpg -s10 -e50 ..
```

然后是可选参数：

- -llinePerPage ：例如`-l72`，表示每一页72行，即将输入文件按72行一页分割，默认为72。
- -f ：按照换页符来打印，默认一行一页。
- -ddstFile：输出文件路径，默认为标准输出，即命令行界面。

如果在检查完上面的参数后发现还有一个参数，即filename，则将其设置为输入文件，否则等待标准输入即命令行输入文件路径。

一些例子：

```
selpg -s1 -e1 input_file
```

将input_file的第1页打印到屏幕。

```
selpg -s10 -e20 input_file >output_file
```

将input_file的第10页到第20页打印output_file。

### 源代码分析

首先定义了一个结构体selpg_args，包含了所有参数。

~~~c
struct selpg_args
{
	int start_page;
	int end_page;
	char in_filename[BUFSIZ];
	int page_len; /* default value, can be overriden by "-l number" on command line */
	int page_type; /* 'l' for lines-delimited, 'f' for form-feed-delimited */
					/* default is 'l' */
	char print_dest[BUFSIZ];
};
~~~

process_args函数解析用户的输入，将对应的参数放入参数结构体sa中。以空格分隔的字符串前带有连字符作为标志判断参数，最后如果还有一个参数的话将其设为输入文件，并且检查所有的输入是否符合要求打印错误信息。

process_input函数负责根据sa的属性设定输入输出源，设置每页的行数并开始打印。

usage函数用于打印帮助信息。

那么我们这次的任务，就是将c语言程序“翻译”成go语言程序。

## golang实现selpg

- 首先是selpgArgs结构体，比较简单

   ~~~go
  type selpgArgs struct {
  	startPage int
  	endPage int
  	inputFile string
  	pageLen int
  	pageType bool
  	printDest string
  }
  ~~~

- ProcessArgs函数：

  ~~~go
  func ProcessArgs(parg *selpgArgs) {
  	/*参数个数必须大于等于3*/
  	if len(os.Args) < 3 {
  		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
  		Usage();
  		os.Exit(1)
  	} 
  	/*第一个参数必须是-s*/
  	tmp := os.Args[1]
  	if tmp[0:2] != "-s" {
  		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page\n",  progname)
  		Usage();
  		os.Exit(2)
  	}
  	/*第二个参数必须是-e*/
  	tmp = os.Args[2]
  	if tmp[0:2] != "-e" {
  		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -eend_page\n",  progname)
  		Usage();
  		os.Exit(3)
  	}
  	/*使用pflag对各参数进行绑定*/
  	pflag.IntVarP(&parg.startPage, "start", "s", -1, "start page.")
  	pflag.IntVarP(&parg.endPage, "end", "e", -1, "end page.")
  	pflag.IntVarP(&parg.pageLen, "line", "l", 72, "lines per page.")
  	pflag.BoolVarP(&parg.pageType, "printdes", "f", false, "flag splits page")
  	pflag.StringVarP(&parg.printDest, "destination", "d", "", "name of printer")
  	pflag.Parse()
  	
  	/*处理最后一个参数，即输入文件*/
  	argsLeft := pflag.Args()
  	if len(argsLeft) > 0 {
  		parg.inputFile = string(argsLeft[0])
  	} else {
  		parg.inputFile = ""
  	}
  	/*检查各参数的数值范围*/
  	if parg.pageLen < 1 {
  		fmt.Fprintf(os.Stderr, "%s: invalid page length %s\n", progname, parg.pageLen)
  		os.Exit(4)
  		Usage();
  	}
      if parg.startPage < 1 {
  		fmt.Fprintf(os.Stderr, "%s: invalid start page %s\n", progname, parg.startPage)
  		os.Exit(5)
  		Usage();
  	} 
  	if (parg.startPage > parg.endPage) || (parg.endPage < 1) {
  		fmt.Fprintf(os.Stderr, "%s: invalid end page %s\n", progname, parg.endPage)
  		os.Exit(6)
  		Usage();
  	}
  }
  ~~~

  根据项目要求第二点：

  > 请使用 pflag 替代 goflag 以满足 Unix 命令行规范， 参考：[Golang之使用Flag和Pflag](https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/)

  这里使用了os.Args和pflag包对参数进行解析和检查。pflag包相比于os.Args可以更容易的获取参数对其赋值，也不必像c语言版进行复杂的循环。在对所有参数都绑定完毕之后使用`pflag.parse()`将所有参数的值都保存到结构体`parg`中供其他函数使用。

- ProcessInput函数：

  该函数从标准输入或文件中获取输入然后输出到标准输出或文件中，也可以使用管道将输出作为子进程的输入。
  
  首先设置输入，若inputFile为空则设为标准输入，否则设为对应文件。
  
  ~~~go
  /*设置输入*/
  if parg.inputFile == "" {
      fin = os.Stdin
  } else {
      var finError error
      fin, finError = os.Open(parg.inputFile)
      if finError != nil {
          fmt.Fprintf(os.Stderr, "%s: could not open input file \"%s\"\n", progname, parg.inputFile);
          os.Exit(7)
      }
  }
  ~~~
  
  然后设置输出，这里是本次项目最难的地方。主要是-d参数的实现，需要使用os/exec包来创建子进程执行外部命令，将输出的数据作为子进程的输入。本次使用了该包的部分功能：
  
  - 使用`exec.Command`设定要执行的外部命令
  - 使用`cmd.StdinPipe()`获取连接到子进程标准输入的管道
  - 使用`cmd.Start()`命令开始非阻塞执行子进程
  
  由于没有连接打印机，所以使用cat命令测试。这里由于对管道和linux命令不熟悉，一开始的做法是用exec.Command()再执行了一遍selpg，将-d的参数重定向，虽然效果是一样的。后来看了某师兄的博客，学到了使用cat命令，将父进程selpg的输出设为cat的输入，cat的输出设为-d参数的目标文件，从而实现输出到文件中，不过目标文件需要先创建。
  
  ~~~go
  /*设置输出*/
  if parg.printDest == "" {
      fout = os.Stdout;
  } else {
      /*由于没有打印机，使用cat命令测试*/
      cmd = exec.Command("cat")
      var err error
      cmd.Stdout, err = os.OpenFile(parg.printDest, os.O_WRONLY|os.O_TRUNC, 0600)
      if err != nil {
          fmt.Fprintf(os.Stderr, "%s: file not exist\n%s\n", progname, parg.printDest);
          os.Exit(8)
      }
      fout, _ = cmd.StdinPipe()
      cmd.Start()
  }	
  ~~~
  
  最后按照设定的行数（默认为10）打印。
  
  ~~~go
  for {
      line, err := buffer.ReadString('\n')
      if err == io.EOF {
          break
      }
      if err != nil {
          fmt.Fprintf(os.Stderr, "%s: could not read file\n%s\n", progname, parg.inputFile);
          os.Exit(8)
      }
      lineCount++
      if lineCount > parg.pageLen {
          pageCount++
          lineCount = 1
      }
      if (pageCount >= parg.startPage) && (pageCount <= parg.endPage) { 
          _, err := fout.Write([]byte(line))
          if err != nil {
              fmt.Fprintf(os.Stderr, "%s: could not read file\n%v\n", progname, err);
              os.Exit(9)
          }
      }
  }
  ~~~
  
## 总结

本次作业在难度上比前两次要大很多。本质上来说是实现了一个“翻译”的任务，实现了golang版的selpg。大二操作系统也有写过C语言的命令行程序，所以对CLI也不是一窍不通，不过这一次对管道，重定向有了更深的了解，也简单使用了golang的命令行参数处理包flag和pflag，执行外部指令的exec，还有os和io的一些操作。过程还是比较坎坷，虽然网上参考的博客很多，但是鱼龙混杂，其中很多都有错的或者理解不到位的地方，自己真正会了才可能完整地做出来，否则即使正确的结果也只是简单的摆拍而已。

github地址 -> [🚪]()

项目中附测试截图

## 参考资料

[go语言子进程](https://shimo.im/sheets/KJXKY9YgchdKJXGx/MODOC/)

[go语言创建读取写入文件](https://studygolang.com/articles/14170)

[某博客](https://blog.csdn.net/C486C/article/details/82990187)

[某另一博客](https://blog.csdn.net/hcm_0079/article/details/82927996)

[开发Linux命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)

[selpg源代码](https://www.ibm.com/developerworks/cn/linux/shell/clutil/selpg.c)

