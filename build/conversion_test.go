package build

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)


func TestLoadYaml(t *testing.T) {

	Convey("setup works", t, func() {
		tasks := readBook(".files/1st.yaml")
		So(tasks, ShouldNotBeNil)
		So(tasks[0].Series, ShouldNotBeNil)
		So(tasks[0].Series.Len(), ShouldEqual, 2)
		So(tasks[0].Series[1], ShouldNotBeNil)
		So(tasks[0].Series[1].Series.Len(), ShouldEqual, 4)
	})
}
