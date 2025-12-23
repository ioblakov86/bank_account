package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
)

type CommandType int

const (
	CmdDeposit CommandType = iota
	CmdWithdraw
	CmdGetBalance
)

type Command struct {
	Type   CommandType
	Amount float64
	Reply  chan Result
}

type Result struct {
	Balance float64
	Err     error
}

var ErrInsufficientFunds = errors.New("недостаточно средств")

type Account struct {
	Owner   string
	Balance float64
}

func AccountActor(owner string) chan<- Command {
	cmds := make(chan Command)

	go func() {
		balance := 0.0

		for cmd := range cmds {
			switch cmd.Type {

			case CmdDeposit:
				if cmd.Amount <= 0 {
					cmd.Reply <- Result{balance, errors.New("amount must be > 0")}
					continue
				}
				balance += cmd.Amount
				cmd.Reply <- Result{balance, nil}

			case CmdWithdraw:
				if cmd.Amount <= 0 {
					cmd.Reply <- Result{balance, errors.New("amount must be > 0")}
					continue
				}
				if cmd.Amount > balance {
					cmd.Reply <- Result{balance, ErrInsufficientFunds}
					continue
				}
				balance -= cmd.Amount
				cmd.Reply <- Result{balance, nil}

			case CmdGetBalance:
				cmd.Reply <- Result{balance, nil}
			}
		}
	}()

	return cmds
}

func main() {
	owner := flag.String("owner", "", "Владелец счёта")
	op := flag.String("op", "", "Операция: deposit | withdraw | balance")
	repeat := flag.Int("repeat", 1, "Количество повторений")
	amount := flag.Float64("amount", 0, "Сумма операции")
	flag.Parse()

	if *owner == "" {
		fmt.Println("owner is required")
		os.Exit(1)
	}

	account := AccountActor(*owner)

	wg := sync.WaitGroup{}
	wg.Add(*repeat)

	for i := 0; i < *repeat; i++ {
		go func() {
			defer wg.Done()

			reply := make(chan Result)

			switch *op {
			case "deposit":
				account <- Command{CmdDeposit, *amount, reply}
			case "withdraw":
				account <- Command{CmdWithdraw, *amount, reply}
			case "balance":
				account <- Command{CmdGetBalance, 0, reply}
			default:
				fmt.Println("unknown operation")
				return
			}

			res := <-reply

			if res.Err != nil {
				fmt.Println("Error:", res.Err)
				return
			}

			fmt.Println("Balance:", res.Balance)
		}()
	}

	wg.Wait()

}
