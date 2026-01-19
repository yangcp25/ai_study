package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"softgen/internal/app"
)

var (
	name  string
	typ   string
	model uint
)

var rootCmd = &cobra.Command{
	Use:   "softgen",
	Short: "软件著作权材料生成工具",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.Run(name, typ, model)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "软件名称")
	rootCmd.Flags().StringVarP(&typ, "type", "t", "", "生成类型: manual | code | all")
	rootCmd.Flags().UintVarP(&model, "model", "m", 1, "模型：0 deepseek-chat | 1 deepseek-reasoner ")

	rootCmd.MarkFlagRequired("name")
	rootCmd.MarkFlagRequired("type")
}
