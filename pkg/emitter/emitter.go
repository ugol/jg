// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package emitter

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jrnd-io/jr/pkg/producers/wasm"

	"github.com/jrnd-io/jr/pkg/configuration"
	"github.com/jrnd-io/jr/pkg/constants"
	jtctx "github.com/jrnd-io/jr/pkg/ctx"
	"github.com/jrnd-io/jr/pkg/functions"
	"github.com/jrnd-io/jr/pkg/producers/awsdynamodb"
	"github.com/jrnd-io/jr/pkg/producers/azblobstorage"
	"github.com/jrnd-io/jr/pkg/producers/azcosmosdb"
	"github.com/jrnd-io/jr/pkg/producers/cassandra"
	"github.com/jrnd-io/jr/pkg/producers/console"
	"github.com/jrnd-io/jr/pkg/producers/elastic"
	"github.com/jrnd-io/jr/pkg/producers/gcs"
	"github.com/jrnd-io/jr/pkg/producers/http"
	"github.com/jrnd-io/jr/pkg/producers/kafka"
	"github.com/jrnd-io/jr/pkg/producers/luascript"
	"github.com/jrnd-io/jr/pkg/producers/mongodb"
	"github.com/jrnd-io/jr/pkg/producers/redis"
	"github.com/jrnd-io/jr/pkg/producers/s3"
	"github.com/jrnd-io/jr/pkg/producers/server"
	"github.com/jrnd-io/jr/pkg/producers/wamp"
	"github.com/jrnd-io/jr/pkg/tpl"
	"github.com/rs/zerolog/log"
)

type Emitter struct {
	Name             string        `mapstructure:"name"`
	Locale           string        `mapstructure:"locale"`
	Num              int           `mapstructure:"num"`
	Frequency        time.Duration `mapstructure:"frequency"`
	Duration         time.Duration `mapstructure:"duration"`
	Preload          int           `mapstructure:"preload"`
	ValueTemplate    string        `mapstructure:"valueTemplate"`
	EmbeddedTemplate string        `mapstructure:"embeddedTemplate"`
	KeyTemplate      string        `mapstructure:"keyTemplate"`
	OutputTemplate   string        `mapstructure:"outputTemplate"`
	Output           string        `mapstructure:"output"`
	Topic            string        `mapstructure:"topic"`
	Kcat             bool          `mapstructure:"kcat"`
	Oneline          bool          `mapstructure:"oneline"`
	Csv              string        `mapstructure:"csv"`
	GeoJson          string        `mapstructure:"geojson"`
	Producer         Producer
	KTpl             tpl.Tpl
	VTpl             tpl.Tpl
}

func (e *Emitter) Initialize(ctx context.Context, conf configuration.GlobalConfiguration) {

	functions.InitCSV(e.Csv)

	functions.InitGeoJson(e.GeoJson)

	templateName := e.ValueTemplate
	if e.EmbeddedTemplate == "" {
		path := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.JR_SYSTEM_DIR, "templates"))
		templateFullPath := fmt.Sprintf("%s/%s.tpl", path, templateName)
		vt, err := os.ReadFile(templateFullPath)
		e.EmbeddedTemplate = string(vt)
		if err != nil {
			log.Printf("Template '%s' not found in %s\n", templateName, path)
		}
	}

	keyTpl, err := tpl.NewTpl("key", e.KeyTemplate, functions.FunctionsMap(), &jtctx.JrContext)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create key template")
	}
	valueTpl, err := tpl.NewTpl("value", e.EmbeddedTemplate, functions.FunctionsMap(), &jtctx.JrContext)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create value template")
	}

	e.KTpl = keyTpl
	e.VTpl = valueTpl

	o, _ := tpl.NewTpl("out", e.OutputTemplate, functions.FunctionsMap(), nil)
	if e.Output == "stdout" {
		e.Producer = &console.Producer{OutputTpl: &o}
		return
	}

	if e.Output == "kafka" {
		e.Producer = createKafkaProducer(ctx, conf, e.Topic, templateName)
		return
	}

	if conf.SchemaRegistry {
		log.Warn().Msg("Ignoring schemaRegistry and/or serializer when output not set to kafka")
	}

	if e.Output == "redis" {
		e.Producer = createRedisProducer(ctx, conf.RedisTtl, conf.RedisConfig)
		return
	}

	if e.Output == "redishash" {
		e.Producer = createRedisHashProducer(ctx, conf.RedisTtl, conf.RedisConfig)
		return
	}

	if e.Output == "redisjson" {
		e.Producer = createRedisJSONProducer(ctx, conf.RedisTtl, conf.RedisConfig)
		return
	}

	if e.Output == "mongo" || e.Output == "mongodb" {
		e.Producer = createMongoProducer(ctx, conf.MongoConfig)
		return
	}

	if e.Output == "elastic" {
		e.Producer = createElasticProducer(ctx, conf.ElasticConfig)
		return
	}

	if e.Output == "s3" {
		e.Producer = createS3Producer(ctx, conf.S3Config)
		return
	}

	if e.Output == "awsdynamodb" {
		e.Producer = createAWSDynamoDB(ctx, conf.AWSDynamoDBConfig)
		return
	}

	if e.Output == "gcs" {
		e.Producer = createGCSProducer(ctx, conf.GCSConfig)
		return
	}

	if e.Output == "azblobstorage" {
		e.Producer = createAZBlobStorageProducer(ctx, conf.AzBlobStorageConfig)
		return
	}
	if e.Output == "azcosmosdb" {
		e.Producer = createAZCosmosDBProducer(ctx, conf.AzCosmosDBConfig)
		return
	}

	if e.Output == "json" {
		e.Producer = &server.JsonProducer{OutputTpl: &o}
		return
	}

	if e.Output == "http" {
		e.Producer = createHTTPProducer(ctx, conf.HTTPConfig)
		return
	}

	if e.Output == "cassandra" {
		e.Producer = createCassandraProducer(ctx, conf.CassandraConfig)
		return
	}
	if e.Output == "luascript" {
		e.Producer = createLUAScriptProducer(ctx, conf.LUAScriptConfig)
		return
	}
	if e.Output == "wasm" {
		e.Producer = createWASMProducer(ctx, conf.LUAScriptConfig)
		return
	}
	if e.Output == "wamp" {
		e.Producer = createWAMPProducer(ctx, conf.WAMPConfig)
		return
	}

}

