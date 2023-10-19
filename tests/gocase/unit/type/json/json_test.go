/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package json

import (
	"context"
	"testing"

	"github.com/apache/kvrocks/tests/gocase/util"
	"github.com/stretchr/testify/require"
)

func TestJson(t *testing.T) {
	srv := util.StartServer(t, map[string]string{})
	defer srv.Close()
	ctx := context.Background()
	rdb := srv.NewClient()
	defer func() { require.NoError(t, rdb.Close()) }()

	t.Run("JSON.SET and JSON.GET basics", func(t *testing.T) {
		require.NoError(t, rdb.Do(ctx, "JSON.SET", "a", "$", ` {"x":1, "y":2} `).Err())
		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a").Val(), `{"x":1,"y":2}`)

		require.NoError(t, rdb.Do(ctx, "JSON.SET", "a", "$.y", `233`).Err())
		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a").Val(), `{"x":1,"y":233}`)

		require.NoError(t, rdb.Do(ctx, "JSON.SET", "a", "$", `[[1], [2]]`).Err())
		require.NoError(t, rdb.Do(ctx, "JSON.SET", "a", "$[*][0]", "3").Err())
		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a").Val(), `[[3],[3]]`)

		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a", "$").Val(), `[[[3],[3]]]`)
		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a", "$[0]").Val(), `[[3]]`)
		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a", "$[0][0]").Val(), `[3]`)

		require.NoError(t, rdb.Do(ctx, "JSON.SET", "a", "$", `{"x":1,"y":{"x":{"y":2},"y":3}}`).Err())
		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a").Val(), `{"x":1,"y":{"x":{"y":2},"y":3}}`)
		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a", "$").Val(), `[{"x":1,"y":{"x":{"y":2},"y":3}}]`)
		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a", "$..x").Val(), `[1,{"y":2}]`)
		require.Equal(t, rdb.Do(ctx, "JSON.GET", "a", "$..x", "$..y").Val(), `{"$..x":[1,{"y":2}],"$..y":[{"x":{"y":2},"y":3},3,2]}`)
	})
}