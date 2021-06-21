package cmd

import (
	"bufio"
	"fmt"
	"github.com/daspoet/gowinkey"
	"github.com/mengdaming/tcr/trace"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/term"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tj/go-terminput"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var cfgFile string
var toolchain string
var autoPush bool

var rootCmd = &cobra.Command{
	Use:   "tcr",
	Short: "TCR (Test && Commit || Revert)",
	Long: `
This application is a tool to practice TCR.
It can be used either in solo, or as a group within a mob or pair session.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO hook up real application here

		trace.HorizontalLine()
		trace.Info("This is an information trace")
		trace.Warning("This is a warning trace")
		//trace.Error("This is an error trace")

		var toolchainTrace = fmt.Sprintf("Toolchain = %v", toolchain)
		trace.Info(toolchainTrace)
		var autoPushTrace = fmt.Sprintf("Auto-Push = %v", autoPush)
		trace.Info(autoPushTrace)

		// Experiments on keystrokes capture
		//exampleBlockingGetKey()
		//exampleKeyStrokeUsingChannel()
		//tryTermInput()
		//tryGoWinKey()
		//tryTerm()
		//tryTermNoEof()
		tryScanner()

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tcr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&toolchain, "toolchain", "t",
		"maven", "indicate the toolchain to be used by TCR")
	rootCmd.Flags().BoolVarP(&autoPush, "auto-push", "p", false, "Enable git push after every commit")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tcr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tcr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// Experiments

func tryTerm() {
	if ! terminal.IsTerminal(0) {
		b, _ := ioutil.ReadAll(os.Stdin)
		fmt.Println(string(b))
	} else {
		fmt.Println("no piped data")
	}
}

func tryTermNoEof() {
	r := strings.NewReader("some io.Reader stream to be read\n")

	buf := make([]byte, 1)
	if _, err := io.ReadAtLeast(os.Stdin, buf, 1); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)

	// buffer smaller than minimal read size.
	shortBuf := make([]byte, 3)
	if _, err := io.ReadAtLeast(r, shortBuf, 1); err != nil {
		fmt.Println("error:", err)
	}

	// minimal read size bigger than io.Reader stream
	longBuf := make([]byte, 64)
	if _, err := io.ReadAtLeast(r, longBuf, 1); err != nil {
		fmt.Println("error:", err)
	}

}


func tryTermInput() {
	t, err := term.Open("/dev/tty")
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	t.SetRaw()
	defer t.Restore()

	fmt.Printf("Type something, use 'q' to exit.\r\n")

	for {
		e, err := terminput.Read(t)
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}

		if e.Key() == terminput.KeyEscape || e.Rune() == 'q' {
			break
		}

		fmt.Printf("%s — shift=%v ctrl=%v alt=%v meta=%v\r\n", e.String(), e.Shift(), e.Ctrl(), e.Alt(), e.Meta())
	}
}

func tryGoWinKey() {
	events, _ := gowinkey.ListenSelective(gowinkey.VK_W, gowinkey.VK_A, gowinkey.VK_S, gowinkey.VK_D)

//	timer := time.AfterFunc(time.Second * 5, stopFn)
//	defer timer.Stop()


	for e := range events {
		switch e.Type {
		case gowinkey.KeyPressed:
			fmt.Println("pressed", e)
		case gowinkey.KeyReleased:
			fmt.Println("released", e)
		}
	}
}

func tryScanner() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println("You typed:", scanner.Text())
	}

	if scanner.Err() != nil {
		// handle error
	}
}

