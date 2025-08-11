package cmd

import (
	"bufio"
	"cut/internal/config"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var cfg config.Options

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cut",
	Short: "Утилита для вывода колонок",
	Long: `Утилита, которая считывает входные данные (STDIN) 
и разбивает каждую строку по заданному разделителю, 
после чего выводит определённые поля (колонки).

Использование:
./cut [-f "1,3-5"] [-d "\t"] [-s]

-f - нужные столбцы
-d - разделитель
-s - показывать только строки с разделителями`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := bufio.NewReaderSize(os.Stdin, 1024*1024)
		w := bufio.NewWriterSize(os.Stdout, 1024*1024)
		err := Process(r, w, cfg)
		if err != nil {
			return err
		}
		_ = w.Flush()
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	rootCmd.Flags().StringP("fields", "f", "", "1,3-5")
	rootCmd.Flags().StringP("delimiter", "d", "\t", "-d \"delimiter\"")
	rootCmd.Flags().BoolVarP(&cfg.SepOnly, "separated", "s", false, "-s")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		f, err := cmd.Flags().GetString("fields")
		if err != nil {
			return err
		}
		if f != "" {
			cfg.Fields, err = parseFields(f)
			if err != nil {
				return err
			}
			cfg.ShowAll = false
		} else {
			cfg.ShowAll = true
		}

		d, err := cmd.Flags().GetString("delimiter")
		if err != nil {
			return err
		}
		if len([]rune(d)) == 0 {
			return fmt.Errorf("delimiter must not be empty")
		}
		cfg.Delimiter = []rune(d)[0]
		return nil
	}
}

func parseFields(s string) (map[int]struct{}, error) {
	if len(s) == 0 {
		return nil, nil
	}
	re := regexp.MustCompile(`\d+(?:-\d+)?`)
	matches := re.FindAllString(s, -1)

	res := make(map[int]struct{})
	for _, m := range matches {
		if strings.Contains(m, "-") {
			parts := strings.SplitN(m, "-", 2)
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			if start > end {
				start, end = end, start
			}
			for i := start; i <= end; i++ {
				res[i] = struct{}{}
			}
		} else {
			num, _ := strconv.Atoi(m)
			res[num] = struct{}{}
		}
	}
	if len(res) == 0 {
		return nil, errors.New("gotta have at least some columns")
	}

	return res, nil
}
