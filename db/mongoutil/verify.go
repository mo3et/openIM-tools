// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongoutil

import (
	"context"

	"github.com/openimsdk/tools/errs"
)

// CheckMongo tests the MongoDB connection without retries.
func Check(ctx context.Context, config *Config) error {
	client, err := NewMongoDB(ctx, config)
	if err != nil {
		if config.MongoMode == ReplicaSetMode {
			return errs.WrapMsg(err, "failed to connect to MongoDB replica set", "config", config)
		} else {
			return errs.WrapMsg(err, "MongoDB connect failed", "URI", config.Uri, "Database", config.Database, "MaxPoolSize", config.MaxPoolSize)
		}
	}

	defer func() {
		if client != nil && client.db != nil {
			client.db.Client().Disconnect(ctx)
		}
	}()

	return nil
}

// ValidateAndSetDefaults validates the configuration and sets default values.
func (c *Config) ValidateAndSetDefaults() error {
	if c.Uri == "" && len(c.Address) == 0 {
		return errs.Wrap(errs.New("either Uri or Address must be provided"))
	}
	if c.Database == "" {
		return errs.Wrap(errs.New("database is required"))
	}
	if c.MaxPoolSize <= 0 {
		c.MaxPoolSize = defaultMaxPoolSize
	}
	if c.MaxRetry <= 0 {
		c.MaxRetry = defaultMaxRetry
	}
	if c.Uri == "" {
		// if authSource is not provided, default to database name
		if c.AuthSource == "" {
			c.Uri = buildMongoURI(c, c.Database)
		} else {
			c.Uri = buildMongoURI(c, c.AuthSource)
		}
	}
	return nil
}
