// Package mcrpc provides a JSON-RPC 2.0 over WebSocket client for managing
// Minecraft Java Edition servers. It implements the Minecraft Server Management
// Protocol (MSMP) and supports server management, player management, bans,
// operators, allowlists, server settings, and real-time notifications.
//
// A client is built, then started:
//
//	client := mcrpc.New("localhost:8080", secret,
//	    mcrpc.WithTLS(nil),
//	    mcrpc.WithHandler(mcrpc.Handler{
//	        OnPlayerJoined: func(p mcrpc.Player) { log.Printf("%s joined", p.Name) },
//	    }),
//	)
//
//	if err := client.Start(ctx); err != nil {
//	    return err
//	}
//	defer client.Close()
//
//	players, err := client.GetPlayers(ctx)
//
// New performs no I/O, so notification handlers are in place before the
// connection exists and can never be missed or raced against. The context
// given to Start bounds the session: cancelling it closes the connection.
package mcrpc
