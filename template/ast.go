package template

import (
	"fmt"
	"text/template/parse"
)

func printer(n parse.Node) bool {
	if x, ok := n.(*parse.TextNode); ok {
		fmt.Printf("%v", x.String())
	} else if x, ok := n.(*parse.StringNode); ok {
		fmt.Printf("%v", x.String())
	} else if x, ok := n.(*parse.NumberNode); ok {
		fmt.Printf("%v", x.String())
	} else if x, ok := n.(*parse.BoolNode); ok {
		fmt.Printf("%v", x.String())
	} else if x, ok := n.(*parse.ActionNode); ok {
		fmt.Print("{{")
		visit(x.Pipe, printer)
		fmt.Print("}}")
		return false
	} else if _, ok := n.(*parse.CommandNode); ok {
	} else if x, ok := n.(*parse.PipeNode); ok {
		if len(x.Decl) > 0 {
			for i, a := range x.Decl {
				visit(a, printer)
				if i < len(x.Decl)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Print(" := ")
		}
		for _, a := range x.Cmds {
			visit(a, printer)
		}
		return false
	} else if x, ok := n.(*parse.VariableNode); ok {
		fmt.Printf("%v ", x.String())
	} else if x, ok := n.(*parse.FieldNode); ok {
		fmt.Printf("%v ", x.String())
	} else if x, ok := n.(*parse.IdentifierNode); ok {
		fmt.Printf("%v ", x.String())
	} else if _, ok := n.(*parse.ListNode); ok {
	} else if x, ok := n.(*parse.RangeNode); ok {
		fmt.Print("{{range ")
		visit(x.Pipe, printer)
		fmt.Print("}}")
		visit(x.List, printer)
		if x.ElseList != nil {
			fmt.Print("{{else}}")
			visit(x.ElseList, printer)
		}
		fmt.Print("{{end}}")
		return false
	} else {
		fmt.Printf("%T\n", n)
	}
	return true
}

func visit(n parse.Node, fn func(parse.Node) bool) bool {
	if n == nil {
		return true
	}
	if !fn(n) {
		return false
	}
	switch l := n.(type) {
	case *parse.ListNode:
		for _, nn := range l.Nodes {
			if !visit(nn, fn) {
				continue
			}
		}

	case *parse.RangeNode:
		visit(l.Pipe, fn)
		if l.List != nil {
			visit(l.List, fn)
		}
		if l.ElseList != nil {
			visit(l.ElseList, fn)
		}

	case *parse.ActionNode:
		for _, c := range l.Pipe.Decl {
			visit(c, fn)
		}
		for _, c := range l.Pipe.Cmds {
			if visit(c, fn) {
				for _, a := range c.Args {
					visit(a, fn)
				}
			}
		}

	case *parse.CommandNode:
		for _, a := range l.Args {
			visit(a, fn)
		}
	}
	return true
}


func listTemplateNodes(node parse.Node, res []string) []string {
   find:= func(n parse.Node)bool {
    if n, ok := n.(*parse.TemplateNode); ok {
        res = append(res, n.Name)
    }
    return true
   } 
   visit(node, find)
  return res
}
