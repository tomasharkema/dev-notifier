package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
	udev "github.com/jochenvg/go-udev"
)

var (
	follow = kingpin.Flag("follow", "Watch mode").Short('f').Bool()
)

func main() {
	root := context.Background()

	ctx, stop := signal.NotifyContext(root, os.Kill, os.Interrupt)
	defer stop()

	kingpin.Parse()

	if *follow {
		runFollow(ctx)
	}
}

func runFollow(ctx context.Context) {
	fmt.Println("Follow")
	u := udev.Udev{}
	m := u.NewMonitorFromNetlink("udev")

	ch, _, err := m.DeviceChan(ctx)
	if err != nil {
		panic(err)
	}

	conn, err := dbus.SessionBusPrivate()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err = conn.Auth(nil); err != nil {
		panic(err)
	}

	if err = conn.Hello(); err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("Started listening on channel")
		for device := range ch {
			handleDevice(conn, device)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("done")
	}

	fmt.Println("Channel closed")
}

func handleDevice(conn *dbus.Conn, device *udev.Device) {
	action := device.Action()
	devnode := device.Devnode()
	properties := device.Properties()
	devlinks := device.Devlinks()

	name := properties["ID_SERIAL"]

	fmt.Println("\n=== New event: ===")
	fmt.Println("Event:", action, device.Syspath())
	fmt.Println(devnode)

	devlinksString := strings.Builder{}
	for k, _ := range devlinks {
		fmt.Println(k)
		fmt.Fprintf(&devlinksString, "%s\n", k)
	}

	// kvbuffer := strings.Builder{}
	for k, v := range properties {
		// fmt.Fprintf(&kvbuffer, " %s: %s;", k, v)
		fmt.Println(k, v)
	}
	// fmt.Println(kvbuffer.String())

	if action == "add" && devnode != "" && name != "" {

		message := fmt.Sprintf("%s\n%s\n%s", name, devnode, devlinksString.String())

		fmt.Println("Send message:", message)

		n := notify.Notification{
			AppName: "Device Notification",
			// ReplacesID: uint32(0),
			// AppIcon:    iconName,
			Summary: "Device Added",
			Body:    message,
			// Actions: []notify.Action{
			//   {Key: "cancel", Label: "Cancel"},
			//   {Key: "open", Label: "Open"},
			// },
			// Hints: map[string]dbus.Variant{
			//   soundHint.ID: soundHint.Variant,
			// },
			ExpireTimeout: time.Minute,
		}

		notifier, err := notify.New(
			conn,
		)
		if err != nil {
			fmt.Println("error sending notification:", err)
		}

		id, err := notifier.SendNotification(n)
		if err != nil {
			fmt.Println("error sending notification:", err)
		}
		fmt.Println("id", id)
	}
}
