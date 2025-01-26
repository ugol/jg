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

package configuration

import "time"

var GlobalCfg GlobalConfiguration

type GlobalConfiguration struct {
	Seed                int64
	KafkaConfig         string
	SchemaRegistry      bool
	RegistryConfig      string
	Serializer          string
	AutoCreate          bool
	RedisTtl            time.Duration
	RedisConfig         string
	MongoConfig         string
	AzBlobStorageConfig string
	AzCosmosDBConfig    string
	ElasticConfig       string
	S3Config            string
	GCSConfig           string
	HTTPConfig          string
	CassandraConfig     string
	AWSDynamoDBConfig   string
	LUAScriptConfig     string
	WASMConfig          string
	WAMPConfig          string
	WAMPRPCConfig       string
	Url                 string
	EmbeddedTemplate    bool
	FileNameTemplate    bool
}
