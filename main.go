package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
	"github.com/creack/pty"
	"golang.org/x/term"
)

var sudoPassword string
var sudoEnabled bool

// Faol jarayonni global saqlash (CTRL+C ushlab qolish uchun)
var activeCmd *exec.Cmd

// **Global signal channel**
var sigs chan os.Signal

func init() {
	// **Bitta signal channel yaratamiz (ko‘p marta chaqirilmasligi uchun)**
	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	// **CTRL+C signalini ushlash (barcha jarayonlar uchun)**
	go func() {
		for range sigs {
			if activeCmd != nil {
				fmt.Println("\n[!] Jarayon to‘xtatildi")
				activeCmd.Process.Signal(syscall.SIGINT) // Faqat faol jarayonni o‘ldiramiz
			}
		}
	}()
}

func main() {
	// **Foydalanuvchining default shell'ini aniqlash**
	userShell := os.Getenv("SHELL")
	if userShell == "" {
		userShell = "/bin/bash" // Agar aniqlanmasa, bash ishlatamiz
		fmt.Println("SHELL aniqlanmadi, bash ishlatiladi")
	}

	// **--sudo flag'ini tekshiramiz**
	if len(os.Args) > 1 && os.Args[1] == "--sudo" {
		fmt.Print("[sudo] parolni kiriting: ")

		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
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

	// **Readline sozlamalari (Tab completion bilan)**
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/shell-history.tmp",
		AutoComplete:    readline.NewPrefixCompleter(getShellCommands()...),
		InterruptPrompt: "^C",
	})
	if err != nil {
		fmt.Println("Xatolik: Terminalni o'qib bo‘lmadi.")
		os.Exit(1)
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil {
			fmt.Println("\nDasturdan chiqildi.")
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// **Agar input ichida "aurora" bo‘lsa, xabar chiqaramiz**
		if strings.Contains(input, "aurora") || strings.Contains(input, "Aurora") {
			// TODO: Aurora bajaryapti...
			fmt.Println("Aurora bajaryapti...")
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Dasturdan chiqildi.")
			break
		}

		args := strings.Fields(input)

		// **Sudo tekshirish**
		if args[0] == "sudo" {
			if !sudoEnabled {
				fmt.Println("Sudo rejimi yoqilmagan. --sudo bilan dastur ishga tushiring yoki parolni kiriting.")
				continue
			}

			args = args[1:]
			cmd := exec.Command(userShell, "-i", "-c", fmt.Sprintf("echo %s | sudo -S -p '' %s", sudoPassword, strings.Join(args, " ")))
			runCommandWithPTY(cmd)
			continue
		}

		// **Shell muhitida ishga tushirish (Ranglar yo‘qolmasligi uchun)**
		cmd := exec.Command(userShell, "-i", "-c", input)
		runCommandWithPTY(cmd)
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
	out, err := exec.Command("bash", "-c", "compgen -c").Output()
	if err != nil {
		fmt.Println("Xatolik: Buyruqlarni olishda muammo bo'ldi.")
		return nil
	}

	commands := strings.Split(string(out), "\n")
	var completions []readline.PrefixCompleterInterface

	for _, cmd := range commands {
		if cmd != "" {
			completions = append(completions, readline.PcItem(cmd))
		}
	}

	extraCommands := []string{"exit", "quit", "clear"}
	for _, cmd := range extraCommands {
		completions = append(completions, readline.PcItem(cmd))
	}

	return completions
}

// **PTY orqali ranglarni saqlab buyruqni ishga tushirish**
func runCommandWithPTY(cmd *exec.Cmd) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("Xatolik:", err)
		return
	}
	defer ptmx.Close()

	// **Faol jarayonni saqlaymiz**
	activeCmd = cmd

	// **PTY output'ni terminalga yo‘naltiramiz**
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				break
			}
			os.Stdout.Write(buf[:n])
		}
	}()

	// **Buyruq bajarilishini kutish**
	cmd.Wait()

	// **Jarayon tugagandan keyin activeCmd'ni tozalash**
	activeCmd = nil
}
