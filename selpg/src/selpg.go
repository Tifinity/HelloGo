package main

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"os/exec"
	"github.com/spf13/pflag"
)

type selpgArgs struct {
	startPage int
	endPage int
	inputFile string
	pageLen int
	pageType bool
	printDest string
}

var progname string

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
	pflag.IntVarP(&parg.pageLen, "line", "l", 10, "lines per page.")
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

func ProcessInput(parg *selpgArgs) {
	var fin *os.File
	var fout io.WriteCloser
	cmd := &exec.Cmd{}
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

	/*开始打印*/
	lineCount := 0
	pageCount := 1
	buffer := bufio.NewReader(fin)
	for {
		line, err := buffer.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not read file\n%s\n", progname, parg.inputFile);
			os.Exit(9)
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
				os.Exit(10)
			}
		}
	}
}

/*打印帮助信息*/
func Usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE: %s -s=start_page -e=end_page [-f | -l=lines_per_page] [-d=dest] [in_filename]\n", progname)
}

func main() {
	progname = os.Args[0]
	var args selpgArgs
	ProcessArgs(&args)
	ProcessInput(&args)
}