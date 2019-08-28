package stacks

type Example []string

//IntStack stack of integers
type IntStack []int

//Push pushes an int into the satck
func (stack *IntStack) Push(value int) {
	*stack = append(*stack, value)
}

//Pop pops an int fromt he stack
func (stack *IntStack) Pop() int {
	var fmt int
	fmt, *stack = (*stack)[len(*stack)-1], (*stack)[:len(*stack)-1]
	return fmt
}

//Empty true if stack is empty
func (stack *IntStack) Empty() bool {
	return len(*stack) == 0
}

//TopIs checks that the top of stack format is of certain value
func (stack *IntStack) TopIs(value int) bool {
	if len(*stack) == 0 {
		return false
	}
	return (*stack)[len(*stack)-1] == value
}
