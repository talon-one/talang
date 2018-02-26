package interpreter

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

func (interp *Interpreter) RegisterTemplate(signature shared.TaTemplate) error {
	signature.Name = strings.ToLower(signature.Name)
	if interp.GetTemplate(signature) != nil {
		return errors.Errorf("Template `%s' is already registered", signature.Name)
	}
	interp.Templates = append(interp.Templates, signature)
	return nil
}

func (interp *Interpreter) UpdateTemplate(signature shared.TaTemplate) error {
	signature.Name = strings.ToLower(signature.Name)
	if s := interp.GetTemplate(signature); s != nil {
		*s = signature
		return nil
	}
	return errors.Errorf("Function `%s' is not registered", signature.Name)
}

func (interp *Interpreter) RemoveTemplate(signature shared.TaTemplate) error {
	signature.Name = strings.ToLower(signature.Name)
	for i := 0; i < len(interp.Templates); i++ {
		if interp.Templates[i].Equal(&signature) {
			fns := interp.Templates[:i]
			interp.Templates = append(fns, interp.Templates[i+1:]...)
			return nil
		}
	}
	return errors.Errorf("Function `%s' is not registered", signature.Name)
}

func (interp *Interpreter) GetTemplate(signature shared.TaTemplate) *shared.TaTemplate {
	signature.Name = strings.ToLower(signature.Name)
	for i := 0; i < len(interp.Templates); i++ {
		if interp.Templates[i].Equal(&signature) {
			return &interp.Templates[i]
		}
	}
	return nil
}

func (interp *Interpreter) registerCoreTemplates() error {
	return nil
}

func (interp *Interpreter) RemoveAllTemplates() error {
	interp.Templates = []shared.TaTemplate{}
	return nil
}

func templateSignature(interp *Interpreter) shared.TaFunction {
	return shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:       "!",
			IsVariadic: true,
			Arguments: []block.Kind{
				block.StringKind,
				block.AnyKind,
			},
			Returns:     block.AnyKind,
			Description: "Resolve a template",
		},
		Func: templateFunc(interp),
	}
}

func templateFunc(interp *Interpreter) func(_ *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
	return func(_ *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 1 {
			return nil, errors.New("invalid or missing arguments")
		}
		templates := interp.AllTemplates()
		blockText := strings.ToLower(args[0].String)
		// iterate trough all functions
		for n := 0; n < len(templates); n++ {
			template := templates[n]
			run, notMatchingDetail, children, err := interp.matchesSignature(&template.CommonSignature, blockText, args[1:])
			if err != nil {
				return nil, errors.Errorf("error in template `%s': %v", template.Name, err)
			}
			if !run {
				if interp.Logger != nil {
					switch notMatchingDetail {
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
			if argc > 1 {
				if _, err := replaceVariables(&b, args[1:]...); err != nil {
					return nil, err
				}
			}
			return &b, nil
		}
		return nil, errors.Errorf("template `%s' not found", args[0].String)
	}
}

func replaceVariables(b *block.Block, args ...*block.Block) (int, error) {
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

func replaceVariable(source *block.Block, name string, replace *block.Block) (replaced int) {
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
// 	return interp.RegisterTemplate(shared.TaTemplate{
// 		CommonSignature: shared.CommonSignature{
// 			Name: name,
// 		},
// 		Template: *block,
// 	})
// }
