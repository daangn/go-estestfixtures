package main

import (
	"context"
	"os"
	"testing"
	"time"
)

var (
	mockESClient *MockESClient
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var err error
	mockESClient, err = NewMockESClient(ctx, "../fixturedata")
	if err != nil {
		panic(err)
	}
	runTests := m.Run()
	os.Exit(runTests)
}

func TestESClient_Foo(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "sample test code case",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockESClient.Reset(tt.args.ctx)
			if err != nil {
				t.Error(err)
			}

			if err = mockESClient.Client.Foo(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Foo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
