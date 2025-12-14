package main

import (
	"fmt"
)

type Account struct {
	Owner   string
	Balance float64
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("сумма пополнения не должна быть отрицательной")
	}
	a.Balance += amount
	return nil

}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("сумма снятия не должна быть отрицательной")
	}
	if a.Balance < amount {
		return fmt.Errorf("недостаточно средств")
	}
	a.Balance -= amount
	return nil
}

func (a Account) GetBalance() float64 {
	return a.Balance
}

func main() {
	var newAccount Account
	newAccount.Owner = "Иван"
	fmt.Printf("Счет для клиента %s создан, баланс: %.2f\n", newAccount.Owner, newAccount.GetBalance())
	err := newAccount.Deposit(10000)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Баланс пополнен. Баланс: %.2f\n", newAccount.GetBalance())
	}
	err = newAccount.Deposit(-1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Баланс пополнен. Баланс: %.2f\n", newAccount.GetBalance())
	}
	err = newAccount.Deposit(-1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Баланс пополнен. Баланс: %.2f\n", newAccount.GetBalance())
	}
	err = newAccount.Withdraw(5000)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Баланс пополнен. Баланс: %.2f\n", newAccount.GetBalance())
	}
	err = newAccount.Withdraw(-1.222)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Баланс пополнен. Баланс: %.2f\n", newAccount.GetBalance())
	}
	fmt.Printf("Итоговый баланс: %.2f\n", newAccount.GetBalance())
}
