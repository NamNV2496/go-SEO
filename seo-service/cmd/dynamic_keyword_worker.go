package cmd

import (
	"log"
	"time"

	"github.com/namnv2496/seo/configs"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var CrawlerWorkerCmd = &cobra.Command{
	Use:   "crawler-worker",
	Short: "A simple web crawler worker",
	Run: func(cmd *cobra.Command, args []string) {
		InvokeCrawlerWorker(startCrawlerWorker)
	},
}

func InvokeCrawlerWorker(invokers ...any) *fx.App {
	config := configs.LoadConfig()
	app := fx.New(
		fx.StartTimeout(time.Second*10),
		fx.StopTimeout(time.Second*10),
		fx.Provide(),
		fx.Supply(
			config,
		),
		fx.Invoke(invokers...),
	)
	return app
}

func startCrawlerWorker(
	lc fx.Lifecycle,
	config *configs.Config,
) {
	log.Println("Start consumer")
	select {}
}
