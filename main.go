package main

import (
	"fmt"

	"github.com/ichiban/prolog"
)

// zanzibar tuple
type Tuple struct {
	resource string
	relation string
	user     string
}

func (t *Tuple) String() string {
	return t.resource + "#" + t.relation + "@" + t.user
}

// return of prolog
type Quest struct {
	X string
	Y string
}

// example
// query = reader(X,Y).

func ZTuples(p *prolog.Interpreter, relation string) {
	query := relation + "(X,Y)."
	sols, err := p.Query(query)
	if err != nil {
		panic(err)
	}
	defer sols.Close()
	var q Quest
	var t Tuple

	for sols.Next() {
		if err := sols.Scan(&q); err != nil {
			panic(err)
		}
		t = Tuple{q.Y, relation, q.X}
		fmt.Printf("%s \n", t.String())
	}

}

func Query(p *prolog.Interpreter, query string) string {
	sols, err := p.Query(query)
	if err != nil {
		panic(err)
	}
	defer sols.Close()
	var q Quest
	var r string = ""

	results := make(map[string]struct{})

	fmt.Printf("? %s -->\n", query)
	for sols.Next() {

		if err := sols.Scan(&q); err != nil {
			panic(err)
		}
		var result string
		if q.Y != "" {
			result = fmt.Sprintf("X= %s  Y= %s ", q.X, q.Y)
		} else {
			result = fmt.Sprintf("X= %s ", q.X)
		}

		if _, exists := results[result]; !exists {
			r = r + fmt.Sprintln(result)
			results[result] = struct{}{}
		}
	}
	return r
}

func main() {
	var r string
	p := prolog.New(nil, nil)

	if err := p.Exec(`
	    /** facts  in prolog are tuples in zanzibar */
		/* zanzibar notation for a tuple is object#relation@user */
		/* reader indicates that the user is a reader on the document. */
		/* writer indicates that the user is a writer on the document. */

	    relation(user1,reader,document1).
		relation(user1,reader,document2).
		relation(user1,reader,document4).
		relation(user1,writer,document4).
		relation(user2,reader,document1).
		relation(user2,reader,document4).
		relation(user1,writer,document1).
		relation(user3,reader,document1).
		relation(user3,reader,document3).
		relation(user3,writer,document3).
		
		reader(X,Y) :- relation(X,reader,Y).
		writer(X,Y) :- relation(X,writer,Y).

		/** rules are permissions in zanzibar */

	    /* edit indicates that the user has permission to edit the document */
		/* view indicates that the user has permission to view the document, */
		/* if they are a reader or have edit permission. */

		edit(X,Y) :- writer(X,Y).
		view(X,Y) :- reader(X,Y) ; edit(X,Y). 


	`); err != nil {
		panic(err)
	}

	// Run the Prolog program.

	// print DATA ztuples
	ZTuples(p, `reader`)
	ZTuples(p, `writer`)

	// expand : Who can view the document Y ou the document4

	r = Query(p, `view(X,Y).`)
	fmt.Println(r)
	r = Query(p, `view(X,document4).`)
	fmt.Println(r)

	// reverse : what can user do ?

	r = Query(p, `edit(user1, X).`)
	fmt.Println(r)
	r = Query(p, `view(user1, X).`)
	fmt.Println(r)

}
