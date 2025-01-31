package main_test

import (
	e "github.com/ffigari/stored-strings/head-first-design-patterns"
)

func (s *Suite) TestPolyphormism() {
	animals := []e.Animal{e.Dog{}, e.Cat{}}

	animalsSounds := make([]string, 2, 2)

	for i, animal := range animals {
		animalsSounds[i] = animal.MakeSound()
	}

	s.Equal("wawaw", animalsSounds[0])
	s.Equal("meow", animalsSounds[1])
}

func (s *Suite) TestDucksBehavior() {
	tom := e.NewMallardDuck("tom")
	s.Equal("flying", tom.PerformFly())
	s.Equal("quack", tom.PerformQuack())
	s.Equal("tom", tom.Name)

	robert := e.NewMallardDuck("robert")
	s.Equal("flying", robert.PerformFly())
	s.Equal("quack", robert.PerformQuack())
	s.Equal("robert", robert.Name)

	rubberDuck := e.NewRubberDuck()
	s.Equal("", rubberDuck.PerformFly())
	s.Equal("squeak", rubberDuck.PerformQuack())

	decoyDuck := e.NewDecoyDuck()
	s.Equal("", decoyDuck.PerformFly())
	s.Equal("", decoyDuck.PerformQuack())

	d1 := e.NewDecoyDuck()
	d2 := e.NewRubberDuck()
	var ds []e.Duck = []e.Duck{&d1, &d2}
	s.Equal(
		[]string{"", "squeak"},
		[]string{ds[0].PerformQuack(), ds[1].PerformQuack()},
	)

	dd := e.NewDecoyDuck()
	var d e.Duck = &dd
	s.Equal("", d.PerformFly())
	d.SetFlyBehavior(e.RocketFlyBehavior{})
	s.Equal("rocket flying", d.PerformFly())
}

func (s *Suite) TestDuckCaller() {
	caller := e.NewDuckCaller()
	s.Equal("quack", caller.PerformQuack())
}
