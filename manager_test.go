package overseer

// Currently using testify/assert here
// and go-test/deep for cmd_test
// Not optimal
import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSimpleOverseer(t *testing.T) {
	assert := assert.New(t)
	ovr := NewOverseer()

	id := "echo"
	ovr.Add(id, "echo", "").Start()
	time.Sleep(timeUnit)

	stat := ovr.Status(id)
	assert.Equal(stat.Exit, 0, "Exit code should be 0")
	assert.NotEqual(ovr.Status("echo").PID, 0, "PID shouldn't be 0")

	id = "list"
	ovr.Add("list", "ls", "/usr/").Start()
	time.Sleep(timeUnit)

	stat = ovr.Status(id)
	assert.Equal(stat.Exit, 0, "Exit code should be 0")
	assert.NotEqual(ovr.Status("list").PID, 0, "PID shouldn't be 0")

	assert.Equal(2, len(ovr.ListAll()), "Expected 2 procs: echo, list")

	// Should not crash
	ovr.StopAll()
}

func TestSleepOverseer(t *testing.T) {
	assert := assert.New(t)
	ovr := NewOverseer()

	id := "sleep"
	ovr.Add(id, "sleep", "10").Start()
	time.Sleep(timeUnit)

	json := ovr.ToJSON(id)
	// JSON status should contain the same info
	assert.False(json.Complete)
	assert.Equal(json.ExitCode, -1)
	assert.True(json.PID > 0)
	assert.Nil(json.Error)

	// success stop
	assert.Nil(ovr.Stop(id))
	time.Sleep(timeUnit * 5)

	// proc was killed
	json = ovr.ToJSON(id)
	assert.Equal(json.ExitCode, -1)
	// assert.True(json.Complete)
}

func TestInvalidOverseer(t *testing.T) {
	assert := assert.New(t)
	ovr := NewOverseer()

	id := "err1"
	ovr.Add(id, "qwertyuiop", "zxcvbnm").Start()

	time.Sleep(timeUnit)
	stat := ovr.Status(id)
	json := ovr.ToJSON(id)

	assert.Equal(stat.Exit, -1, "Exit code should be negative")
	assert.NotEqual(stat.Error, nil, "Error shouldn't be nil")
	// JSON status should contain the same info
	assert.Equal(stat.Exit, json.ExitCode)
	assert.Equal(stat.Error, json.Error)
	assert.Equal(stat.PID, json.PID)
	assert.Equal(stat.Complete, json.Complete)

	// try to stop a dead process
	assert.Nil(ovr.Stop(id))

	id = "err2"
	ovr.Add(id, "ls", "/some_random_not_existent_path").Start()

	time.Sleep(timeUnit)
	stat = ovr.Status(id)

	assert.True(stat.Exit > 0, "Exit code should be positive")
	assert.Nil(stat.Error, "Error should be nil")
}

func TestSimpleSupervise(t *testing.T) {
	assert := assert.New(t)
	ovr := NewOverseer()

	ovr.Add("echo", "echo", "")
	id := "sleep"
	ovr.Add(id, "sleep", "1")

	ovr.Supervise(id) // To supervise sleep. How cool is that?
	stat := ovr.Status(id)

	assert.Equal(stat.Exit, 0, "Exit code should be 0")
	assert.Nil(stat.Error, "Error should be nil")

	assert.Equal(2, len(ovr.ListAll()), "Expected 2 procs: echo, sleep")
}
