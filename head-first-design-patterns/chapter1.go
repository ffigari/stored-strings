package main

type Animal interface {
	MakeSound() string
}

type Dog struct {}

func (d Dog) MakeSound() string {
	return "wawaw"
}

type Cat struct {}

func (d Cat) MakeSound() string {
	return "meow"
}

//

type flyBehavior interface {
	Fly() string
}

type quackBehavior interface {
	Quack() string
}

type Duck interface {
	PerformFly() string
	PerformQuack() string
	SetFlyBehavior(fb flyBehavior)
}

type duck struct {
	flyBehavior flyBehavior
	quackBehavior quackBehavior
}

func (d *duck) PerformFly() string {
	return d.flyBehavior.Fly()
}

func (d *duck) SetFlyBehavior(fb flyBehavior) {
	d.flyBehavior = fb
}

func (d duck) PerformQuack() string {
	return d.quackBehavior.Quack()
}

func newDuck(fb flyBehavior, qb quackBehavior) duck {
	return duck{
		flyBehavior: fb,
		quackBehavior: qb,
	}
}

type MallardDuck struct {
	duck
	Name string
}

func NewMallardDuck(name string) MallardDuck {
	return MallardDuck{
		duck: newDuck(YesFlyBehavior{}, QuackQuackBehavior{}),
		Name: name,
	}
}

type RubberDuck struct {
	duck
}

func NewRubberDuck() RubberDuck {
	return RubberDuck{
		duck: newDuck(NoFlyBehavior{}, SqueakQuackBehavior{}),
	}
}

type DecoyDuck struct {
	duck
}

func NewDecoyDuck() DecoyDuck {
	return DecoyDuck{
		duck: newDuck(NoFlyBehavior{}, NoQuackBehavior{}),
	}
}

type YesFlyBehavior struct {}

func (b YesFlyBehavior) Fly() string {
	return "flying"
}

type RocketFlyBehavior struct {}

func (b RocketFlyBehavior) Fly() string {
	return "rocket flying"
}

type NoFlyBehavior struct {}

func (b NoFlyBehavior) Fly() string {
	return ""
}

type QuackQuackBehavior struct {}

func (b QuackQuackBehavior) Quack() string {
	return "quack"
}

type SqueakQuackBehavior struct {}

func (b SqueakQuackBehavior) Quack() string {
	return "squeak"
}

type NoQuackBehavior struct {}

func (b NoQuackBehavior) Quack() string {
	return ""
}

//

type DuckCaller struct {
	quackBehavior quackBehavior
}

func NewDuckCaller() DuckCaller {
	return DuckCaller{
		quackBehavior: QuackQuackBehavior{},
	}
}

func (c *DuckCaller) PerformQuack() string {
	return c.quackBehavior.Quack()
}