func (e *Emitter) Run(ctx context.Context, num int, o any) {

	for i := 0; i < num; i++ {

		k := e.KTpl.Execute()
		v := e.VTpl.Execute()
		kInValue := functions.GetV("KEY")

		if kInValue != "" {
			e.Producer.Produce(ctx, []byte(kInValue), []byte(v), o)
		} else {
			e.Producer.Produce(ctx, []byte(k), []byte(v), o)
		}
		jtctx.JrContext.GeneratedObjects++
		jtctx.JrContext.GeneratedBytes += int64(len(v))

	}

}

func createRedisProducer(_ context.Context, ttl time.Duration, redisConfig string) Producer {
	rProducer := &redis.Producer{
		Ttl: ttl,
	}
	rProducer.Initialize(redisConfig)
	return rProducer
}

func createRedisHashProducer(_ context.Context, ttl time.Duration, redisConfig string) Producer {
	rProducer := &redis.HashProducer{
		Ttl: ttl,
	}
	rProducer.Initialize(redisConfig)
	return rProducer
}

func createRedisJSONProducer(_ context.Context, ttl time.Duration, redisConfig string) Producer {
	rProducer := &redis.JSONProducer{
		Ttl: ttl,
	}
	rProducer.Initialize(redisConfig)
	return rProducer
}

func createMongoProducer(ctx context.Context, mongoConfig string) Producer {
	mProducer := &mongodb.MongoProducer{}
	mProducer.Initialize(ctx, mongoConfig)

	return mProducer
}

func createElasticProducer(_ context.Context, elasticConfig string) Producer {
	eProducer := &elastic.Producer{}
	eProducer.Initialize(elasticConfig)

	return eProducer
}

func createS3Producer(ctx context.Context, s3Config string) Producer {
	sProducer := &s3.Producer{}
	sProducer.Initialize(ctx, s3Config)

	return sProducer
}

func createAWSDynamoDB(ctx context.Context, config string) Producer {
	producer := &awsdynamodb.Producer{}
	producer.Initialize(ctx, config)

	return producer
}

func createAZBlobStorageProducer(ctx context.Context, azConfig string) Producer {
	producer := &azblobstorage.Producer{}
	producer.Initialize(ctx, azConfig)

	return producer
}

func createAZCosmosDBProducer(_ context.Context, azConfig string) Producer {
	producer := &azcosmosdb.Producer{}
	producer.Initialize(azConfig)

	return producer
}

func createGCSProducer(ctx context.Context, gcsConfig string) Producer {
	gProducer := &gcs.Producer{}
	gProducer.Initialize(ctx, gcsConfig)

	return gProducer
}

func createHTTPProducer(_ context.Context, httpConfig string) Producer {
	httpProducer := &http.Producer{}
	httpProducer.Initialize(httpConfig)

	return httpProducer
}

func createCassandraProducer(_ context.Context, config string) Producer {
	producer := &cassandra.Producer{}
	producer.Initialize(config)

	return producer
}

func createLUAScriptProducer(_ context.Context, config string) Producer {
	producer := &luascript.Producer{}
	producer.Initialize(config)

	return producer
}

func createWASMProducer(ctx context.Context, config string) Producer {
	producer := &wasm.Producer{}
	producer.Initialize(ctx, config)

	return producer
}

func createWAMPProducer(ctx context.Context, config string) Producer {
	producer := &wamp.Producer{}
	producer.Initialize(ctx, config)

	return producer
}

func createKafkaProducer(ctx context.Context, conf configuration.GlobalConfiguration, topic string, templateType string) *kafka.Manager {

	kManager := &kafka.Manager{
		Serializer:   conf.Serializer,
		Topic:        topic,
		TemplateType: templateType,
	}

	kManager.Initialize(conf.KafkaConfig)

	if conf.SchemaRegistry {
		kManager.InitializeSchemaRegistry(conf.RegistryConfig)
	}
	if conf.AutoCreate {
		kManager.CreateTopic(ctx, topic)
	}
	return kManager
}
