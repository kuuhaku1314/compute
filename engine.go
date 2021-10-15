package compute

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	LeftBrackets  = "("
	RightBrackets = ")"
	PlusSign      = "+"
	SubSign       = "-"
	MulSign       = "*"
	DivSign       = "/"
)

type OperatorMap map[string]int

func (m OperatorMap) isOperator(s string) bool {
	_, ok := m[s]
	return ok
}

var operatorMap = OperatorMap{
	LeftBrackets: 1,
	RightBrackets: 4,
	PlusSign: 2,
	SubSign: 2,
	MulSign: 3,
	DivSign: 3,
}

type Engine interface {
	Parse(interface{}) error
	Run() (interface{}, error)
	Reset()
}

type ComputeEngine struct {
	Stacker
	postfixExpression []interface{}
	infixExpression   []interface{}
	isParsed          bool
}

func NewComputeEngine() Engine {
	return &ComputeEngine{Stacker: NewStack()}
}

func (ce *ComputeEngine) Parse(in interface{}) error {
	if ce.isParsed {
		return nil
	}
	if err := ce.parseToInfixExpression(in); err != nil {
		return err
	}
	if err := ce.parseToPostfixExpression(); err != nil {
		return err
	}
	ce.isParsed = true
	return nil
}

func (ce *ComputeEngine) Run() (interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	return ce.run()
}

func (ce *ComputeEngine) Reset() {
	if !ce.isParsed {
		return
	}
	ce.Stacker = NewStack()
	ce.infixExpression = nil
	ce.postfixExpression = nil
	ce.isParsed = false
}

func (ce *ComputeEngine) run() (interface{}, error) {
	if !ce.isParsed {
		return nil, errors.New("param is not parsed")
	}
	postfixExpression := ce.postfixExpression
	for _, element := range postfixExpression {
		switch element.(type) {
		case string:
			if value, err := ce.compute(element.(string)); err != nil {
				return nil, err
			} else {
				ce.Push(value)
			}
		case float64:
			ce.Push(element)
		default:
			return nil, errors.New("type is wrong")
		}
	}
	if ce.Size() != 1 {
		return nil, errors.New("param is wrong")
	}
	return ce.Pop(), nil
}

func (ce *ComputeEngine) parseToInfixExpression(in interface{}) error {
	str, ok := in.(string)
	if !ok {
		return errors.New("param is not string type")
	}
	str = prevHandle(str)
	runes := []rune(str)
	var operatorIndex []int
	for i, r := range runes {
		if isOperator(string(r)) {
			operatorIndex = append(operatorIndex, i)
		}
	}
	if len(operatorIndex) == 0 {
		return errors.New("param is wrong because operator is none")
	}
	start := 0
	for _, end := range operatorIndex {
		if end-start != 0 {
			float, err := parseToNumber(runes[start:end])
			if err != nil {
				return err
			}
			ce.infixExpression = append(ce.infixExpression, float)
		}
		ce.infixExpression = append(ce.infixExpression, string(runes[end]))
		start = end + 1
	}
	if start != len(runes) {
		float, err := parseToNumber(runes[start:])
		if err != nil {
			return err
		}
		ce.infixExpression = append(ce.infixExpression, float)
	}
	return nil
}

func (ce *ComputeEngine) parseToPostfixExpression() error {
	infixExpression := ce.infixExpression
	for _, element := range infixExpression {
		switch element.(type) {
		case string:
			if err := ce.handleOperator(element.(string)); err != nil {
				return err
			}
		case float64:
			ce.postfixExpression = append(ce.postfixExpression, element)
		default:
			return errors.New("type is wrong")
		}
	}
	for !ce.IsEmpty() {
		s := ce.Pop().(string)
		if isLeftBrackets(s) {
			return errors.New("param is wrong because Brackets")
		} else {
			ce.postfixExpression = append(ce.postfixExpression, s)
		}
	}
	return nil
}

func (ce *ComputeEngine) compute(s string) (float64, error) {
	if ce.Size() < 2 {
		return 0, errors.New("param is wrong because number < 2")
	}
	right := ce.Pop().(float64)
	left := ce.Pop().(float64)
	switch s {
	case PlusSign:
		return right + left, nil
	case SubSign:
		return left - right, nil
	case MulSign:
		return left * right, nil
	case DivSign:
		return left / right, nil
	}
	return 0, errors.New("param is wrong because operator type is unknown")
}

func (ce *ComputeEngine) handleOperator(s string) error {
	if isLeftBrackets(s) {
		ce.Push(s)
		return nil
	}
	if isRightBrackets(s) {
		for !ce.IsEmpty() {
			s := ce.Pop().(string)
			if isLeftBrackets(s) {
				return nil
			} else {
				ce.postfixExpression = append(ce.postfixExpression, s)
			}
		}
		return errors.New("param is wrong because Brackets")
	}
	if ce.IsEmpty() {
		ce.Push(s)
		return nil
	}
	for !ce.IsEmpty() && ge(ce.Peek().(string), s) {
		ce.postfixExpression = append(ce.postfixExpression, ce.Pop().(string))
	}
	ce.Push(s)
	return nil
}

func ge(s1 string, s2 string) bool {
	return operatorMap[s1] >= operatorMap[s2]
}

func isLeftBrackets(s string) bool {
	return s == LeftBrackets
}

func isRightBrackets(s string) bool {
	return s == RightBrackets
}

func parseToNumber(rs []rune) (float64, error) {
	return strconv.ParseFloat(string(rs), 64)
}

func isOperator(s string) bool {
	return operatorMap.isOperator(s)
}

func prevHandle(s string) string {
	s = strings.Replace(s, " ", "", -1)
	if strings.HasPrefix(s, "-") {
		s = "0" + s
	}
	in := []rune(s)
	builder := &strings.Builder{}
	for i, r := range in {
		if string(r) == "-" && string(in[i - 1]) == "("{
			builder.WriteString("0")
		}
		builder.WriteRune(r)
	}
	return builder.String()
}
