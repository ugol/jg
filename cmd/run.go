//Copyright © 2022 Ugo Landini <ugo.landini@gmail.com>
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ugol/jr/functions"
	"github.com/ugol/jr/producers/console"
	"github.com/ugol/jr/producers/kafka"
	"github.com/ugol/jr/producers/redis"
	"log"
	"os"
	"os/signal"
	"regexp"
	"text/template"
	"time"
)

type Producer interface {
	Close()
	Produce(k []byte, v []byte)
}

var runCmd = &cobra.Command{
	Use:   "run [template]",
	Short: "Execute a template",
	Long: `Execute a template. 
  Without any other flag, [template] is just the name of a template in the templates directory, which by default is in '$HOME/.jr/templates' Example: 
jr run net-device
  With the --template flag, [template] is a string containing a full template. Example:
jr run --template "{{name}}"
 With the -templateFileName flag [template] is a file name with a template. Example:
jr run --templateFileName ~/.jr/templates/net-device.tpl
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		keyTemplate, _ := cmd.Flags().GetString("key")
		outputTemplate, _ := cmd.Flags().GetString("outputTemplate")
		embeddedTemplate, _ := cmd.Flags().GetBool("template")
		templateFileName, _ := cmd.Flags().GetBool("templateFileName")
		kcat, _ := cmd.Flags().GetBool("kcat")
		output, _ := cmd.Flags().GetString("output")
		oneline, _ := cmd.Flags().GetBool("oneline")
		locales, _ := cmd.Flags().GetStringSlice("locales")

		num, _ := cmd.Flags().GetInt("num")
		frequency, _ := cmd.Flags().GetDuration("frequency")
		duration, _ := cmd.Flags().GetDuration("duration")
		seed, _ := cmd.Flags().GetInt64("seed")
		kafkaConfig, _ := cmd.Flags().GetString("kafkaConfig")
		registryConfig, _ := cmd.Flags().GetString("registryConfig")
		topic, _ := cmd.Flags().GetString("topic")

		templateDir, _ := cmd.Flags().GetString("templateDir")
		templateDir = os.ExpandEnv(templateDir)

		autocreate, _ := cmd.Flags().GetBool("autocreate")
		schemaRegistry, _ := cmd.Flags().GetBool("schemaRegistry")
		serializer, _ := cmd.Flags().GetString("serializer")

		redisTtl, _ := cmd.Flags().GetDuration("redis.ttl")
		redisConfig, _ := cmd.Flags().GetString("redisConfig")

		if kcat {
			oneline = true
			output = "stdout"
			outputTemplate = "{{.K}},{{.V}}\n"
		}

		var valueTemplate []byte
		var err error

		if embeddedTemplate {
			valueTemplate = []byte(args[0])
		} else if templateFileName {
			valueTemplate, err = os.ReadFile(os.ExpandEnv(args[0]))
			functions.JrContext.TemplateType = args[0]
		} else {
			templatePath := fmt.Sprintf("%s/%s.tpl", templateDir, args[0])
			valueTemplate, err = os.ReadFile(templatePath)
			functions.JrContext.TemplateType = args[0]
		}
		if err != nil {
			log.Fatal(err)
		}

		outTemplate, err := template.New("out").Parse(outputTemplate)
		if err != nil {
			log.Fatal(err)
		}

		key, err := template.New("key").Funcs(functions.FunctionsMap()).Parse(keyTemplate)
		if err != nil {
			log.Fatal(err)
		}

		value, err := template.New("value").Funcs(functions.FunctionsMap()).Parse(string(valueTemplate))
		if err != nil {
			log.Fatal(err)
		}

		var producer Producer

		if output == "stdout" {
			producer = &console.ConsoleProducer{OutTemplate: outTemplate}
		}

		if output == "kafka" {
			producer = createKafkaProducer(serializer, topic, kafkaConfig, schemaRegistry, registryConfig, kcat, autocreate)
		} else {
			if schemaRegistry {
				log.Println("Ignoring schemaRegistry and/or serializer when output not set to kafka")
			}
		}

		if output == "redis" {
			producer = createRedisProducer(redisTtl, redisConfig)
		}

		if output == "mongo" {
			log.Fatal("Not yet implemented")
		}

		functions.Random.Seed(seed)
		functions.JrContext.Num = num
		functions.JrContext.Range = make([]int, num)
		functions.JrContext.Frequency = frequency
		functions.JrContext.Locales = locales
		functions.JrContext.Seed = seed
		functions.JrContext.TemplateDir = templateDir

		infinite := true
		if duration > 0 {
			timer := time.NewTimer(duration)

			go func() {
				<-timer.C
				infinite = false
			}()
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		if frequency != -1 {
		Infinite:
			for ok := true; ok; ok = infinite {
				select {
				case <-time.After(frequency):
					for range functions.JrContext.Range {
						k, v, _ := executeTemplate(key, value, oneline)
						producer.Produce([]byte(k), []byte(v))
					}
				case <-ctx.Done():
					stop()
					break Infinite
				}
			}
		} else {
			for range functions.JrContext.Range {
				k, v, _ := executeTemplate(key, value, oneline)
				producer.Produce([]byte(k), []byte(v))

			}
		}

		producer.Close()

		time.Sleep(100 * time.Millisecond)
		writeStats()

	},
}

func createRedisProducer(ttl time.Duration, redisConfig string) Producer {
	rProducer := &redis.RedisProducer{
		Ttl: ttl,
	}
	rProducer.Initialize(redisConfig)
	return rProducer
}

func createKafkaProducer(serializer string, topic string, kafkaConfig string, schemaRegistry bool, registryConfig string, kcat bool, autocreate bool) *kafka.KafkaManager {
	kManager := &kafka.KafkaManager{
		Serializer:   serializer,
		Topic:        topic,
		TemplateType: functions.JrContext.TemplateType,
	}

	kManager.Initialize(kafkaConfig)

	if schemaRegistry {
		kManager.InitializeSchemaRegistry(registryConfig)
		if kcat {
			log.Println("Ignoring kcat when schemaRegistry is enabled")
		}
	}
	if autocreate {
		kManager.CreateTopic(topic)
	}
	return kManager
}

func writeStats() {
	_, _ = fmt.Fprintln(os.Stderr)
	elapsed := time.Since(functions.JrContext.StartTime)
	_, _ = fmt.Fprintf(os.Stderr, "Elapsed time: %v\n", elapsed.Round(1*time.Second))
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (Objects): %d\n", functions.JrContext.GeneratedObjects)
	_, _ = fmt.Fprintf(os.Stderr, "Data Generated (bytes): %d\n", functions.JrContext.GeneratedBytes)
	_, _ = fmt.Fprintf(os.Stderr, "Throughput (bytes per second): %9.f\n", float64(functions.JrContext.GeneratedBytes)/elapsed.Seconds())
	_, _ = fmt.Fprintln(os.Stderr)
}

func executeTemplate(key *template.Template, value *template.Template, oneline bool) (string, string, error) {

	var kBuffer, vBuffer bytes.Buffer
	var err error

	if err = key.Execute(&kBuffer, functions.JrContext); err != nil {
		log.Println(err)
	}
	k := kBuffer.String()

	if err = value.Execute(&vBuffer, functions.JrContext); err != nil {
		log.Println(err)
	}
	v := vBuffer.String()

	if oneline {
		re := regexp.MustCompile(`\r?\n?`)
		v = re.ReplaceAllString(v, "")
	}

	functions.JrContext.GeneratedObjects++
	functions.JrContext.GeneratedBytes += int64(len(v))

	return k, v, err
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntP("num", "n", functions.JrContext.Num, "Number of elements to create for each pass")
	runCmd.Flags().DurationP("frequency", "f", -1, "how much time to wait for next generation pass")
	runCmd.Flags().DurationP("duration", "d", 0, "If frequency is enabled, with Duration you can set a finite amount of time")

	runCmd.Flags().Int64("seed", functions.JrContext.Seed, "Seed to init pseudorandom generator")

	runCmd.Flags().String("templateDir", functions.JrContext.TemplateDir, "directory containing templates")
	runCmd.Flags().StringP("kafkaConfig", "F", "./kafka/config.properties", "Kafka configuration")
	runCmd.Flags().String("registryConfig", "./kafka/registry.properties", "Kafka configuration")
	runCmd.Flags().Bool("templateFileName", false, "If enabled, [template] must be a template file")
	runCmd.Flags().Bool("template", false, "If enabled, [template] must be a string containing a template, to be embedded directly in the script")

	runCmd.Flags().StringP("key", "k", "key", "A template to generate a key")
	runCmd.Flags().StringP("topic", "t", "test", "Kafka topic name")

	runCmd.Flags().Bool("kcat", false, "If you want to pipe jr with kcat, use this flag: it is equivalent to --output stdout --outputTemplate '{{key}},{{value}}' --oneline")
	runCmd.Flags().StringP("output", "o", "stdout", "can be one of stdout, kafka, redis, mongo")
	runCmd.Flags().String("outputTemplate", "{{.V}}\n", "Formatting of K,V on standard output")
	runCmd.Flags().BoolP("oneline", "l", false, "strips /n from output, for example to be pipelined to tools like kcat")
	runCmd.Flags().BoolP("autocreate", "a", false, "if enabled, autocreate topics")
	runCmd.Flags().StringSlice("locales", functions.JrContext.Locales, "List of locales")

	runCmd.Flags().BoolP("schemaRegistry", "s", false, "If you want to use Confluent Schema Registry")
	runCmd.Flags().String("serializer", "json-schema", "Type of serializer: json-schema, avro-generic, avro, protobuf")
	runCmd.Flags().Duration("redis.ttl", 1*time.Minute, "If output is redis, ttl of the object")
	runCmd.Flags().String("redisConfig", "./redis/config.json", "Redis configuration")

}
