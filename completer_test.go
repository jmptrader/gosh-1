package gosh

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mySimpleCommand struct{}

func (c *mySimpleCommand) SubCommands() CommandMap {
	return nil
}

func (c *mySimpleCommand) Exec(arguments []Argument) error {
	return nil
}

func NewSimpleCommand() *mySimpleCommand {
	return &mySimpleCommand{}
}

type myComplexCommand struct {
	*mySimpleCommand
	subCommands CommandMap
}

func (c *myComplexCommand) SubCommands() CommandMap {
	return c.subCommands
}

func (c *myComplexCommand) Exec(arguments []Argument) error {
	return nil
}

func NewComplexCommand(subCommands CommandMap) *myComplexCommand {
	return &myComplexCommand{NewSimpleCommand(), subCommands}
}

var _ = Describe("Gosh", func() {
	Describe("Completer behavior", func() {
		var completer *Completer
		BeforeEach(func() {
			completer = NewCompleter(CommandMap{
				"john":  NewSimpleCommand(),
				"james": NewSimpleCommand(),
				"mary":  NewSimpleCommand(),
				"nancy": NewSimpleCommand(),
			})
		})

		It("Should return all the top level strings when the empty string is supplied", func() {
			wanted := []string{"james", "john", "mary", "nancy"}
			Expect(completer.Complete("")).To(Equal(wanted))
		})

		It("Should return strings that match the input prefix", func() {
			wanted := []string{"james", "john"}
			Expect(completer.Complete("j")).To(Equal(wanted))
		})
	})

	Describe("Second level HierarchyCompleter response", func() {
		var completer *Completer
		BeforeEach(func() {
			completer = NewCompleter(CommandMap{
				"john": NewComplexCommand(CommandMap{
					"jacob":        NewSimpleCommand(),
					"jingleheimer": NewSimpleCommand(),
					"schmidt":      NewSimpleCommand(),
				}),
				"james": NewSimpleCommand(),
				"mary":  NewSimpleCommand(),
				"nancy": NewSimpleCommand(),
			})
		})

		It("Should return all the second level tokens when there is an exact match for the first field and no second field", func() {
			wanted := []string{"john jacob", "john jingleheimer", "john schmidt"}
			Expect(completer.Complete("john ")).To(Equal(wanted))
		})

		It("Should return only matching second level tokens when there is an exact match for the first field and second field", func() {
			wanted := []string{"john jacob", "john jingleheimer"}
			Expect(completer.Complete("john j")).To(Equal(wanted))
		})
	})
})
