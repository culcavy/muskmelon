package repl

import (
	"bufio"
	"fmt"
	"github.com/hollykbuck/muskmelon/lexer"
	"github.com/hollykbuck/muskmelon/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	for {
		_, err := fmt.Fprintf(out, PROMPT)
		if err != nil {
			return fmt.Errorf("输出失败: %w", err)
		}
		scanned := scanner.Scan()
		if !scanned {
			return err
		}
		line := scanner.Text()
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			_, err = fmt.Fprintf(out, "%+v\n", tok)
			if err != nil {
				return fmt.Errorf("输出失败: %w", err)
			}
		}
	}
}
