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
		return ErrInsufficientFunds
	}
	a.Balance -= amount
	return nil
}

func (a Account) GetBalance() float64 {
	return a.Balance
}

func perform(op string, err error, balance float64) {
	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			fmt.Printf("%s — ошибка: %w. Баланс: %.2f\n", op, balance)
			return
		}
		fmt.Println(op, "— ошибка:", err)
		return
	}
	fmt.Printf("%s выполнено. Баланс: %.2f\n", op, balance)
}

func main() {
	account := Account{Owner: "Иван"}

	fmt.Printf("Счёт для %s создан. Баланс: %.2f\n", account.Owner, account.GetBalance())

	perform("Пополнение", account.Deposit(10000), account.GetBalance())
	perform("Пополнение", account.Deposit(-1), account.GetBalance())
	perform("Снятие", account.Withdraw(5000), account.GetBalance())
	perform("Снятие", account.Withdraw(-10), account.GetBalance())

	fmt.Printf("Итоговый баланс: %.2f\n", account.GetBalance())
}
