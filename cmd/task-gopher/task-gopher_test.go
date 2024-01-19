package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		new  Task
		old  Task
		want Task
	}{
		{
			new: Task{
				ID:      1,
				Name:    "name",
				Tag:     "",
				Desc:    "test",
				Status:  invalidStatus,
				Created: time.Date(2023, 11, 18, 7, 43, 34, 1, time.UTC),
			},
			old: Task{
				ID:      1,
				Name:    "",
				Tag:     "tag",
				Desc:    "",
				Status:  inProgress,
				Created: time.Date(2023, 11, 18, 7, 42, 34, 1, time.UTC),
			},
			want: Task{
				ID:      1,
				Name:    "name",
				Tag:     "tag",
				Desc:    "test",
				Status:  inProgress,
				Created: time.Date(2023, 11, 18, 7, 42, 34, 1, time.UTC),
			},
		},
	}
	for _, tc := range tests {
		tc.old.merge(tc.new)
		if !reflect.DeepEqual(tc.old, tc.want) {
			t.Fatalf("got: %#v, want %#v", tc.new, tc.want)
		}
	}
}

func TestAddTask(t *testing.T) {

	var tests = []struct {
		name  string
		input Task
		want  Task
	}{
		{"empty task should match", Task{}, Task{}},
		{"task with name only should match", Task{Name: "test"}, Task{Name: "test"}},
		{"task with all fields should match", Task{
			Name:   "full",
			Tag:    "tag",
			Desc:   "desc",
			Status: inProgress,
		}, Task{
			Name:   "full",
			Tag:    "tag",
			Desc:   "desc",
			Status: inProgress,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db = setupTests()
			defer teardownTests(db)
			id, err := addTask(db, tt.input.Name, tt.input.Desc, tt.input.Status, tt.input.Tag)
			if err != nil {
				log.Fatal(err)
			}
			ans, err := getTask(db, id)
			if err != nil {
				log.Fatal(err)
			}
			// fields that we don't know in advance
			tt.want.ID = id
			tt.want.Created = ans.Created
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
			err = delTask(db, id)
			if err != nil {
				log.Fatal(err)
			}
		})
	}
}

func TestEditTask(t *testing.T) {

	var tests = []struct {
		name  string
		input Task
		want  Task
	}{
		{"edit name", Task{Name: "test2"}, Task{Name: "test2", Status: todo}},
		{"edit tag", Task{Tag: "x"}, Task{Name: "test", Tag: "x", Status: todo}},
		{"edit status to inProgress", Task{Status: inProgress}, Task{Name: "test", Status: inProgress}},
		{"edit status to todo", Task{Status: todo}, Task{Name: "test", Status: todo}},
		{"edit description", Task{Desc: "asdf"}, Task{Name: "test", Status: todo, Desc: "asdf"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db = setupTests()
			defer teardownTests(db)
			// put task in db
			id, err := addTask(db, "test", "", todo, "")
			if err != nil {
				log.Fatal(err)
			}
			tt.input.ID = id
			// edit it
			err = editTask(db, tt.input)
			if err != nil {
				log.Fatal(err)
			}
			// get updated task
			ans, err := getTask(db, id)
			if err != nil {
				log.Fatal(err)
			}
			// set fields that we don't know in advance
			tt.want.ID = id
			tt.want.Created = ans.Created
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
			err = delTask(db, id)
			if err != nil {
				log.Fatal(err)
			}
		})
	}
}

func TestDelTask(t *testing.T) {

	var tests = []struct {
		name string
		want Task
	}{
		{"deletes a task", Task{Name: "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db = setupTests()
			defer teardownTests(db)
			// put task in db
			id, err := addTask(db, "test", "", todo, "")
			if err != nil {
				log.Fatal(err)
			}
			// set fields that we don't know in advance
			// get created task
			ans, err := getTask(db, id)
			if err != nil {
				log.Fatal(err)
			}
			tt.want.ID = ans.ID
			tt.want.Created = ans.Created
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
			err = delTask(db, id)
			if err != nil {
				log.Fatal(err)
			}
			tasks, err := getTasks(db)
			if err != nil {
				log.Fatal(err)
			}
			if len(tasks) != 0 {
				t.Errorf("Expected tasks table to be empty but got %v", tasks)
			}
		})
	}
}

func setupTests() *sql.DB {
	var dbPath = filepath.Join(os.TempDir(), "test.db")
	// start the database
	var db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Query("SELECT * FROM tasks"); err == nil {

	} else {
		sqlStatement := `
            CREATE TABLE "tasks" (
                "id" INTEGER NOT NULL PRIMARY KEY,
                "name" TEXT NOT NULL,
                "description" TEXT,
                "status" INTEGER,
                "created" TEXT,
                "tag" TEXT
            );
            DELETE FROM tasks;
        `
		_, err = db.Exec(sqlStatement)
		if err != nil {
			log.Fatal(err)
		}
	}
	return db
}

func teardownTests(db *sql.DB) {
	db.Close()
}
