package ctxkit

import (
	"context"
	"fmt"
	"time"
)

func ExampleDetach() {
	// Mock the request context.
	ctx := context.Background()
	ctx = context.WithValue(ctx, "foo", "bar")
	ctx, cancel := context.WithCancel(ctx)

	// Fork the request context, and trigger cancellation.
	dCtx := Detach(ctx)
	cancel()
	fmt.Println("ctx.Err():", ctx.Err())                 // ctx.Err(): context canceled
	fmt.Println("dCtx.Err():", dCtx.Err())               // dCtx.Err(): <nil>
	fmt.Println(`dCtx.Value("foo"):`, dCtx.Value("foo")) // dCtx.Value("foo"): bar

	// Create a child context of dCtx with a a deadline.
	cdCtx, cancel := context.WithTimeout(dCtx, time.Millisecond)
	cancel()
	_, ok := cdCtx.Deadline()
	fmt.Println("cdCtx.Deadline():", ok)     // cdCtx.Deadline(): true, because cancel() is triggered
	fmt.Println("cdCtx.Err():", cdCtx.Err()) // cdCtx.Err(): context canceled, as above
	fmt.Println("dCtx.Err():", dCtx.Err())   // dCtx.Err(): <nil>, cancel() impacts cdCtx not dCtx

	// Output:
	// ctx.Err(): context canceled
	// dCtx.Err(): <nil>
	// dCtx.Value("foo"): bar
	// cdCtx.Deadline(): true
	// cdCtx.Err(): context canceled
	// dCtx.Err(): <nil>
}
