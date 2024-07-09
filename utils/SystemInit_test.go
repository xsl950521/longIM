package utils

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInitRedis(t *testing.T) {
	//InitConfig()
	InitRedis()
	Convey("TestInitRedis", t, func() {
		ctx := context.Background()
		err := RDB.Set(context.Background(), "test", "111", 0).Err()
		So(err, ShouldBeNil)
		val, err := RDB.Get(ctx, "test").Result()
		So(val, ShouldEqual, "111")
	})
}
