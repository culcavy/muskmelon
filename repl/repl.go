package repl

import (
	"bufio"
	"fmt"
	"github.com/hollykbuck/muskmelon/evaluator"
	"github.com/hollykbuck/muskmelon/lexer"
	"github.com/hollykbuck/muskmelon/parser"
	"io"
)

const PROMPT = ">> "

const MONKEY_FACE = `
     ______
    /_____/\
   /_____\\ \
  /_____\ \\ / 
 /_____/ \/ / / 
/_____/ /   \//\ 
\_____\//\   / / 
 \_____/ / /\ /
  \_____/ \\ \ 
   \_____\ \\ 
    \_____\/
`

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
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			err = printParserErrors(out, p.Errors())
			if err != nil {
				return err
			}
			continue
		}
		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			_, err = io.WriteString(out, evaluated.Inspect())
			if err != nil {
				return err
			}
			_, err = io.WriteString(out, "\n")
			if err != nil {
				return err
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) error {
	_, err := io.WriteString(out, MONKEY_FACE)
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, " parser errors:\n")
	if err != nil {
		return err
	}
	for _, msg := range errors {
		_, err = io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			return err
		}
	}
	return nil
}
