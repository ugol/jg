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

package emitter

import (
	"fmt"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/ctx"
	"github.com/ugol/jr/pkg/functions"
	"github.com/ugol/jr/pkg/loop"
	"github.com/ugol/jr/pkg/producers/console"
	"github.com/ugol/jr/pkg/producers/kafka"
	"github.com/ugol/jr/pkg/producers/mongoDB"
	"github.com/ugol/jr/pkg/producers/redis"
	"github.com/ugol/jr/pkg/tpl"
	"log"
	"os"
	"time"
)

type Emitter struct {
	Name           string        `mapstructure:"name"`
	Locale         string        `mapstructure:"locale"`
	Num            int           `mapstructure:"num"`
	Frequency      time.Duration `mapstructure:"frequency"`
	Duration       time.Duration `mapstructure:"duration"`
	Preload        int           `mapstructure:"preload"`
	ValueTemplate  string        `mapstructure:"valueTemplate"`
	KeyTemplate    string        `mapstructure:"keyTemplate"`
	OutputTemplate string        `mapstructure:"outputTemplate"`
	Output         string        `mapstructure:"output"`
	Topic          string        `mapstructure:"topic"`
	Kcat           bool          `mapstructure:"kcat"`
	Oneline        bool          `mapstructure:"oneline"`
	Producer       loop.Producer
}

func (e *Emitter) RunPreload(conf configuration.GlobalConfiguration) {

	keyTpl, err := tpl.NewTpl("key", e.KeyTemplate, functions.FunctionsMap(), &ctx.JrContext)
	if err != nil {
		log.Println(err)
	}
	templatePath := fmt.Sprintf("%s/%s.tpl", os.ExpandEnv(conf.TemplateDir), e.ValueTemplate)
	vt, err := os.ReadFile(templatePath)
	valueTpl, err := tpl.NewTpl("value", string(vt), functions.FunctionsMap(), &ctx.JrContext)
	if err != nil {
		log.Println(err)
	}

	// Preload
	for i := 0; i < e.Preload; i++ {
		k := keyTpl.Execute()
		v := valueTpl.Execute()
		e.Producer.Produce([]byte(k), []byte(v), nil)
		ctx.JrContext.GeneratedObjects++
		ctx.JrContext.GeneratedBytes += int64(len(v))
	}

}

func (e *Emitter) Initialize(conf configuration.GlobalConfiguration) {

	o, _ := tpl.NewTpl("out", e.OutputTemplate, functions.FunctionsMap(), nil)
	if e.Output == "stdout" {
		e.Producer = &console.KonsoleProducer{OutputTpl: &o}
		return
	}

	if e.Output == "kafka" {
		e.Producer = createKafkaProducer(conf, e.Topic, e.ValueTemplate)
		return
	} else {
		if conf.SchemaRegistry {
			log.Println("Ignoring schemaRegistry and/or serializer when output not set to kafka")
		}
	}

	if e.Output == "redis" {
		e.Producer = createRedisProducer(conf.RedisTtl, conf.RedisConfig)
		return
	}

	if e.Output == "mongo" || e.Output == "mongodb" {
		e.Producer = createMongoProducer(conf.MongoConfig)
		return
	}

	if e.Output == "http" {
		//e.Producer = &server.JsonProducer{OutTemplate: &o}
		// return
	}
}

/*
func (e *Emitter) CreateProducer() loop.Producer {
	o, _ := tpl.NewTpl("out", e.OutputTemplate, functions.FunctionsMap(), nil)
	return &console.KonsoleProducer{OutputTpl: &o}
}
*/

func createRedisProducer(ttl time.Duration, redisConfig string) loop.Producer {
	rProducer := &redis.RedisProducer{
		Ttl: ttl,
	}
	rProducer.Initialize(redisConfig)
	return rProducer
}

func createMongoProducer(mongoConfig string) loop.Producer {
	mProducer := &mongoDB.MongoProducer{}
	mProducer.Initialize(mongoConfig)

	return mProducer
}

func createKafkaProducer(conf configuration.GlobalConfiguration, topic string, templateType string) *kafka.KafkaManager {

	kManager := &kafka.KafkaManager{
		Serializer:   conf.Serializer,
		Topic:        topic,
		TemplateType: templateType,
	}

	kManager.Initialize(conf.KafkaConfig)

	if conf.SchemaRegistry {
		kManager.InitializeSchemaRegistry(conf.RegistryConfig)
	}
	if conf.AutoCreate {
		kManager.CreateTopic(topic)
	}
	return kManager
}
