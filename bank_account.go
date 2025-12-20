package main

import (
	"errors"
	"fmt"
)

type Account struct {
	Owner   string
	Balance float64
}

var ErrInsufficientFunds = errors.New("недостаточно средств")

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("сумма должна быть больше нуля")
	}
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("сумма должна быть больше нуля")
	}
	if amount > a.Balance {
		return fmt.Errorf("%w для снятия %.2f", ErrInsufficientFunds, amount)
	}
	a.Balance -= amount
	return nil
}

func (a Account) GetBalance() float64 {
	return a.Balance
}

func perform(op string, acc *Account, err error) {
	cGreen := "\033[32m"
	cRed := "\033[31m"
	cYellow := "\033[33m"
	cReset := "\033[0m"

	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			fmt.Println(cYellow, op, "— ошибка:", err, "| Баланс:", acc.GetBalance(), cReset)
			return
		}
		fmt.Println(cRed, op, "— ошибка:", err, cReset)
		return
	}
	fmt.Printf(" %s%s выполнено. Баланс: %.2f%s\n", cGreen, op, acc.GetBalance(), cReset)
}

func NewAccount(owner string) (*Account, error) {
	if owner == "" {
		return nil, fmt.Errorf("owner cannot be empty")
	}
	return &Account{Owner: owner}, nil
}

func main() {
	account, err := NewAccount("Иван")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Счёт для %s создан. Баланс: %.2f\n", account.Owner, account.GetBalance())

	perform("Пополнение", account, account.Deposit(10000))
	perform("Пополнение", account, account.Deposit(-1))
	perform("Снятие", account, account.Withdraw(5000))
	perform("Снятие", account, account.Withdraw(-10))
	perform("Снятие", account, account.Withdraw(10000))

	fmt.Printf("Итоговый баланс: %.2f\n", account.GetBalance())
}
