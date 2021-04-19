package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/emar-kar/course_project/internal/logger"
	"github.com/emar-kar/course_project/pkg/request"
)

var buildVersion string
var finished []string = []string{
	":metal_tone2:",
	":clap_tone2:",
	":five-thirty:",
	":hand_over_mouth:",
	":mechanical_arm:",
	":relieved:",
	":vulcan_tone1:",
	":nail_care_tone1:",
	":monocle_face:",
}

type config struct {
	logFolder      string
	separator      string
	requestTimeout int
}

func main() {
	rand.Seed(time.Now().Unix())

	conf, err := initConfig()
	if err != nil {
		panic(fmt.Sprintf("cannot initialize config: %v", err))
	}

	logger.InitLogger(conf.logFolder)
	l := logger.GetLogger()
	defer l.Info.Sync()
	defer l.Error.Sync()

	printWelcome()

	quitCtx, quitCancel := context.WithCancel(context.Background())
	defer quitCancel()
	go watchSignals(quitCancel, l.Error)

	reqCh := make(chan string, 1)
	defer close(reqCh)
	nextCh := make(chan struct{}, 1)
	defer close(nextCh)
	nextCh <- struct{}{}
	go readRequest(quitCtx, reqCh, nextCh)

	for {
		select {
		case <-quitCtx.Done():
			emoji.Println("\nGood bye! :pensive_face:")
			return
		case requestString := <-reqCh:
			l.Info.Sugar().Infof("user request: %s", requestString)

			ctx, cancel := context.WithTimeout(quitCtx, time.Duration(conf.requestTimeout)*time.Second)
			defer cancel()

			start := time.Now()
			req, err := request.NewRequest(requestString)
			if err != nil {
				emoji.Printf("error: %s :sad_but_relieved_face:\n", err)
				l.Error.Sugar().Errorf("error: %s", err)
				nextCh <- struct{}{}
				continue
			}
			result, err := req.Do(ctx, conf.separator)
			if errors.Is(err, context.DeadlineExceeded) {
				emoji.Println(
					"deadline exceeded, try to increase the processing time in config file or specify the request :hammer_and_wrench:",
				)
				l.Error.Error(
					"deadline exceeded, try to increase the processing time in config file or specify the request",
				)
			} else if err != nil {
				emoji.Printf("error during request: %s :sweat:\n", err)
				l.Error.Sugar().Errorf("error during request: %s", err)
			}
			if result != nil {
				result.Print()
				emoji.Printf("request done in %s %s\n", time.Since(start), finished[rand.Intn(len(finished))])
			}
			nextCh <- struct{}{}
		}
	}
}

func printWelcome() {
	filePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("cannot define path to executed file")
	}

	emoji.Println("Welcome to csv-queuer :smiling_face_with_heart-eyes::smiling_face_with_halo::punch:")
	emoji.Println(
		"This tool allows you to play queue selection on given csv files ",
		"and print results with :smiling_face_with_three_hearts:",
	)
	emoji.Println("Please, read the instructions first and let's get started :nerd_face:")
	fmt.Printf("My current version: %s\n", buildVersion)
	fmt.Printf("App path: %s\n", filePath)
}

func watchSignals(quitCancel context.CancelFunc, logError *zap.Logger) {
	osSignalChan := make(chan os.Signal, 1)

	signal.Notify(osSignalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-osSignalChan
	fmt.Println()
	quitCancel()
	logError.Error("exited by user")
}

func readRequest(quitCtx context.Context, reqCh chan<- string, nextCh <-chan struct{}) {
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-quitCtx.Done():
			return
		case <-nextCh:
			fmt.Print("Your request: ")
			// Since ReadString in case of an error still returns something what was already read,
			// we can pass it to the request channel. It will fail later during request string
			// parsing procedure.
			requestString, _ := reader.ReadString(';')
			reqCh <- requestString
			reader.Reset(os.Stdin)
		}
	}
}

func initConfig() (*config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	separator := viper.GetString("csv_separator")
	if separator == "" {
		separator = ","
	}
	requestTimeout := viper.GetInt("request_timeout")
	if requestTimeout == 0 {
		requestTimeout = 5
	}

	logFolder := viper.GetString("log_folder")
	if logFolder == "" {
		logFolder = "./logs"
	}

	return &config{
		logFolder:      logFolder,
		separator:      separator,
		requestTimeout: requestTimeout,
	}, nil
}
