package interpreter

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"
)

func (interp *Interpreter) RegisterTemplate(signatures ...TaTemplate) error {
	for i := 0; i < len(signatures); i++ {
		signature := signatures[i]
		signature.sanitize()
		if interp.GetTemplate(&signature) != nil {
			return errors.Errorf("Template `%s' is already registered", signature.Name)
		}
		interp.Templates = append(interp.Templates, signature)
	}
	return nil
}

func (interp *Interpreter) MustRegisterTemplate(signatures ...TaTemplate) {
	if err := interp.RegisterTemplate(signatures...); err != nil {
		panic(err)
	}
}

func (interp *Interpreter) UpdateTemplate(signature TaTemplate) error {
	signature.sanitize()
	if s := interp.GetTemplate(&signature); s != nil {
		*s = signature
		return nil
	}
	return errors.Errorf("Function `%s' is not registered", signature.Name)
}

func (interp *Interpreter) RemoveTemplate(signature TaTemplate) error {
	signature.sanitize()
	for i := 0; i < len(interp.Templates); i++ {
		if interp.Templates[i].Equal(&signature) {
			fns := interp.Templates[:i]
			interp.Templates = append(fns, interp.Templates[i+1:]...)
			return nil
		}
	}
	return errors.Errorf("Function `%s' is not registered", signature.Name)
}

func (interp *Interpreter) GetTemplate(signature *TaTemplate) *TaTemplate {
	signature.sanitize()
	for i := 0; i < len(interp.Templates); i++ {
		if interp.Templates[i].Equal(signature) {
			return &interp.Templates[i]
		}
	}
	return nil
}

func (interp *Interpreter) registerCoreTemplates() error {
	return nil
}

func (interp *Interpreter) RemoveAllTemplates() error {
	interp.Templates = []TaTemplate{}
	return nil
}

var templateSignature = TaFunction{
	CommonSignature: CommonSignature{
		Name:       "!",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.String,
			block.Any,
		},
		Returns:     block.Any,
		Description: "Resolve a template",
		Example: `
(! Template1)                                                    ; executes the Template1
(! Template2 "Hello World")                                      ; executes Template2 with "Hello World" as parameter
`,
	},
	Func: func(interp *Interpreter, args ...*block.TaToken) (*block.TaToken, error) {
		walker := templateWalker{interp: interp}
		blockText := strings.ToLower(args[0].String)
		// iterate trough all functions
		for template := walker.Next(); template != nil; template = walker.Next() {
			run, detail, children, err := interp.matchesSignature(&template.CommonSignature, blockText, args[1:])
			if err != nil {
				return nil, errors.Errorf("error in template `%s': %v", template.Name, err)
			}
			if !run {
				if interp.Logger != nil {
					switch detail {
					case invalidSignature:
						interp.Logger.Printf("NOT Running template `%s' (not matching signature)\n", template.String())
					case errorInChildrenEvaluation:
						interp.Logger.Printf("NOT Running template `%s' (errors in child evaluation)\n", template.String())
					}
				}
				continue
			}
			if interp.Logger != nil {
				interp.Logger.Printf("Running template `%s' with `%v'\n", template.String(), block.BlockArguments(children).ToHumanReadable())
			}
			b := template.Template
			if len(args) > 1 {
				if _, err := replaceVariables(&b, args[1:]...); err != nil {
					return nil, err
				}
			}
			return &b, nil
		}
		return nil, errors.Errorf("template `%s' not found", args[0].String)
	},
}

func replaceVariables(b *block.TaToken, args ...*block.TaToken) (int, error) {
	total := 0

	var replaced int

replace:
	replaced = 0
	for i := 0; i < len(args); i++ {
		replaced += replaceVariable(b, strconv.Itoa(i), args[i])
	}

	total += replaced

	if replaced > 0 {
		goto replace
	}
	return total, nil
}

func replaceVariable(source *block.TaToken, name string, replace *block.TaToken) (replaced int) {
	if len(source.Children) <= 0 {
		return replaced
	}

	if source.String == "#" {
		if strings.EqualFold(source.Children[0].String, name) {
			*source = *replace
			replaced++
		}
	}

	for i := 0; i < len(source.Children); i++ {
		replaced += replaceVariable(source.Children[i], name, replace)
	}

	return replaced
}

// func (interp *Interpreter) SetTemplate(name string, str string) error {
// 	block, err := lexer.Lex(str)
// 	if err != nil {
// 		return err
// 	}
// 	return interp.RegisterTemplate(TaTemplate{
// 		CommonSignature: CommonSignature{
// 			Name: name,
// 		},
// 		Template: *block,
// 	})
// }
