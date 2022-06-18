package estestfixtures

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/olivere/elastic/v7"
)

var esDsnTestSourceContainer = ""
var esDsnTestTargetContainer = ""

func TestMain(m *testing.M) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		esDsn, err := runESContainer("esfixture-test-source-es")
		if err != nil {
			panic(err)
		}
		esDsnTestSourceContainer = fmt.Sprintf("http://%s", esDsn)
		wg.Done()
	}()
	go func() {
		esDsn, err := runESContainer("esfixture-test-target-es")
		if err != nil {
			panic(err)
		}
		esDsnTestTargetContainer = fmt.Sprintf("http://%s", esDsn)
		wg.Done()
	}()
	go func() {
		<-time.After(60 * time.Second)
		wg.Done()
		wg.Done()
	}()
	wg.Wait()

	if esDsnTestSourceContainer == "" || esDsnTestTargetContainer == "" {
		panic("test es containers booting is too long. please check your docker and machine's resources")
	}
	fmt.Println("esDsnTestSourceContainer => ", esDsnTestSourceContainer)
	fmt.Println("esDsnTestTargetContainer => ", esDsnTestTargetContainer)
	fmt.Println("---------- test containers are start ------------")
	<-time.After(5 * time.Second)
	os.Exit(m.Run())
}

func TestNewLoader(t *testing.T) {
	type fields struct {
		limit       int
		dir         string
		targetNames []string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success case",
			fields: fields{
				limit: 10,
				targetNames: []string{
					"posts",         // using alias mode
					"faqs_20220529", // not using, just index name
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceLoadLoader, err := NewLoader(
				tt.args.ctx,
				esDsnTestSourceContainer,
				WithLimit(tt.fields.limit),
				WithTargetNames(tt.fields.targetNames...),
				WithDirectory("./fixturedata"),
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("sourceLoader - NewLoader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = sourceLoadLoader.Load(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("sourceLoader Loader.Load() error = %v", err)
			}
			t.Log("sourceLoadLoader  Loader.Load() is done")

			sourceDumpLoader, err := NewLoader(
				tt.args.ctx,
				esDsnTestSourceContainer,
				WithLimit(tt.fields.limit),
				WithTargetNames(tt.fields.targetNames...),
				WithDirectory("./testdata"),
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("sourceDumpLoader - NewLoader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = sourceDumpLoader.Dump(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("sourceDumpLoader Loader.Dump() error = %v", err)
			}
			t.Log("sourceDumpLoader Loader.Dump() is done")

			targetLoader, err := NewLoader(
				tt.args.ctx,
				esDsnTestTargetContainer,
				WithLimit(tt.fields.limit),
				WithTargetNames(tt.fields.targetNames...),
				WithDirectory("./testdata"),
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("targetLoader - NewLoader() error = %v", err)
			}
			err = targetLoader.Load(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("targetLoader - Loader.Load() error = %v", err)
			}
			t.Log("targetLoader Loader.Load() is done")

		})
	}
}

func TestLoader_clearElasticsearch(t *testing.T) {
	type fields struct {
		dsn         string
		es          *elastic.Client
		limit       int
		dir         string
		searchFunc  func(c *elastic.Client, indexes []string) *elastic.SearchService
		targetNames []string
	}
	type args struct {
		ctx context.Context
		dsn string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success case",
			fields: fields{
				limit: 10,
				targetNames: []string{
					"posts",         // using alias mode
					"faqs_20220529", // not using, just index name
				},
				dir: "./testdata",
			},
			args: args{
				ctx: context.Background(),
				dsn: esDsnTestSourceContainer,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := NewLoader(
				tt.args.ctx,
				tt.args.dsn,
				WithSearchFunc(func(c *elastic.Client, targetNames []string) *elastic.SearchService {
					return c.Search(targetNames...).Size(0).From(10)
				}),
				WithLimit(tt.fields.limit),
				WithTargetNames(tt.fields.targetNames...),
				WithDirectory("./fixturedata"),
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Loader.Load() error = %v", err)
			}

			if err = l.ClearElasticsearch(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Loader.clearElasticsearch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
