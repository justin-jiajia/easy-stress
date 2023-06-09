package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	time_amount   int
	time_sent     int
	time_finished int
	time_ok       int
	time_err      int
	time_used     int
	website       string
	chanel        chan bool
	end           bool
	use_file      bool
	last_time     int
	file_name     string
	client        http.Client
	print_mu      sync.Mutex
	time_startrun time.Time
	write_mu      sync.Mutex
	csvwriter     *csv.Writer
)

func start(this_id int) {
	this_res_id := 1
	for time_sent < time_amount {
		time_sent++
		time_start := time.Now()
		_, err := client.Get(website)
		// resp, err := http.Get(website)
		time_spend := time.Since(time_start).Abs().Milliseconds()
		if use_file {
			write_mu.Lock()
			if err != nil {
				csvwriter.Write([]string{fmt.Sprintf("%f", time.Since(time_startrun).Seconds()), "error", fmt.Sprintf("%d", time_spend)})
			} else {
				csvwriter.Write([]string{fmt.Sprintf("%f", time.Since(time_startrun).Seconds()), "success", fmt.Sprintf("%d", time_spend)})
			}
			write_mu.Unlock()
		}
		if err != nil {
			time_err++
			print_mu.Lock()
			fmt.Print("\033[33mError: \033[0m")
			fmt.Println(err)
			print_mu.Unlock()
			// defer resp.Body.Close()
		} else {
			time_ok++
			last_time = int(time_spend)
			time_used += int(time_spend)
		}
		time_finished++
		this_res_id++
	}
	chanel <- true
}
func print_im() {
	for !end {
		print_mu.Lock()
		if time_ok != 0 {
			fmt.Printf("Sent:%d Success:%d Error:%d Average:%dms Last time:%dms        \r", time_finished, time_ok, time_err, time_used/time_ok, last_time)
		} else {
			fmt.Printf("Sent:%d Success:%d Error:%d      \r", time_finished, time_ok, time_err)
		}
		print_mu.Unlock()
	}
}
func main() {
	var c int
	app := &cli.App{
		Name:    "EasyStress",
		Usage:   "Send a lot of requests to test a website's stress resistance",
		Version: "v1.3",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "licence",
				Usage:   "Show the licence",
				Aliases: []string{"l"},
				Count:   &c,
			},
			&cli.IntFlag{
				Name:    "timeout",
				Usage:   "The timeout (unit: second)",
				Aliases: []string{"o"},
				Value:   5,
			},
			&cli.IntFlag{
				Name:    "time",
				Usage:   "The time of sending requests",
				Aliases: []string{"t"},
			},
			&cli.IntFlag{
				Name:    "worker",
				Usage:   "The amount of workers to send requests",
				Aliases: []string{"w"},
				Value:   4,
			},
			&cli.StringFlag{
				Name:        "file",
				Usage:       "The name of the csv file which contains every request's time and error",
				Aliases:     []string{"f"},
				DefaultText: "none",
			},
		},
		Action: func(ctx *cli.Context) error {
			if c != 0 {
				fmt.Print(`You should have known that this program is just made for test stress resistance of websites, not for DDOS attack.
To use this program, you must agree these terms:
- All the Consequences made by using this program should be borne by you, not by me (the author of this program).
- You use this program means you agree to all the terms.
- If you don't agree to these terms, you MUST stop using this program AT ONCE.`)
				return cli.Exit("", 0)
			}
			if ctx.NArg() != 1 {
				return cli.Exit("Only 1 argument(website) is allowed!", -1)
			}
			time_amount = ctx.Int("time")
			if time_amount == 0 {
				return cli.Exit("Required flag \"t\" not set", -1)
			}
			website = ctx.Args().Get(0)
			if website[:8] != "https://" && website[:7] != "http://" {
				return cli.Exit("Website should start with 'http://' or 'https://'.", -1)
			}
			use_file = ctx.String("file") != ""
			if use_file {
				file_name = ctx.String("file")
				var err error
				_, err = os.Stat(file_name)
				if err == nil {
					fmt.Println("\033[33mThe csv file exits.\033[0m")
					fmt.Print("(c)lean the file or (s)top this program: (c/s)")
					var chose string
					fmt.Scanf("%s", &chose)
					if chose == "S" || chose == "s" {
						return cli.Exit("", 0)
					}
				}
				var file *os.File
				file, err = os.Create(file_name)
				if err != nil {
					fmt.Print("\033[33mError: \033[0m")
					fmt.Println(err)
					return cli.Exit("", -1)
				} else {
					defer file.Close()
				}
				csvwriter = csv.NewWriter(file)
				csvwriter.Write([]string{"FinishTime(s)", "Status", "TimeCost(ms)"})
			}
			worker_amount := ctx.Int("worker")
			if worker_amount > time_amount {
				fmt.Println("\033[33mWarning:\033[0m workers more than requests")
			}
			timeout := ctx.Int("timeout")
			client = http.Client{Timeout: time.Duration(timeout) * time.Second}
			fmt.Print("\033[33mStarting...\033[0m")
			fmt.Printf(`
	target :%s
	worker :%d
	amount :%d
			`, website, worker_amount, time_amount)
			chanel = make(chan bool)
			end = false
			go print_im()
			time_startrun = time.Now()
			for i := 0; i < worker_amount; i++ {
				go start(i + 1)
			}
			for i := 0; i < worker_amount; i++ {
				<-chanel
			}
			end = true
			fmt.Print("\r\033[33mFinish!\033[0m")
			if use_file {
				csvwriter.Flush()
				fmt.Print(" Logs saved to ", file_name, " as csv")
			}
			fmt.Print("                                                              ")
			fmt.Printf(`
	-Sent:%d
	--Successful:%d
	--Error:%d
	-Average:%dms`, time_finished, time_ok, time_err, time_used/time_finished)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
