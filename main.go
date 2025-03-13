package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
	"golang.org/x/term"
)

var sudoPassword string
var sudoEnabled bool

func main() {
	// --sudo flag'ini tekshiramiz
	if len(os.Args) > 1 && os.Args[1] == "--sudo" {
		fmt.Print("[sudo] parolni kiriting: ")

		// Terminaldan yashirin parol olish
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // Yangi qator qo'shamiz, chunki ReadPassword qatorni chiqarmaydi
		if err != nil {
			fmt.Println("Xatolik: Parolni o‘qib bo‘lmadi.")
			os.Exit(1)
		}
		sudoPassword = strings.TrimSpace(string(bytePassword))
		sudoEnabled = true

		// **Parolni tekshirish**
		if !checkSudoPassword(sudoPassword) {
			fmt.Println("Xatolik: Noto‘g‘ri parol!")
			os.Exit(1)
		}

		fmt.Println("Sudo rejimi faollashtirildi!")
	}

	// **Readline sozlamalarini yaratamiz (Tab completion bilan)**
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/shell-history.tmp",                           // Tarixni saqlash
		AutoComplete:    readline.NewPrefixCompleter(getShellCommands()...), // Tab completion
		InterruptPrompt: "^C",
	})
	if err != nil {
		fmt.Println("Xatolik: Terminalni o'qib bo‘lmadi.")
		os.Exit(1)
	}
	defer rl.Close()

	for {
		// Buyruqni o'qish (↑ va ↓ ham ishlaydi)
		input, err := rl.Readline()
		if err != nil { // Agar foydalanuvchi CTRL+D bossa, chiqib ketamiz
			fmt.Println("\nDasturdan chiqildi.")
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Agar chiqish buyruqlari bo'lsa, dasturdan chiqamiz
		if input == "exit" || input == "quit" {
			fmt.Println("Dasturdan chiqildi.")
			break
		}

		// Buyruqni bo'laklarga ajratamiz
		args := strings.Fields(input)

		// Sudo tekshirish
		if args[0] == "sudo" {
			if !sudoEnabled {
				fmt.Println("Sudo rejimi yoqilmagan. --sudo bilan dastur ishga tushiring yoki parolni kiriting.")
				continue
			}

			// `sudo` ni olib tashlab, haqiqiy buyruqni olamiz
			args = args[1:]

			// Parolni `echo` orqali yuboramiz va `-S` flag'idan foydalanamiz
			cmd := exec.Command("bash", "-c", fmt.Sprintf("echo %s | sudo -S %s", sudoPassword, strings.Join(args, " ")))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println("Xatolik:", err)
			}
			continue
		}

		// Oddiy buyruqni ishga tushirish
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("Xatolik:", err)
		}
	}
}

// **Parolni tekshirish funksiyasi**
func checkSudoPassword(password string) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo %s | sudo -S -v", password))
	err := cmd.Run()
	return err == nil
}

// **Tab completion uchun mavjud buyruqlarni olish**
func getShellCommands() []readline.PrefixCompleterInterface {
	// Unix tizimdagi barcha mavjud buyruqlarni olish uchun `compgen -c` ishlatamiz
	out, err := exec.Command("bash", "-c", "compgen -c").Output()
	if err != nil {
		fmt.Println("Xatolik: Buyruqlarni olishda muammo bo'ldi.")
		return nil
	}

	// Buyruqlar ro‘yxatini olish
	commands := strings.Split(string(out), "\n")
	var completions []readline.PrefixCompleterInterface

	// Har bir buyruqni completion sifatida qo‘shamiz
	for _, cmd := range commands {
		if cmd != "" {
			completions = append(completions, readline.PcItem(cmd))
		}
	}

	// Ichki buyruqlar (custom commands) qo‘shish
	extraCommands := []string{"exit", "quit", "clear"}
	for _, cmd := range extraCommands {
		completions = append(completions, readline.PcItem(cmd))
	}

	return completions
}
