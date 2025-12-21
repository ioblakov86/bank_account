package main

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"golang.org/x/term"
)

const (
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func green(s string) string {
	if !isTerminal() {
		return s
	}
	return colorGreen + s + colorReset
}

func red(s string) string {
	if !isTerminal() {
		return s
	}
	return colorRed + s + colorReset
}

func yellow(s string) string {
	if !isTerminal() {
		return s
	}
	return colorYellow + s + colorReset
}

type Account struct {
	Owner   string
	Balance float64
	mu      sync.Mutex
}

var ErrInsufficientFunds = errors.New("недостаточно средств")

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("сумма должна быть больше нуля")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("сумма должна быть больше нуля")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if amount > a.Balance {
		return fmt.Errorf("%w для снятия %.2f", ErrInsufficientFunds, amount)
	}

	a.Balance -= amount
	return nil
}

func (a *Account) GetBalance() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Balance
}

func perform(op string, acc *Account, err error) {
	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			fmt.Println(yellow(fmt.Sprint(op, " — ошибка:", err, "| Баланс:", acc.GetBalance())))
			return
		}
		fmt.Println(red(fmt.Sprint(op, " — ошибка:", err)))
		return
	}
	fmt.Println(green(fmt.Sprintf(" %s выполнено. Баланс: %.2f", op, acc.GetBalance())))
}

func NewAccount(owner string) (*Account, error) {
	if owner == "" {
		return nil, fmt.Errorf("owner cannot be empty")
	}
	return &Account{Owner: owner}, nil
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(5)

	account, err := NewAccount("Иван")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Счёт для %s создан. Баланс: %.2f\n", account.Owner, account.GetBalance())

	go func() {
		defer wg.Done()
		perform("Пополнение", account, account.Deposit(10000))
	}()

	go func() {
		defer wg.Done()
		perform("Пополнение", account, account.Deposit(-1))
	}()

	go func() {
		defer wg.Done()
		perform("Снятие", account, account.Withdraw(5000))
	}()

	go func() {
		defer wg.Done()
		perform("Снятие", account, account.Withdraw(-10))
	}()

	go func() {
		defer wg.Done()
		perform("Снятие", account, account.Withdraw(10000))
	}()
	wg.Wait()

	fmt.Printf("Итоговый баланс: %.2f\n", account.GetBalance())
}
