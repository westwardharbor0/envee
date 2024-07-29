package envee

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestEnvee_SetPrefix(t *testing.T) {
	t.Parallel()

	e := New()
	e.SetPrefix("APP_")

	if e.prefix != "APP_" {
		t.Errorf("prefix should be APP_ but is %q", e.prefix)
	}
}

type testConfig struct {
	Name  string        `env:"NAME,required"`
	Age   uint32        `default:"42"        env:"AGE"`
	Score float32       `env:"SCORE"`
	TTL   time.Duration `env:"TTL"`
}

func TestEnvee_Parse(t *testing.T) {
	os.Clearenv()

	_ = os.Setenv("NAME", "username")
	_ = os.Setenv("AGE", "18")
	_ = os.Setenv("SCORE", "12.4")
	_ = os.Setenv("TTL", "1h")

	var c testConfig

	e := New()
	if err := e.Parse(&c); err != nil {
		t.Error(err)
	}

	if c.Name != "username" {
		t.Errorf("name should be username but is %q", c.Name)
	}

	if c.Age != 18 {
		t.Errorf("age should be 18 but is %d", c.Age)
	}

	if c.Score != 12.4 {
		t.Errorf("score should be 12.4 but is %f", c.Score)
	}

	if c.TTL != time.Hour {
		t.Errorf("ttl should be 1h but is %s", c.TTL)
	}
}

type testConfigAdvanced struct {
	Service string     `default:"default-service" env:"SERVICE"`
	Config  testConfig `prefix:"BASIC_"`
}

func TestEnvee_Parse_Advanced(t *testing.T) {
	os.Clearenv()

	_ = os.Setenv("BASIC_NAME", "username")
	_ = os.Setenv("BASIC_AGE", "21")
	_ = os.Setenv("BASIC_SCORE", "13.4")
	_ = os.Setenv("BASIC_TTL", "1h")

	var c testConfigAdvanced

	e := New()
	if err := e.Parse(&c); err != nil {
		t.Error(err)
	}

	if c.Service != "default-service" {
		t.Errorf("service should be default-service but is %s", c.Service)
	}

	if c.Config.Name != "username" {
		t.Errorf("name should be username but is %q", c.Config.Name)
	}

	if c.Config.Age != 21 {
		t.Errorf("age should be 21 but is %d", c.Config.Age)
	}

	if c.Config.Score != 13.4 {
		t.Errorf("score should be 13.4 but is %f", c.Config.Score)
	}

	if c.Config.TTL != time.Hour {
		t.Errorf("ttl should be 1h but is %s", c.Config.TTL)
	}
}

func TestEnvee_Parse_Default(t *testing.T) {
	os.Clearenv()

	_ = os.Setenv("NAME", "username2")
	_ = os.Setenv("SCORE", "0.4")
	_ = os.Setenv("TTL", "4h")

	var c testConfig

	e := New()
	if err := e.Parse(&c); err != nil {
		t.Error(err)
	}

	if c.Name != "username2" {
		t.Errorf("name should be username2 but is %q", c.Name)
	}

	if c.Age != 42 {
		t.Errorf("age should be 42 but is %d", c.Age)
	}
}

func TestEnvee_Parse_Required(t *testing.T) {
	os.Clearenv()

	var c testConfig

	e := New()
	err := e.Parse(&c)

	if err == nil {
		t.Error("expected error but got nil")
	}

	if err != nil && !errors.Is(err, ErrMissingRequired) {
		t.Errorf("expected ErrMissingRequired but was %s", err.Error())
	}
}

type testConfigDuration struct {
	Service        string                         `default:"default-service" env:"SERVICE"`
	ShutdownPeriod time.Duration                  `env:"SHUTDOWN_PERIOD"`
	Mongo          testConfigurationDurationMongo `prefix:"MONGO_"`
}

type testConfigurationDurationMongo struct {
	Host    string        `default:"localhost" env:"HOST"`
	Port    int           `default:"27017"     env:"PORT"`
	Timeout time.Duration `default:"13s"       env:"TIMEOUT"`
}

func TestEnvee_Parse_Duration(t *testing.T) {
	os.Clearenv()

	_ = os.Setenv("APP_SERVICE", "test-service")
	_ = os.Setenv("APP_SHUTDOWN_PERIOD", "9s")

	_ = os.Setenv("APP_MONGO_HOST", "not-localhost")
	_ = os.Setenv("APP_MONGO_TIMEOUT", "2s")

	var c testConfigDuration

	e := New()
	e.SetPrefix("APP_")

	if err := e.Parse(&c); err != nil {
		t.Error(err)
	}

	if c.Service != "test-service" {
		t.Errorf("service should be test-service but is %q", c.Service)
	}

	if c.ShutdownPeriod != 9*time.Second {
		t.Errorf("shutdown period should be 9s but is %q", c.ShutdownPeriod)
	}

	if c.Mongo.Host != "not-localhost" {
		t.Errorf("mongo host should be not-localhost but is %q", c.Mongo.Host)
	}

	if c.Mongo.Port != 27017 {
		t.Errorf("mongo port should be 27017 but is %q", c.Mongo.Port)
	}

	if c.Mongo.Timeout != 2*time.Second {
		t.Errorf("mongo timeout should be 2s but is %q", c.Mongo.Timeout)
	}
}

func TestField_isRequired(t *testing.T) {
	f := field{
		name:     "test",
		defValue: "yes",
		required: false,
	}

	if f.isRequired() {
		t.Error("field should not be required but is true")
	}

	f = field{
		name:     "test",
		defValue: "",
		required: true,
	}

	if !f.isRequired() {
		t.Error("field should be required but is false")
	}

	f = field{
		name:     "test",
		defValue: "",
		required: false,
	}

	if !f.isRequired() {
		t.Error("field should be required but is false")
	}
}

func TestParseTypeValue(t *testing.T) {
	v, err := parseTypeValue("3s", reflect.ValueOf(time.Second).Type())
	if err != nil {
		t.Error(err)
	}

	if v, _ := (v).(time.Duration); v != 3*time.Second {
		t.Errorf("value should be 3s but is %v", v)
	}

	_, err = parseTypeValue("3", reflect.ValueOf(time.Second).Type())
	if err == nil {
		t.Error("expected error but got nil")
	}

	if errors.Is(err, ErrUnsupportedType) {
		t.Error("expected other error than ErrUnsupportedType")
	}

	_, err = parseTypeValue("3", reflect.ValueOf(2).Type())
	if err == nil {
		t.Error("expected error but got nil")
	}

	if !errors.Is(err, ErrUnsupportedType) {
		t.Errorf("expected ErrUnsupportedType but got %s", err.Error())
	}
}

func TestParseKindValue(t *testing.T) {
	v, err := parseKindValue("8", reflect.ValueOf(uint(8)).Kind())
	if err != nil {
		t.Error(err)
	}

	if v != uint(8) {
		t.Errorf("value should be 8 but is %q", v)
	}

	v, err = parseKindValue("444444.444444", reflect.ValueOf(float64(8)).Kind())
	if err != nil {
		t.Error(err)
	}

	if v != float64(444444.444444) {
		t.Errorf("value should be 444444.444444 but is %q", v)
	}

	_, err = parseKindValue("asd", reflect.ValueOf(complex(8.1, 3)).Kind())
	if err == nil {
		t.Error("expected error but got nil")
	}

	if !errors.Is(err, ErrUnsupportedType) {
		t.Errorf("expected ErrUnsupportedType but got %s", err.Error())
	}
}
